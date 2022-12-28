package mock

// TODO: does not work yet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/artnoi43/superwatcher"
	"github.com/ethereum/go-ethereum/common"
)

type mockPoller struct {
	currentIndex  int
	lastGood      uint64
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
	reorgedBlock := p.reorgedBlocks[p.currentIndex]
	noMoreReorg := len(p.reorgedBlocks) == p.currentIndex+1

	for number := fromBlock; number <= toBlock; number++ {
		p.seen[number]++
	}

	result := &superwatcher.FilterResult{FromBlock: fromBlock, ToBlock: toBlock}

	if p.seen[reorgedBlock] < 1 || noMoreReorg {
		result.LastGoodBlock = toBlock
		p.lastGood = toBlock
		for number := fromBlock; number <= toBlock; number++ {
			result.GoodBlocks = append(result.GoodBlocks, &superwatcher.BlockInfo{
				Number: number,
				Hash:   common.BigToHash(big.NewInt(int64(number))),
			})
		}
		return result, nil
	}

	p.currentIndex++

	if reorgedBlock == fromBlock {
		return nil, superwatcher.ErrFromBlockReorged
	}

	for number := fromBlock; number < reorgedBlock; number++ {
		result.GoodBlocks = append(result.GoodBlocks, &superwatcher.BlockInfo{
			Number: number,
			Hash:   common.BigToHash(big.NewInt(int64(number))),
		})
	}
	for number := p.lastGood + 1; number <= toBlock; number++ {
		result.GoodBlocks = append(result.GoodBlocks, &superwatcher.BlockInfo{
			Number: number,
			Hash:   common.BigToHash(big.NewInt(int64(number))),
		})
	}

	for number := reorgedBlock; number <= p.lastGood; number++ {
		fmt.Println("ello", number)
		result.ReorgedBlocks = append(result.ReorgedBlocks, &superwatcher.BlockInfo{
			Number: number,
			Hash:   common.BigToHash(big.NewInt(int64(number))),
		})
	}

	result.LastGoodBlock = superwatcher.LastGoodBlock(result)
	p.lastGood = result.LastGoodBlock

	return result, nil
}

func (p *mockPoller) SetDoReorg(bool)                {}
func (p *mockPoller) DoReorg() bool                  { return true }
func (p *mockPoller) Addresses() []common.Address    { return nil }
func (p *mockPoller) Topics() [][]common.Hash        { return nil }
func (p *mockPoller) AddAddresses(...common.Address) {}
func (p *mockPoller) AddTopics(...[]common.Hash)     {}
func (p *mockPoller) SetAddresses([]common.Address)  {}
func (p *mockPoller) SetTopics([][]common.Hash)      {}
