package mock

// TODO: does not work yet

import (
	"context"
	"math/big"

	"github.com/artnoi43/superwatcher"
	"github.com/ethereum/go-ethereum/common"
)

type mockPoller struct {
	currentIndex  int
	reorgedBlocks []uint64
	seen          map[uint64]int
}

func NewPoller(
	reorgedBlocks []uint64,
) superwatcher.EmitterPoller {
	return &mockPoller{
		reorgedBlocks: reorgedBlocks,
		seen:          make(map[uint64]int),
	}
}

func (p *mockPoller) Poll(
	ctx context.Context,
	fromBlock, toBlock uint64,
) (
	*superwatcher.FilterResult,
	error,
) {
	result := &superwatcher.FilterResult{FromBlock: fromBlock, ToBlock: toBlock}

	reorgedBlock := p.reorgedBlocks[p.currentIndex]
	noMoreReorg := len(p.reorgedBlocks) == p.currentIndex+1

	for number := fromBlock; number <= toBlock; number++ {
		p.seen[number]++
	}

	if p.seen[reorgedBlock] <= 1 || p.seen[reorgedBlock] > 3 {
		for n := fromBlock; n <= toBlock; n++ {
			result.GoodBlocks = append(result.GoodBlocks, &superwatcher.BlockInfo{
				Number: n,
				Hash:   common.BigToHash(big.NewInt(int64(n))),
			})
		}

		result.LastGoodBlock = toBlock

		return result, nil
	}

	if !noMoreReorg {
		p.currentIndex++
		p.seen = make(map[uint64]int)
	}

	// reorgBlock is somewhere between fromBlock -> toBlock
	if fromBlock != reorgedBlock {
		for n := fromBlock; n < reorgedBlock; n++ {
			result.GoodBlocks = append(result.GoodBlocks, &superwatcher.BlockInfo{
				Number: n,
				Hash:   common.BigToHash(big.NewInt(int64(n))),
			})
		}
		for n := reorgedBlock; n <= toBlock; n++ {
			result.ReorgedBlocks = append(result.ReorgedBlocks, &superwatcher.BlockInfo{
				Number: n,
				Hash:   common.BigToHash(big.NewInt(int64(n))),
			})
		}

		result.LastGoodBlock = superwatcher.LastGoodBlock(result)
		return result, nil
	}

	return nil, superwatcher.ErrFromBlockReorged
}

func (p *mockPoller) SetDoReorg(bool)                {}
func (p *mockPoller) DoReorg() bool                  { return true }
func (p *mockPoller) Addresses() []common.Address    { return nil }
func (p *mockPoller) Topics() [][]common.Hash        { return nil }
func (p *mockPoller) AddAddresses(...common.Address) {}
func (p *mockPoller) AddTopics(...[]common.Hash)     {}
func (p *mockPoller) SetAddresses([]common.Address)  {}
func (p *mockPoller) SetTopics([][]common.Hash)      {}
