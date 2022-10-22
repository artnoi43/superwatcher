package demoengine

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase"
)

type (
	// demoKey is used to track various states of various items from different contracts.
	demoKey interface {
		engine.ItemKey
		GetUseCase() usecase.UseCase
	}

	// demoEngine wraps "subservices' engines"
	demoEngine struct {
		usecases map[common.Address]usecase.UseCase
		services map[usecase.UseCase]engine.ServiceEngine[demoKey, engine.ServiceItem[demoKey]]

		// fsm is a map[usecase.UseCase]engine.ServiceFSM[DemoKey].
		// i.e. it wraps subservice FSM, to be returned by *demoEngine.ServiceStateTracker().
		// *engine.Engine calls ServiceStateTracker before entering a loop, so the one returned
		// must have access to all of the subservices' FSMs
		fsm *demoFSM
	}
)

func newDemoEngine(
	usecases map[common.Address]usecase.UseCase,
	services map[usecase.UseCase]engine.ServiceEngine[demoKey, engine.ServiceItem[demoKey]],
	fsm *demoFSM,
) engine.ServiceEngine[demoKey, engine.ServiceItem[demoKey]] {
	return &demoEngine{
		usecases: usecases,
		services: services,
		fsm:      fsm,
	}
}
