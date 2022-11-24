package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	filenames := os.Args[1:]
	if len(filenames) == 0 {
		return
	}

	for _, filename := range filenames {
		b, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("error reading", filename, err.Error())
			continue
		}

		var logs []types.Log
		if err := json.Unmarshal(b, &logs); err != nil {
			fmt.Println("failed to unmarshal logs in", filename)
		}

		for _, log := range logs {
			fmt.Println("block", log.BlockNumber, "blockHash", log.BlockHash.String(), "txHash", log.TxHash.String())
		}
	}
}
