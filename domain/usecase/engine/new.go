package engine

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher/domain/usecase/emitter/reorg"
)

func New[K ItemKey, T ServiceItem[K]](
	serviceEngine ServiceEngine[K, T],
	// TODO: For prod, should we create chans inside this func instead?
	logChan chan *types.Log,
	blockChan chan *reorg.BlockInfo,
	reorgChan chan *reorg.BlockInfo,
	errChan chan error,
	debug bool,
) WatcherEngine {
	emitterClient := NewEmitterClientDebug[T](
		logChan,
		blockChan,
		reorgChan,
		errChan,
	)

	return newWatcherEngine(
		emitterClient,
		serviceEngine,
		debug,
	)
}
