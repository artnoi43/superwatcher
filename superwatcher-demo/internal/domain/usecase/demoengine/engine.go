package demoengine

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/pkg/superwatcher"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/subengines"
)

type (
	// demoEngine wraps "subservices' engines"
	demoEngine struct {
		usecases map[common.Address]subengines.SubEngineEnum
		services map[subengines.SubEngineEnum]superwatcher.ServiceEngine
	}
)

func New(
	usecases map[common.Address]subengines.SubEngineEnum,
	services map[subengines.SubEngineEnum]superwatcher.ServiceEngine,
) superwatcher.ServiceEngine {
	return &demoEngine{
		usecases: usecases,
		services: services,
	}
}
