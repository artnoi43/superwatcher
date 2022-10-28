package uniswapv3factoryengine

import (
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
)

func (e *uniswapv3PoolFactoryEngine) handlePoolCreated(
	pool *entity.Uniswapv3PoolCreated,
) error {
	logger.Info("got poolCreated, writing to db", zap.Any("pool", pool))

	return nil
}

func (e *uniswapv3PoolFactoryEngine) revertPoolCreated(
	pool *entity.Uniswapv3PoolCreated,
) error {
	logger.Info("reverting poolCreated", zap.Any("pool", pool))

	return nil
}
