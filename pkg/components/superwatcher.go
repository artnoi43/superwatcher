package components

import (
	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
)

// superWatcher implements superwatcher.SuperWatcher.
// It is a meta-type in that it merely wraps Emitter and Engine,
// and only provides superWatcher.Run as its original method.
type superWatcher struct {
	emitter  superwatcher.Emitter
	engine   superwatcher.Engine
	debugger *debugger.Debugger
}

func NewSuperWatcherOptions(options ...Option) superwatcher.SuperWatcher {
	var conf initConfig
	for _, opt := range options {
		opt(&conf)
	}

	logLevel := gslutils.Max(conf.logLevel, conf.conf.LogLevel)

	poller := NewPoller(
		conf.addresses,
		conf.topics,
		conf.conf.DoReorg || conf.doReorg,
		conf.filterRange,
		conf.ethClient,
		logLevel,
	)

	emitter := NewEmitter(
		conf.conf,
		conf.ethClient,
		conf.getStateDataGateway,
		poller,
		conf.syncChan,
		conf.filterResultChan,
		conf.errChan,
	)

	emitterClient := NewEmitterClient(
		conf.conf,
		conf.syncChan,
		conf.filterResultChan,
		conf.errChan,
	)

	engine := NewEngine(
		emitterClient,
		conf.serviceEngine,
		conf.setStateDataGateway,
		logLevel,
	)

	return NewSuperWatcher(emitter, engine, logLevel)
}

func NewSuperWatcherDefault(
	conf *config.Config,
	ethClient superwatcher.EthClient,
	getStateDataGateway superwatcher.GetStateDataGateway,
	setStateDataGateway superwatcher.SetStateDataGateway,
	serviceEngine superwatcher.ServiceEngine,
	addresses []common.Address,
	topics [][]common.Hash,
) superwatcher.SuperWatcher {
	emitter, engine := NewDefault(
		conf,
		ethClient,
		getStateDataGateway,
		setStateDataGateway,
		serviceEngine,
		addresses,
		topics,
	)

	return NewSuperWatcher(emitter, engine, conf.LogLevel)
}

func NewSuperWatcher(
	emitter superwatcher.Emitter,
	engine superwatcher.Engine,
	logLevel uint8,
) superwatcher.SuperWatcher {
	return &superWatcher{
		emitter:  emitter,
		engine:   engine,
		debugger: debugger.NewDebugger("SuperWatcher", logLevel),
	}
}
