package uniswapv3factoryengine

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/entity"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/lib/logutils"
)

func (e *uniswapv3PoolFactoryEngine) handlePoolCreated(
	pool *entity.Uniswapv3PoolCreated,
) error {
	e.debugger.Debug(1, "got poolCreated, writing to db", zap.Any("pool", pool))
	return e.dataGateway.SetPool(context.Background(), pool)
}

func (e *uniswapv3PoolFactoryEngine) revertPoolCreated(
	pool *entity.Uniswapv3PoolCreated,
) error {
	e.debugger.Debug(1, "reverting poolCreated", zap.Any("pool", pool))

	err := e.dataGateway.DelPool(context.Background(), pool)
	if err != nil {
		return errors.Wrapf(err, "failed to revert pool %s", pool.Address.String())
	}

	return nil
}

// parsePoolCreatedUnpackedMap collects unpacked log.Data into *entity.Uniswapv3PoolCreated.
// Other fields not available in the log byte data is populated elsewhere.
func parsePoolCreatedUnpackedMap(unpacked map[string]interface{}) (*entity.Uniswapv3PoolCreated, error) {
	poolAddr, err := logutils.ExtractFieldFromUnpacked[common.Address](unpacked, "pool")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unpack poolAddr from map %v", unpacked)
	}
	return &entity.Uniswapv3PoolCreated{
		Address: poolAddr,
	}, nil
}
