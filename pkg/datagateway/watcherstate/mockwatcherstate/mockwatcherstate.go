package mockwatcherstate

import (
	"context"

	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate"
)

type fakeRedis struct {
	lastRecordedBlock uint64
}

func New(lastRecordedBlock uint64) watcherstate.StateDataGateway {
	return &fakeRedis{
		lastRecordedBlock: lastRecordedBlock,
	}
}

func (m *fakeRedis) GetLastRecordedBlock(context.Context) (uint64, error) {
	return m.lastRecordedBlock, nil
}

func (m *fakeRedis) SetLastRecordedBlock(ctx context.Context, v uint64) error {
	m.lastRecordedBlock = v
	return nil
}

func (m *fakeRedis) Shutdown() error {
	return nil
}
