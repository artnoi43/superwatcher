package engine

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
)

// engineBlocks are used to aggregate multiple blocks' information
// from superwatcher.FilterResult to pass to ServiceEngine methods.
type engineBlocks struct {
	logs []*types.Log

	metadata  []*blockMetadata
	artifacts []superwatcher.Artifact
}

func (e *engine) handleResults(ctx context.Context) error {
	// Get emitterConfig to clear tracker metadata based on FilterRange
	emitterConfig := e.emitterClient.WatcherConfig()

	for {
		result := e.emitterClient.WatcherResult()
		if result == nil {
			e.debugger.Debug(1, "handleLogs got nil result, emitterClient was probably shutdown, returning..")
			return nil
		}

		var reorged bool
		var shouldCallServiceEngine bool

		var reorgedBlocks engineBlocks
		for _, block := range result.ReorgedBlocks {
			reorged = true
			metadata := e.metadataTracker.GetBlockMetadata(callerReorgedLogs, block.Number, block.String())
			e.debugger.Debug(
				3, "* got reorged metadata",
				zap.Uint64("blockNumber", block.Number),
				zap.String("blockHash", block.String()),
				zap.Any("metadata artifacts", metadata.artifacts),
			)

			metadata.state.Fire(EventReorg)
			// Only process block with Reorged state
			if metadata.state != StateReorged {
				e.debugger.Debug(
					1, "skip bad reorged block logs",
					zap.String("state", metadata.state.String()),
					zap.Uint64("blockNumber", metadata.blockNumber),
					zap.String("blockHash", metadata.blockHash),
				)

				continue
			}

			shouldCallServiceEngine = true

			reorgedBlocks.logs = append(reorgedBlocks.logs, block.Logs...)
			reorgedBlocks.metadata = append(reorgedBlocks.metadata, metadata)
			reorgedBlocks.artifacts = append(reorgedBlocks.artifacts, metadata.artifacts)
		}

		var reorgedArtifacts map[common.Hash][]superwatcher.Artifact
		if shouldCallServiceEngine {
			var err error
			reorgedArtifacts, err = e.serviceEngine.HandleReorgedLogs(reorgedBlocks.logs, reorgedBlocks.artifacts)
			if err != nil {
				return errors.Wrap(err, "serviceEngine.HandleReorgedBlockLogs failed")
			}

			shouldCallServiceEngine = false
		}

		// Check debug here so we dont have to iterate over all |artifacts| members
		// before checking
		if e.debug {
			for k, v := range reorgedArtifacts {
				e.debugger.Debug(2, "got handleReorgedLogs artifacts", zap.Any("k", k), zap.Any("v", v))
			}
		}

		// Update metadata for reorged blocks
		for _, metadata := range reorgedBlocks.metadata {
			metadata.state.Fire(EventHandleReorg)
			metadata.artifacts = reorgedArtifacts[common.HexToHash(metadata.blockHash)]

			e.debugger.Debug(
				4, "* saving reorgedBlock metadata",
				zap.Uint64("blockNumber", metadata.blockNumber),
				zap.String("blockHash", metadata.blockHash),
				zap.Any("metadata artifacts", metadata.artifacts),
			)
			e.metadataTracker.SetBlockMetadata(callerReorgedLogs, metadata)
		}

		var goodBlocks engineBlocks
		for _, block := range result.GoodBlocks {
			metadata := e.metadataTracker.GetBlockMetadata(callerGoodLogs, block.Number, block.String())
			metadata.state.Fire(EventGotLog)

			// Update state to tracker
			e.metadataTracker.SetBlockMetadata(callerGoodLogs, metadata)

			// Only process block with Seen state
			if metadata.state != StateSeen {
				e.debugger.Debug(
					1, "skip goodBlock",
					zap.String("state", metadata.state.String()),
					zap.Uint64("blockNumber", metadata.blockNumber),
					zap.String("blockHash", metadata.blockHash),
				)

				continue
			}

			shouldCallServiceEngine = true

			goodBlocks.logs = append(goodBlocks.logs, block.Logs...)
			goodBlocks.metadata = append(goodBlocks.metadata, metadata)
			reorgedBlocks.artifacts = append(reorgedBlocks.artifacts, metadata.artifacts)
		}

		var artifacts map[common.Hash][]superwatcher.Artifact
		if shouldCallServiceEngine {
			var err error
			artifacts, err = e.serviceEngine.HandleGoodLogs(goodBlocks.logs, goodBlocks.artifacts)
			if err != nil {
				return errors.Wrap(err, "serviceEngine.HandleGoodBlockLogs failed")
			}
		}

		for _, metadata := range goodBlocks.metadata {
			metadata.state.Fire(EventProcess)
			metadata.artifacts = artifacts[common.HexToHash(metadata.blockHash)]

			e.debugger.Debug(
				3, "* saving goodBlock metadata",
				zap.Uint64("blockNumber", metadata.blockNumber),
				zap.String("blockHash", metadata.blockHash),
				zap.Any("metadata artifacts", metadata.artifacts),
			)

			e.metadataTracker.SetBlockMetadata(callerGoodLogs, metadata)
		}

		// TODO: How many should we clear?
		e.metadataTracker.ClearUntil(
			result.LastGoodBlock - (emitterConfig.FilterRange * emitterConfig.GoBackRetries),
		)

		var lastRecordedBlock uint64
		if reorged {
			lastRecordedBlock = result.LastGoodBlock
		} else {
			lastRecordedBlock = result.ToBlock
		}

		if err := e.stateDataGateway.SetLastRecordedBlock(ctx, result.LastGoodBlock); err != nil {
			return errors.Wrapf(err, "failed to save lastRecordedBlock %d", lastRecordedBlock)
		}

		e.emitterClient.WatcherEmitterSync()
	}
}
