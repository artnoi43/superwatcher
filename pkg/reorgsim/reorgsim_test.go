package reorgsim

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

func initChains() (blockChain, blockChain) {
	mappedLogs := InitLogs()
	return NewBlockChain(mappedLogs, 15944444)
}

func TestFoo(t *testing.T) {
	chain, reorgedChain := initChains()

	fmt.Println("old chain")
	prontBlockChain(chain)

	fmt.Println("reorged chain")
	prontBlockChain(reorgedChain)

	sim := NewReorgSim(5, 15944400, 15944444)
	filterLogs, err := sim.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(15944401),
		ToBlock:   big.NewInt(15944500),
	})

	if err != nil {
		t.Fatal(err.Error())
	}
	filterLogsMapped := mapLogsToNumber(filterLogs)
	fmt.Println("filterLogs")
	prontMapLen(filterLogsMapped, "blockNumber", "len(logs)")
}

func prontMapLen[T comparable, U any](m map[T][]U, keyString, lenString string) {
	for k, arr := range m {
		fmt.Println(keyString, k, lenString, len(arr))
	}
}

func prontLogs(logs []types.Log) {
	for _, log := range logs {
		fmt.Println("blockNumber", log.BlockNumber, "addr", log.Address)
	}
}

func prontBlockChain(chain blockChain) {
	for _, b := range chain {
		fmt.Println("blockNumber", b.blockNumber, "len(logs)", len(b.logs), "reorgedHere", b.reorgedHere)
	}
}
