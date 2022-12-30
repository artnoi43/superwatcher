package datagateway

import (
	"context"
	"fmt"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/superwatcher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/entity"
)

type mockDataGatewayPoolFactory struct {
	m map[string]*entity.Uniswapv3PoolCreated
}

func NewMockDataGatewayPoolFactory() RepositoryPoolFactory {
	return &mockDataGatewayPoolFactory{
		m: make(map[string]*entity.Uniswapv3PoolCreated),
	}
}

func (s *mockDataGatewayPoolFactory) SetPool(
	ctx context.Context,
	pool *entity.Uniswapv3PoolCreated,
) error {
	addr := gslutils.StringerToLowerString(pool.Address)
	s.m[addr] = pool

	return nil
}

func (s *mockDataGatewayPoolFactory) GetPool(
	ctx context.Context,
	lpAddress common.Address,
) (
	*entity.Uniswapv3PoolCreated,
	error,
) {
	addr := gslutils.StringerToLowerString(lpAddress)
	fmt.Println("SET", addr)
	pool, ok := s.m[addr]
	if !ok {
		return nil, errors.Wrapf(superwatcher.ErrRecordNotFound, "lp %s not found", addr)
	}

	return pool, nil
}

func (s *mockDataGatewayPoolFactory) GetPools(ctx context.Context) ([]*entity.Uniswapv3PoolCreated, error) {
	var pools []*entity.Uniswapv3PoolCreated //nolint:prealloc
	for _, pool := range s.m {
		pools = append(pools, pool)
	}

	return pools, nil
}

func (s *mockDataGatewayPoolFactory) DelPool(
	ctx context.Context,
	pool *entity.Uniswapv3PoolCreated,
) error {
	addr := gslutils.StringerToLowerString(pool.Address)
	fmt.Println("DEL", addr)
	pool, ok := s.m[addr] //nolint:staticcheck
	if !ok {
		return errors.Wrapf(superwatcher.ErrRecordNotFound, "lp %s not found", addr)
	}

	s.m[addr] = nil

	return nil
}
