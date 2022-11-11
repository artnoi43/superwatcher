package reorgsim

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestFoo(t *testing.T) {
	mappedLogs := initLogs()
	chain, reorgedChain := newBlockChain(mappedLogs, 15944444)

	fmt.Println("old chain")
	prontBlockChain(chain)

	fmt.Println("reorged chain")
	prontBlockChain(reorgedChain)
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
