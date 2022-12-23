package components

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/poller"
)

func NewPoller(
	addresses []common.Address,
	topics [][]common.Hash,
	doReorg bool,
	filterRange uint64,
	filterFunc superwatcher.FilterFunc,
	logLevel uint8,
) superwatcher.EmitterPoller {
	return poller.New(
		addresses,
		topics,
		doReorg,
		filterRange,
		filterFunc,
		logLevel,
	)
}
