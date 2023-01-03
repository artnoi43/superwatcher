package superwatcher

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type SuperWatcher interface {
	// Run is the entry point for SuperWatcher
	Run(context.Context, context.CancelFunc) error
	Emitter() Emitter
	Engine() Engine
	Shutdown()

	Controller
}

// Controller gives users the means and methods to change some of EmitterPoller parameters
type Controller interface {
	// SetDoReorg makes the EmitterPoller engage chain reorg detection logic
	SetDoReorg(bool)
	// DoReorg returns if EmitterPoller is currently processing chain reorg inside EmitterPoller.Poll
	DoReorg() bool
	// SetDoHeader makes the EmitterPoller fetch block header for every block with interesting logs
	SetDoHeader(bool)
	// DoHeader returns if the EmitterPoller will fetch block headers for blocks with interesting logs
	DoHeader() bool
	// Addresses reads EmitterPoller's current event log addresses for filter query
	Addresses() []common.Address
	// Topics reads EmitterPoller's current event log topics for filter query
	Topics() [][]common.Hash
	// AddAddresses adds (appends) addresses to EmitterPoller's filter query
	AddAddresses(...common.Address)
	// AddTopics adds (appends) topics to EmitterPoller's filter query
	AddTopics(...[]common.Hash)
	// SetAddresses changes EmitterPoller's event log addresses on-the-fly
	SetAddresses([]common.Address)
	// SetTopics changes EmitterPoller's event log topics on-the-fly
	SetTopics([][]common.Hash)
}
