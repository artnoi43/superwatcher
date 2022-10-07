package watcher

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
	"github.com/artnoi43/superwatcher/domain/usecase/watcher/reorg"
	"github.com/artnoi43/superwatcher/lib/enums"
	"github.com/artnoi43/superwatcher/lib/logger"
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
type Watcher interface {
	Loop(context.Context) error
}

// watcher implements Watcher, and other than Config,
// other fields of this structure are defined as ifaces,
// to facil mock testing.
type watcher struct {
	// These fields are used for filtering event logs
	config           *config.Config
	client           ethClient
	dataGateway      datagateway.DataGateway
	stateDataGateway datagateway.StateDataGateway
	tracker          *reorg.Tracker
	startBlock       uint64
	addresses        []common.Address
	topics           [][]common.Hash

	// These fields are comms with for other services
	logChan   chan<- *types.Log
	reorgChan chan<- *struct{}
	errChan   chan<- error
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
func NewWatcher(
	conf *config.Config,
	client ethClient,
	dataGateway datagateway.DataGateway,
	stateDataGateway datagateway.StateDataGateway,
	addresses []common.Address,
	topics [][]common.Hash,
	logChan chan<- *types.Log,
	errChan chan<- error,
	reorgChan chan<- *struct{},
) *watcher {
	logger.Debug("initializing watcher", zap.Any("addresses", addresses), zap.Any("topics", topics))
	return &watcher{
		config:           conf,
		client:           client,
		dataGateway:      dataGateway,
		stateDataGateway: stateDataGateway,
		tracker:          reorg.NewTracker(),
		startBlock:       conf.StartBlock,
		addresses:        addresses,
		topics:           topics,
		logChan:          logChan,
		errChan:          errChan,
		reorgChan:        reorgChan,
	}
}

// Loop wraps loopFilterLogs with graceful shutdown code.
func (w *watcher) Loop(ctx context.Context) error {
	for {
		// NOTE: this is not clean, but a workaround to prevent infinity loop
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := w.loopFilterLogs(ctx); err != nil {
				w.errChan <- errors.Wrap(err, "error in loopFilterLogs")
			}
		}
	}
}
