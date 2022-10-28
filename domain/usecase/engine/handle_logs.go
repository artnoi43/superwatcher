package engine

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
)

func (e *engine) handleLogs(ctx context.Context) error {
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

			blockState.Fire(EngineBlockEventHandleReorg)

			metadata.artifacts = artifacts
			metadata.state = blockState

			e.metadataTracker.SetBlockMetadata(block, metadata)
		}

		for _, block := range result.GoodBlocks {
			metadata := e.metadataTracker.GetBlockMetadata(block)

			blockState := metadata.state
			blockState.Fire(EngineBlockEventGotLog)

			if !blockState.IsValid() {
				logger.Panic("invalid block state encountered", zap.Uint8("state", uint8(blockState)))
			}

			// Only process block with Seen state
			if blockState != EngineBlockStateSeen {
				continue
			}

			artifacts, err := e.serviceEngine.HandleGoodLogs(block.Logs, metadata.artifacts)
			if err != nil {
				return errors.Wrap(err, "serviceEngine.HandleGoodBlockLogs failed")
			}

			blockState.Fire(EngineBlockEventProcess)

			metadata.artifacts = artifacts
			metadata.state = blockState

			e.metadataTracker.SetBlockMetadata(block, metadata)
		}

		// TODO: How many should we clear?
		e.metadataTracker.ClearUntil(
			result.LastGoodBlock - (emitterConfig.LookBackBlocks * emitterConfig.LookBackRetries),
		)

		e.stateDataGateway.SetLastRecordedBlock(ctx, result.LastGoodBlock)

		e.emitterClient.WatcherEmitterSync()
	}
}
