package emitter

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/domain/datagateway"
	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
	"github.com/artnoi43/superwatcher/lib/enums"
	"github.com/artnoi43/superwatcher/lib/logger/debug"
)

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
	config     *config.Config
	client     ethClient
	tracker    *reorg.Tracker
	startBlock uint64
	addresses  []common.Address
	topics     [][]common.Hash

	// Redis-store for tracking last recorded block
	stateDataGateway datagateway.StateDataGateway

	// These fields are gateways via which
	// external services interact with emitter
	filterResultChan chan<- *FilterResult
	errChan          chan<- error
	syncChan         <-chan struct{}

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

// Loop wraps loopFilterLogs to provide graceful shutdown mechanism for emitter.
// When ctx is camcled else where, Loop calls *emitter.shutdown and returns ctx.Err()
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
	close(e.filterResultChan)
	close(e.errChan)
}

func (e *emitter) debugMsg(msg string, fields ...zap.Field) {
	debug.DebugMsg(e.debug, msg, fields...)
}
