package engine

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
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

			blockState := metadata.state
			blockState.Fire(EngineBlockEventReorg)

			if !blockState.IsValid() {
				logger.Panic("invalid block state encountered", zap.Uint8("state", uint8(blockState)))
			}

			// Only process block with Reorged state
			if blockState != EngineBlockStateReorged {
				continue
			}

			e.metadataTracker.SetBlockState(block, blockState)

			artifacts, err := e.serviceEngine.HandleReorgedLogs(block.Logs, metadata.artifacts)
			if err != nil {
				return errors.Wrap(err, "e.serviceEngine.HandleReorgedBlockLogs failed")
			}
			for k, v := range artifacts {
				e.debugMsg("got handleReorgedLogs artifacts", zap.Any("k", k), zap.Any("v", v))
			}

			blockState.Fire(EngineBlockEventHandleReorg)

			metadata.artifacts = artifacts
			metadata.state = blockState

			e.debugMsg("saving handleReorgedLogs metadata for block", zap.Any("metadata", metadata))
			e.metadataTracker.SetBlockMetadata(block, metadata)
		}

		for _, block := range result.GoodBlocks {
			metadata := e.metadataTracker.GetBlockMetadata(block)

			metadata.state.Fire(EngineBlockEventGotLog)
			if !metadata.state.IsValid() {
				logger.Panic("invalid block state encountered", zap.Uint8("state (uint8)", uint8(metadata.state)))
			}

			// Only process block with Seen state
			if metadata.state != EngineBlockStateSeen {
				e.debugMsg(
					"skip seen good block logs",
					zap.Uint64("blockNumber", metadata.blockNumber),
					zap.String("blockHash", metadata.blockHash),
				)

				continue
			}

			artifacts, err := e.serviceEngine.HandleGoodLogs(block.Logs, metadata.artifacts)
			if err != nil {
				return errors.Wrap(err, "serviceEngine.HandleGoodBlockLogs failed")
			}

			metadata.state.Fire(EngineBlockEventProcess)
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
