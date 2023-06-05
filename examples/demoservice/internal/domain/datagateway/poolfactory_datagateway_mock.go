package datagateway

import (
	"context"
	"fmt"
	"sync"

	"github.com/soyart/gsl"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/soyart/superwatcher"

	"github.com/soyart/superwatcher/examples/demoservice/internal/domain/entity"
)

type MockDataGatewayPoolFactory struct {
	sync.RWMutex

	// m maps pool address to service entity pool
	m map[string]*entity.Uniswapv3PoolCreated

	// writeLogs is used to record all write operations done on mockDataGatewayENS.
	// It is useful in demotest.
	writeLogs []WriteLog
}

func NewMockDataGatewayPoolFactory() RepositoryPoolFactory {
	return &MockDataGatewayPoolFactory{
		m: make(map[string]*entity.Uniswapv3PoolCreated),
	}
}

func (s *MockDataGatewayPoolFactory) SetPool(
	ctx context.Context,
	pool *entity.Uniswapv3PoolCreated,
) error {
	s.Lock()
	defer s.Unlock()

	addr := gsl.StringerToLowerString(pool.Address)
	hash := gsl.StringerToLowerString(pool.BlockHash)

	s.m[addr] = pool
	s.writeLogs = append(
		s.writeLogs,
		WriteLog(
			fmt.Sprintf("SET POOL %s BLOCK %d HASH %s", addr, pool.BlockCreated, hash),
		),
	)

	return nil
}

func (s *MockDataGatewayPoolFactory) GetPool(
	ctx context.Context,
	lpAddress common.Address,
) (
	*entity.Uniswapv3PoolCreated,
	error,
) {
	s.RLock()
	defer s.RUnlock()

	addr := gsl.StringerToLowerString(lpAddress)
	pool, ok := s.m[addr]
	if !ok {
		return nil, errors.Wrapf(superwatcher.ErrRecordNotFound, "lp %s not found", addr)
	}

	return pool, nil
}

func (s *MockDataGatewayPoolFactory) GetPools(ctx context.Context) ([]*entity.Uniswapv3PoolCreated, error) {
	s.RLock()
	defer s.RUnlock()

	var pools []*entity.Uniswapv3PoolCreated //nolint:prealloc
	for _, pool := range s.m {
		pools = append(pools, pool)
	}

	return pools, nil
}

func (s *MockDataGatewayPoolFactory) DelPool(
	ctx context.Context,
	pool *entity.Uniswapv3PoolCreated,
) error {
	s.Lock()
	defer s.Unlock()

	addr := gsl.StringerToLowerString(pool.Address)
	hash := gsl.StringerToLowerString(pool.BlockHash)
	pool, ok := s.m[addr] //nolint:staticcheck
	if !ok {
		return errors.Wrapf(superwatcher.ErrRecordNotFound, "lp %s not found", addr)
	}

	s.m[addr] = nil
	s.writeLogs = append(
		s.writeLogs,
		WriteLog(
			fmt.Sprintf("DEL POOL %s BLOCK %d HASH %s", addr, pool.BlockCreated, hash),
		),
	)

	return nil
}

func (s *MockDataGatewayPoolFactory) WriteLogs() []WriteLog {
	return s.writeLogs
}
