package routerengine

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"

	"github.com/artnoi43/superwatcher/pkg/logger/debugger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines"
)

type (
	// routerEngine wraps "subservices' engines"
	routerEngine struct {
		routes   map[subengines.SubEngineEnum]map[common.Address][]common.Hash
		services map[subengines.SubEngineEnum]superwatcher.ServiceEngine
		debugger *debugger.Debugger
	}
)

func New(
	routes map[subengines.SubEngineEnum]map[common.Address][]common.Hash,
	services map[subengines.SubEngineEnum]superwatcher.ServiceEngine,
	logLevel uint8,
) superwatcher.ServiceEngine {
	return &routerEngine{
		routes:   routes,
		services: services,
		debugger: debugger.NewDebugger("routerEngine", logLevel),
	}
}
