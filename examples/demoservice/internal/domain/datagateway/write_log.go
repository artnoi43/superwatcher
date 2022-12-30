package datagateway

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type WriteLog string

// For demotest code
func (s WriteLog) Unmarshal() (string, string, uint64, common.Hash, error) {
	splitted := strings.Split(string(s), " ")
	if len(splitted) < 7 {
		return "", "", 0, common.Hash{}, fmt.Errorf("malformed ensWriteLogs: %s", s)
	}

	if splitted[0] != "SET" && splitted[0] != "DEL" {
		return "", "", 0, common.Hash{}, fmt.Errorf("malformed ensWriteLogs - unexpected key: %s from %s", splitted[0], splitted)
	}

	if splitted[3] != "BLOCK" {
		return "", "", 0, common.Hash{}, fmt.Errorf("malformed ensWriteLogs - missing \"BLOCK\" at index 3: %s", s)
	}

	if splitted[5] != "HASH" {
		return "", "", 0, common.Hash{}, fmt.Errorf("malformed ensWriteLogs - missing \"HASH\" at index 5: %s", s)
	}

	blockNumber, err := strconv.ParseInt(splitted[4], 10, 64)
	if err != nil {
		return "", "", 0, common.Hash{}, fmt.Errorf("malformed ensWriteLogs: %s", s)
	}

	blockHash := common.HexToHash(splitted[6])

	return splitted[0], splitted[2], uint64(blockNumber), blockHash, nil
}
