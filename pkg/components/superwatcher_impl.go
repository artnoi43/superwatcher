package components

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
)

func (spw *superWatcher) Run(
	ctx context.Context,
	cancel context.CancelFunc,
) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := spw.emitter.Loop(ctx); err != nil {
			cancel()
			return
		}
	}()

	return errors.Wrap(spw.engine.Loop(ctx), "engine.Loop exited")
}

func (spw *superWatcher) Emitter() superwatcher.Emitter {
	return spw.emitter
}

func (spw *superWatcher) Engine() superwatcher.Engine {
	return spw.engine
}

func (spw *superWatcher) Shutdown() {
	spw.emitter.Shutdown()
}

func (spw *superWatcher) SetDoReorg(doReorg bool) {
	spw.emitter.Poller().SetDoReorg(doReorg)
}

func (spw *superWatcher) DoReorg() bool {
	return spw.emitter.Poller().DoReorg()
}

func (spw *superWatcher) Addresses() []common.Address {
	return spw.emitter.Poller().Addresses()
}

func (spw *superWatcher) Topics() [][]common.Hash {
	return spw.emitter.Poller().Topics()
}

func (spw *superWatcher) AddAddresses(addresses ...common.Address) {
	spw.emitter.Poller().AddAddresses(addresses...)
}

func (spw *superWatcher) AddTopics(topics ...[]common.Hash) {
	spw.emitter.Poller().AddTopics(topics...)
}

func (spw *superWatcher) SetAddresses(addresses []common.Address) {
	spw.emitter.Poller().SetAddresses(addresses)
}

func (spw *superWatcher) SetTopics(topics [][]common.Hash) {
	spw.emitter.Poller().SetTopics(topics)
}
