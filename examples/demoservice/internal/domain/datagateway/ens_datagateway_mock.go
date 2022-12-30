package datagateway

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/entity"
)

type MockDataGatewayENS struct {
	sync.RWMutex
	// m maps ID to *mockENS
	m map[string]*entity.ENS

	// WriteLogs is used to record all write operations done on mockDataGatewayENS.
	// It is useful in demotest.
	WriteLogs []WriteLog
}

func NewMockDataGatewayENS() RepositoryENS {
	return &MockDataGatewayENS{
		m: make(map[string]*entity.ENS),
	}
}

func (s *MockDataGatewayENS) SetENS(ctx context.Context, ens *entity.ENS) error {
	s.Lock()
	defer s.Unlock()

	writeLog := WriteLog(fmt.Sprintf("DEL ID %s BLOCK %d HASH %s", ens.ID, ens.BlockNumber, ens.BlockHash))
	log.Println(writeLog)

	s.m[ens.ID] = ens
	s.WriteLogs = append(s.WriteLogs, writeLog)

	return nil
}

func (s *MockDataGatewayENS) GetENS(ctx context.Context, ensID string) (*entity.ENS, error) {
	s.RLock()
	defer s.RUnlock()

	saved, ok := s.m[ensID]
	if !ok || saved == nil {
		return nil, errors.Wrapf(superwatcher.ErrRecordNotFound, "ens not found for key %s", ensID)
	}

	return saved, nil
}

func (s *MockDataGatewayENS) GetENSes(context.Context) ([]*entity.ENS, error) {
	s.RLock()
	defer s.RUnlock()

	var enses []*entity.ENS //nolint:prealloc
	for _, saved := range s.m {
		enses = append(enses, saved)
	}

	return enses, nil
}

func (s *MockDataGatewayENS) DelENS(ctx context.Context, ens *entity.ENS) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.m[ens.ID]
	if !ok {
		return errors.Wrapf(superwatcher.ErrRecordNotFound, "ens not found for key %s", ens.Name)
	}

	writeLog := WriteLog(fmt.Sprintf("DEL ID %s BLOCK %d HASH %s", ens.ID, ens.BlockNumber, ens.BlockHash))
	log.Println(writeLog)

	s.m[ens.ID] = nil
	s.WriteLogs = append(s.WriteLogs, writeLog)

	return nil
}
