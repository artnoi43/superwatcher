package demoengine

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/subengines"
)

type (
	// demoEngine wraps "subservices' engines"
	demoEngine struct {
		usecases map[common.Address]subengines.SubEngineEnum
		services map[subengines.SubEngineEnum]engine.ServiceEngine
	}
)

func New(
	usecases map[common.Address]subengines.SubEngineEnum,
	services map[subengines.SubEngineEnum]engine.ServiceEngine,
) engine.ServiceEngine {
	return &demoEngine{
		usecases: usecases,
		services: services,
	}
}
