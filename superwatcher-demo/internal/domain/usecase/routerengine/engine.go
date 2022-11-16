package routerengine

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/subengines"
)

type (
	// routerEngine wraps "subservices' engines"
	routerEngine struct {
		routes   map[subengines.SubEngineEnum]map[common.Address][]common.Hash
		services map[subengines.SubEngineEnum]superwatcher.ServiceEngine
	}
)

func New(
	routes map[subengines.SubEngineEnum]map[common.Address][]common.Hash,
	services map[subengines.SubEngineEnum]superwatcher.ServiceEngine,
) superwatcher.ServiceEngine {
	return &routerEngine{
		routes:   routes,
		services: services,
	}
}
