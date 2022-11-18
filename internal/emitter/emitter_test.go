package emitter

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

func TestEmitter(t *testing.T) {
	// to run test all cases a once, set timeout to ~1 min
	t.Run("case testEmitterCase1", func(t *testing.T) {
		emitterTestTemplate(t, testCases[0])
	})
	t.Run("case testEmitterCase2", func(t *testing.T) {
		emitterTestTemplate(t, testCases[1])
	})
	t.Run("case testEmitterCase3", func(t *testing.T) {
		emitterTestTemplate(t, testCases[2])
	})
	t.Run("case testEmitterCase4", func(t *testing.T) {
		emitterTestTemplate(t, testCases[3])
	})
	t.Run("case testEmitterCase5", func(t *testing.T) {
		emitterTestTemplate(t, testCases[4])
	})
}

func emitterTestTemplate(t *testing.T, tc testConfig) {
	b, _ := json.Marshal(tc)
	t.Logf("testConfig: %s", b)

	syncChan := make(chan struct{})
	filterResultChan := make(chan *superwatcher.FilterResult)
	errChan := make(chan error)

	conf, _ := config.ConfigYAML("../../config/config.yaml")
	// Override LoopInterval
	conf.LoopInterval = 0

	fakeDataGateway := &mockStateDataGateway{value: tc.FromBlock - 1}
	sim := reorgsim.NewReorgSim(conf.LookBackBlocks, tc.FromBlock-1, tc.ReorgedAt, tc.LogsFiles)
	e := New(conf, sim, fakeDataGateway, nil, nil, syncChan, filterResultChan, errChan, false)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		e.(*emitter).loopFilterLogs(ctx)
	}()

	tracker := newTracker()
	var seenLogs []*types.Log
	for {
		result := <-filterResultChan
		if result.LastGoodBlock > tc.ToBlock {
			cancel()
			break
		}
		lastGoodBlock := result.LastGoodBlock

		// check ReorgedBlocks
		for i, block := range result.ReorgedBlocks {
			blockNumber := block.Number
			hash := block.Hash

			if b, exists := tracker.getTrackerBlockInfo(blockNumber); exists {
				if b.Hash == hash {
					t.Fatalf("ReorgedBlocks[%d] is not reorg: hash=%v", i, hash)
				}
			} else {
				t.Fatalf("ReorgedBlocks[%d] didn't check before", i)
			}

			tracker.addTrackerBlockInfo(block)

			// check LastGoodBlock
			if blockNumber < lastGoodBlock {
				t.Fatalf("LastGoodBlock is wrong: ReorgedBlocks[%d] blockNumber=%v, LastGoodBlock=%v", i, blockNumber, lastGoodBlock)
			}

			for _, log := range block.Logs {
				if !gslutils.Contains(seenLogs, log) {
					t.Fatalf(
						"reorgedLog not seen before: blockHash %s, txHash %s, addr %s, topics[0] %s",
						log.BlockHash.String(), log.TxHash.String(), log.Address.String(), log.Topics[0].String(),
					)
				}
			}
		}

		// check GoodBlocks
		for i, block := range result.GoodBlocks {
			blockNumber := block.Number
			hash := block.Hash

			if b, exists := tracker.getTrackerBlockInfo(blockNumber); exists {
				if b.Hash != hash {
					t.Fatalf("GoodBlocks[%d] is reorg: hash(before)=%v hash(after)=%v", i, b.Hash, hash)
				}
			}
			tracker.addTrackerBlockInfo(block)

			for _, log := range block.Logs {
				appendUnique(seenLogs, log)
			}
		}

		e.(*emitter).stateDataGateway.SetLastRecordedBlock(ctx, result.LastGoodBlock)
		syncChan <- struct{}{}
	}
}

func appendUnique[T comparable](arr []T, item T) []T {
	if gslutils.Contains(arr, item) {
		return append(arr, item)
	}

	return arr
}

type mockStateDataGateway struct {
	value uint64
}

func (m *mockStateDataGateway) GetLastRecordedBlock(context.Context) (uint64, error) {
	return m.value, nil
}

func (m *mockStateDataGateway) SetLastRecordedBlock(ctx context.Context, v uint64) error {
	m.value = v
	return nil
}

func (m *mockStateDataGateway) Shutdown() error {
	return nil
}
