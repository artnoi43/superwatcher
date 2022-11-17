package emitter

import (
	"context"
	"testing"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/ethereum/go-ethereum/common"
)

func TestEmitterCase1(t *testing.T) {
	syncChan := make(chan struct{})
	filterResultChan := make(chan *superwatcher.FilterResult)
	errChan := make(chan error)

	conf, _ := config.ConfigYAML("../../config/config.yaml")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tc := testCases[0]
	fakeDataGateway := &mockStateDataGateway{value: tc.fromBlock}
	sim := reorgsim.NewReorgSim(5, tc.fromBlock-1, tc.reorgedAt, tc.logs)
	e := New(conf, sim, fakeDataGateway, nil, nil, syncChan, filterResultChan, errChan, true)

	go func() {
		e.(*emitter).loopFilterLogs(ctx)
	}()

	results := make(map[uint64]common.Hash)
	for {
		result := <-filterResultChan
		lastGoodBlock := result.LastGoodBlock

		// check ReorgedBlocks
		for i, block := range result.ReorgedBlocks {
			blockNumber := block.Number
			hash := block.Hash

			if val, exists := results[blockNumber]; exists {
				if val == hash {
					t.Fatalf("ReorgedBlocks[%d] is not reorg: hash=%v", i, val)
				}
			} else {
				t.Fatalf("ReorgedBlocks[%d] didn't check before", i)
			}
			results[blockNumber] = hash

			// check LastGoodBlock
			if blockNumber < lastGoodBlock {
				t.Fatalf("LastGoodBlock is wrong: ReorgedBlocks[%d] blockNumber=%v, LastGoodBlock=%v", i, blockNumber, lastGoodBlock)
			}
		}


		// check GoodBlocks
		for i, block := range result.GoodBlocks {
			blockNumber := block.Number
			hash := block.Hash

			if val, exists := results[blockNumber]; exists {
				if val != hash {
					t.Fatalf("GoodBlocks[%d] is reorg: hash(before)=%v hash(after)=%v", i, val, hash)
				}
			}
			results[blockNumber] = hash
		}

		syncChan <- struct{}{}
	}
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
