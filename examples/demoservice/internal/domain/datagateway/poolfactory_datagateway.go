package datagateway

import (
	"context"
	"encoding/json"

	"github.com/artnoi43/gsl"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/entity"
)

const PoolFactoryRedisKey = "demo:poolfactory"

type RepositoryPoolFactory interface {
	SetPool(context.Context, *entity.Uniswapv3PoolCreated) error
	GetPool(context.Context, common.Address) (*entity.Uniswapv3PoolCreated, error)
	GetPools(context.Context) ([]*entity.Uniswapv3PoolCreated, error)
	DelPool(context.Context, *entity.Uniswapv3PoolCreated) error
}

type repositoryPoolFactory struct {
	redisCli *redis.Client
}

func NewDataGatewayPoolFactory(redisCli *redis.Client) RepositoryPoolFactory {
	return &repositoryPoolFactory{
		redisCli: redisCli,
	}
}

func (s *repositoryPoolFactory) SetPool(
	ctx context.Context,
	pool *entity.Uniswapv3PoolCreated,
) error {
	poolJSON, err := json.Marshal(pool)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal pool %s", pool.Address.String())
	}

	addr := gsl.StringerToLowerString(pool.Address)
	if err := s.redisCli.HSet(ctx, PoolFactoryRedisKey, addr, poolJSON).Err(); err != nil {
		return handleRedisErr(err, "HSet pool", addr)
	}

	return nil
}

func (s *repositoryPoolFactory) GetPool(
	ctx context.Context,
	lpAddress common.Address,
) (
	*entity.Uniswapv3PoolCreated,
	error,
) {
	addr := gsl.StringerToLowerString(lpAddress)
	poolJSON, err := s.redisCli.HGet(ctx, PoolFactoryRedisKey, addr).Result()
	if err != nil {
		return nil, handleRedisErr(err, "HGET pool", addr)
	}

	var pool entity.Uniswapv3PoolCreated
	if err := json.Unmarshal([]byte(poolJSON), &pool); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal poolJSON")
	}

	return &pool, nil
}

func (s *repositoryPoolFactory) GetPools(
	ctx context.Context,
) (
	[]*entity.Uniswapv3PoolCreated,
	error,
) {
	resultMap, err := s.redisCli.HGetAll(ctx, PoolFactoryRedisKey).Result()
	if err != nil {
		return nil, handleRedisErr(err, "HGETALL pool", "null")
	}

	var pools []*entity.Uniswapv3PoolCreated //nolint:prealloc
	for lpAddress, poolJSON := range resultMap {
		var pool entity.Uniswapv3PoolCreated
		if err := json.Unmarshal([]byte(poolJSON), &pool); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal poolJSON %s", lpAddress)
		}

		pools = append(pools, &pool)
	}

	return pools, nil
}

func (s *repositoryPoolFactory) DelPool(
	ctx context.Context,
	pool *entity.Uniswapv3PoolCreated,
) error {
	addr := gsl.StringerToLowerString(pool.Address)
	if err := s.redisCli.HDel(ctx, PoolFactoryRedisKey, addr).Err(); err != nil {
		return handleRedisErr(err, "HDEL pool", addr)
	}

	return nil
}
