package emitter

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/domain/datagateway"
	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/lib/enums"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

// ethClient is an interface representing *ethclient.Client methods
// called in watcher methods. Defined as an iface to facil mock testing.
type ethClient interface {
	BlockNumber(context.Context) (uint64, error)
	BlockByNumber(context.Context, *big.Int) (*types.Block, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)

	// Not sure if needed
	// BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// Code that imports watcher should only use this method.
type WatcherEmitter interface {
	Loop(context.Context) error
	shutdown()
}

// emitter implements Watcher, and other than Config,
// other fields of this structure are defined as ifaces,
// to facil mock testing.
type emitter struct {
	// These fields are used for filtering event logs
	config           *config.Config
	client           ethClient
	dataGateway      datagateway.DataGateway // TODO: remove?
	stateDataGateway datagateway.StateDataGateway
	tracker          *reorg.Tracker
	startBlock       uint64
	addresses        []common.Address
	topics           [][]common.Hash

	// These fields are comms with for other services
	logChan   chan<- *types.Log
	blockChan chan<- *reorg.BlockInfo
	reorgChan chan<- *reorg.BlockInfo
	errChan   chan<- error

	// For debugging
	debug bool
}

// Config represents the configuration structure for watcher
type Config struct {
	Chain           enums.ChainType `mapstructure:"chain" json:"chain"`
	Node            string          `mapstructure:"node_url" json:"node"`
	StartBlock      uint64          `mapstructure:"start_block" json:"startBlock"`
	LookBackBlocks  uint64          `mapstructure:"lookback_blocks" json:"lookBackBlock"`
	LookBackRetries uint64          `mapstructure:"lookback_retries" json:"lookBackRetries"`
	IntervalSecond  int             `mapstructure:"interval_second" json:"intervalSecond"`
}

// NewWatcher initializes contract info from config
func NewWatcherDebug(
	conf *config.Config,
	client ethClient,
	dataGateway datagateway.DataGateway,
	stateDataGateway datagateway.StateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	logChan chan<- *types.Log,
	blockChan chan<- *reorg.BlockInfo,
	reorgChan chan<- *reorg.BlockInfo,
	errChan chan<- error,
) WatcherEmitter {
	e := New(
		conf,
		client,
		dataGateway,
		stateDataGateway,
		addresses,
		topics,
		logChan,
		blockChan,
		reorgChan,
		errChan,
	)

	e.(*emitter).debug = true

	return e
}

// Loop wraps loopFilterLogs with graceful shutdown code.
func (e *emitter) Loop(ctx context.Context) error {
	for {
		// NOTE: this is not clean, but a workaround to prevent infinity loop
		select {
		case <-ctx.Done():
			e.shutdown()
			return ctx.Err()
		default:
			if err := e.loopFilterLogs(ctx); err != nil {
				e.errChan <- errors.Wrap(err, "error in loopFilterLogs")
			}
		}
	}
}

func (e *emitter) shutdown() {
	close(e.logChan)
	close(e.blockChan)
	close(e.reorgChan)
	close(e.errChan)
}

func (e *emitter) debugMsg(msg string, fields ...zap.Field) {
	debug.DebugMsg(e.debug, msg, fields...)
}
