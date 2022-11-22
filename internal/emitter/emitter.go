package emitter

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/enums"
	"github.com/artnoi43/superwatcher/pkg/logger/debug"
)

// emitter implements Watcher, and other than Config,
// other fields of this structure are defined as ifaces,
// to facil mock testing.
type emitter struct {
	// These fields are used for filtering event logs
	config     *config.Config
	client     superwatcher.EthClient
	tracker    *blockInfoTracker
	startBlock uint64
	addresses  []common.Address
	topics     [][]common.Hash

	// Redis-store for tracking last recorded block
	stateDataGateway watcherstate.StateDataGateway

	// These fields are gateways via which
	// external services interact with emitter
	filterResultChan chan<- *superwatcher.FilterResult
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
	status := new(filterLogStatus)

	for {
		// NOTE: this is not clean, but a workaround to prevent infinity loop
		select {
		case <-ctx.Done():
			e.debugMsg("shutting down emitter", zap.Any("emitterStatus", status))
			e.Shutdown()
			return errors.Wrap(ctx.Err(), ErrEmitterShutdown.Error())

		default:
			if err := e.loopFilterLogs(ctx, status); err != nil {
				e.debugMsg("loopFilterLogs returned", zap.Any("status", status), zap.Error(err))
				e.emitError(errors.Wrap(err, "error in loopFilterLogs"))
			}
		}
	}
}

func (e *emitter) Shutdown() {
	close(e.filterResultChan)
	close(e.errChan)
}

func (e *emitter) SyncsWithEngine() {
	e.debugMsg("emitter: waiting for engine sync")
	<-e.syncChan
	e.debugMsg("emitter: synced with engine")
}

func (e *emitter) debugMsg(msg string, fields ...zap.Field) {
	debug.DebugMsg(e.debug, fmt.Sprintf("emitter: %s", msg), fields...)
}
