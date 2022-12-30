package mock

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
)

type fakeRedisFile struct {
	sync.Mutex

	filename string
	ok       bool
}

func (m *fakeRedisFile) GetLastRecordedBlock(ctx context.Context) (uint64, error) {
	m.Lock()
	defer m.Unlock()

	if !m.ok {
		return 0, errors.Wrap(superwatcher.ErrRecordNotFound, "key not found")
	}

	b, err := os.ReadFile(m.filename)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get lastRecordedBlock from file %s", m.filename)
	}

	lastRecordedBlock, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse file %s content %s to uint64", m.filename, b)
	}

	return lastRecordedBlock, nil
}

func (m *fakeRedisFile) SetLastRecordedBlock(ctx context.Context, lastRecordedBlock uint64) error {
	m.Lock()
	defer m.Unlock()

	m.ok = true
	return writeLastRecordedBlockToFile(m.filename, lastRecordedBlock)
}

func (m *fakeRedisFile) Shutdown() error {
	return nil
}

func writeLastRecordedBlockToFile(filename string, lastRecordedBlock uint64) error {
	s := fmt.Sprintf("%d", lastRecordedBlock)

	return errors.Wrap(os.WriteFile(filename, []byte(s), os.ModePerm), "failed to write fakeRedisFile db")
}
