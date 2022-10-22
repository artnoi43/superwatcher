package demoengine

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase"
)

type (
	// demoEngine wraps "subservices' engines"
	demoEngine struct {
		usecases map[common.Address]usecase.UseCase
		services map[usecase.UseCase]engine.ServiceEngine[usecase.DemoKey, engine.ServiceItem[usecase.DemoKey]]

		// fsm is a map[usecase.UseCase]engine.ServiceFSM[usecase.DemoKey].
		// i.e. it wraps subservice FSM, to be returned by *demoEngine.ServiceStateTracker().
		// *engine.Engine calls ServiceStateTracker before entering a loop, so the one returned
		// must have access to all of the subservices' FSMs
		fsm *demoFSM
	}
)

func New(
	usecases map[common.Address]usecase.UseCase,
	services map[usecase.UseCase]engine.ServiceEngine[usecase.DemoKey, engine.ServiceItem[usecase.DemoKey]],
	fsm *demoFSM,
) engine.ServiceEngine[usecase.DemoKey, engine.ServiceItem[usecase.DemoKey]] {
	return &demoEngine{
		usecases: usecases,
		services: services,
		fsm:      fsm,
	}
}
