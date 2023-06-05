package mock

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/soyart/superwatcher"
)

type fakeRedisMem struct {
	sync.RWMutex

	lastRecordedBlock uint64
	ok                bool
}

func (m *fakeRedisMem) GetLastRecordedBlock(ctx context.Context) (uint64, error) {
	m.RLock()
	defer m.RUnlock()

	if m.ok {
		return m.lastRecordedBlock, nil
	}

	return 0, errors.Wrap(superwatcher.ErrRecordNotFound, "key not found")
}

func (m *fakeRedisMem) SetLastRecordedBlock(ctx context.Context, v uint64) error {
	m.Lock()
	defer m.Unlock()

	m.ok = true
	m.lastRecordedBlock = v

	return nil
}

func (m *fakeRedisMem) Shutdown() error {
	return nil
}
