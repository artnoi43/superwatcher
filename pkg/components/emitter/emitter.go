package emitter

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitter"
)

// New returns a new, default superwatcher.WatcherEmitter.
func New(
	conf *config.EmitterConfig,
	client superwatcher.EthClient,
	stateDataGateway superwatcher.GetStateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	syncChan <-chan struct{}, // Send-receive so that emitter can close this chan
	filterResultChan chan<- *superwatcher.FilterResult,
	errChan chan<- error,
) superwatcher.WatcherEmitter {
	return emitter.New(
		conf,
		client,
		stateDataGateway,
		addresses,
		topics,
		syncChan,
		filterResultChan,
		errChan,
	)
}
