package uniswapv3factoryengine

import (
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/subengines"
)

type poolFactoryArtifact map[entity.Uniswapv3PoolCreated]uniswapv3PoolFactoryState

func (a poolFactoryArtifact) ForSubEngine() subengines.SubEngineEnum {
	return subengines.SubEngineUniswapv3Factory
}
