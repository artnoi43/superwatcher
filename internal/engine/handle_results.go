package engine

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (e *engine) handleResults(ctx context.Context) error {
	// Get emitterConfig to clear tracker metadata based on lookBackBlocks
	emitterConfig := e.emitterClient.WatcherConfig()

	for {
		result := e.emitterClient.WatcherResult()
		if result == nil {
			e.debugMsg("handleLogs got nil result, emitterClient was probably shutdown, returning..")
			return nil
		}

		for _, block := range result.ReorgedBlocks {
			metadata := e.metadataTracker.GetBlockMetadata(block)
			metadata.state.Fire(EventReorg)

			// Only process block with Reorged state
			if metadata.state != StateReorged {
				e.debugMsg(
					"skip bad reorged block logs",
					zap.String("state", metadata.state.String()),
					zap.Uint64("blockNumber", metadata.blockNumber),
					zap.String("blockHash", metadata.blockHash),
				)

				continue
			}

			e.metadataTracker.SetBlockState(block, metadata.state)

			artifacts, err := e.serviceEngine.HandleReorgedLogs(block.Logs, metadata.artifacts)
			if err != nil {
				return errors.Wrap(err, "e.serviceEngine.HandleReorgedBlockLogs failed")
			}

			// Check debug here so we dont have to iterate over all keys in map artifacts before checking in `e.debugMsg`
			if e.debug {
				for k, v := range artifacts {
					e.debugMsg("got handleReorgedLogs artifacts", zap.Any("k", k), zap.Any("v", v))
				}
			}

			metadata.state.Fire(EventHandleReorg)
			metadata.artifacts = artifacts

			e.debugMsg("saving handleReorgedLogs metadata for block", zap.Any("metadata", metadata))
			e.metadataTracker.SetBlockMetadata(block, metadata)
		}

		for _, block := range result.GoodBlocks {
			metadata := e.metadataTracker.GetBlockMetadata(block)
			metadata.state.Fire(EventGotLog)

			// Only process block with Seen state
			if metadata.state != StateSeen {
				e.debugMsg(
					"skip seen good block logs",
					zap.String("state", metadata.state.String()),
					zap.Uint64("blockNumber", metadata.blockNumber),
					zap.String("blockHash", metadata.blockHash),
				)

				continue
			}

			artifacts, err := e.serviceEngine.HandleGoodLogs(block.Logs, metadata.artifacts)
			if err != nil {
				return errors.Wrap(err, "serviceEngine.HandleGoodBlockLogs failed")
			}

			metadata.state.Fire(EventProcess)
			metadata.artifacts = artifacts

			e.debugMsg("saving metadata for block", zap.Any("metadata", metadata))
			e.metadataTracker.SetBlockMetadata(block, metadata)
		}

		// TODO: How many should we clear?
		e.metadataTracker.ClearUntil(
			result.LastGoodBlock - (emitterConfig.LookBackBlocks * emitterConfig.LookBackRetries),
		)

		if err := e.stateDataGateway.SetLastRecordedBlock(ctx, result.LastGoodBlock); err != nil {
			return errors.Wrapf(err, "failed to save lastRecordedBlock %d", result.LastGoodBlock)
		}

		e.emitterClient.WatcherEmitterSync()
	}
}
