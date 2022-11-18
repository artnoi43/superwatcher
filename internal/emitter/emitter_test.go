package emitter

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

func TestEmitterAllCases(t *testing.T) {
	for i := range testCases {
		testName := fmt.Sprintf("Case:%d", i+1)
		t.Run(testName, func(t *testing.T) {
			emitterTestTemplate(t, i+1)
		})
	}
}

// Run this from the root of the repo with:
// go test -v ./internal/emitter -run TestEmitterByCase -case 69
// Go test binary already called `flag.Parse`, so we just simply
// need to name our flag so that the flag package knows to parse it too.
var flagCase = flag.Int("case", -1, "Emitter test case")

func TestEmitterByCase(t *testing.T) {
	var caseNumber int = 1 // default case 1
	if flagCase != nil {
		caseNumber = *flagCase
	}

	if caseNumber < 0 {
		TestEmitterAllCases(t)
		return
	}

	if len(testCases) > caseNumber {
		testName := fmt.Sprintf("Case:%d", caseNumber)
		t.Run(testName, func(t *testing.T) {
			emitterTestTemplate(t, caseNumber)
		})

		return
	}

	t.Skipf("no such test case: %d", caseNumber)
}

func emitterTestTemplate(t *testing.T, caseNumber int) {
	tc := testCases[caseNumber-1]
	b, _ := json.Marshal(tc)
	t.Logf("testConfig for case %d: %s", caseNumber, b)

	syncChan := make(chan struct{})
	filterResultChan := make(chan *superwatcher.FilterResult)
	errChan := make(chan error)

	conf, _ := config.ConfigYAML("../../config/config.yaml")
	// Override LoopInterval
	conf.LoopInterval = 0

	fakeRedis := &mockStateDataGateway{value: tc.FromBlock - 1}
	sim := reorgsim.NewReorgSim(conf.LookBackBlocks, tc.FromBlock-1, tc.ReorgedAt, tc.LogsFiles)
	testEmitter := New(conf, sim, fakeRedis, nil, nil, syncChan, filterResultChan, errChan, false)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		testEmitter.(*emitter).loopFilterLogs(ctx)
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

			// Check LastGoodBlock
			if lastGoodBlock > blockNumber {
				t.Fatalf(
					"invalid LastGoodBlock: ReorgedBlocks[%d] blockNumber=%v, LastGoodBlock=%v",
					i, blockNumber, lastGoodBlock,
				)
			}

			// Check that all the reorged logs were seen before in |seenLogs|
			for _, log := range block.Logs {
				if !gslutils.Contains(seenLogs, log) {
					t.Fatalf(
						"reorgedLog not seen before: blockHash %s, txHash %s, addr %s, topics[0] %s",
						log.BlockHash.String(), log.TxHash.String(), log.Address.String(), log.Topics[0].String(),
					)
				}
			}

			tracker.addTrackerBlockInfo(block)
		}

		for i, block := range result.GoodBlocks {
			blockNumber := block.Number
			hash := block.Hash

			if b, exists := tracker.getTrackerBlockInfo(blockNumber); exists {
				if b.Hash != hash {
					t.Fatalf(
						"GoodBlocks[%d] is reorged: hash(before)=%v hash(after)=%v",
						i, b.Hash.String(), hash.String(),
					)
				}
			}
			for _, log := range block.Logs {
				var ok bool
				seenLogs, ok = appendUnique(seenLogs, log)
				if !ok {
					t.Fatalf(
						"duplicate good logs seen: blockHash %s, txHash %s, addr %s, topics[0]: %s",
						log.BlockHash.String(), log.TxHash.String(), log.Address.String(), log.Topics[0].String(),
					)
				}
			}

			tracker.addTrackerBlockInfo(block)
		}

		// Sets before syncs
		testEmitter.(*emitter).stateDataGateway.SetLastRecordedBlock(ctx, result.LastGoodBlock)
		syncChan <- struct{}{}
	}
}

// appendUnique appends item to arr if arr does not contain item.
func appendUnique[T comparable](arr []T, item T) ([]T, bool) {
	if !gslutils.Contains(arr, item) {
		return append(arr, item), true
	}

	return arr, false
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
