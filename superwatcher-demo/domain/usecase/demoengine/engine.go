package demoengine

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/subengines"
)

type (
	// demoEngine wraps "subservices' engines"
	demoEngine struct {
		usecases map[common.Address]subengines.SubEngine
		services map[subengines.SubEngine]engine.ServiceEngine[subengines.DemoKey, engine.ServiceItem[subengines.DemoKey]]

		// fsm is a map[subengines.UseCase]engine.ServiceFSM[subengines.DemoKey].
		// i.e. it wraps subservice FSM, to be returned by *demoEngine.ServiceStateTracker().
		// *engine.Engine calls ServiceStateTracker before entering a loop, so the one returned
		// must have access to all of the subservices' FSMs
		fsm *demoFSM
	}
)

func New(
	usecases map[common.Address]subengines.SubEngine,
	services map[subengines.SubEngine]engine.ServiceEngine[subengines.DemoKey, engine.ServiceItem[subengines.DemoKey]],
	fsm engine.ServiceFSM[subengines.DemoKey],
) engine.ServiceEngine[subengines.DemoKey, engine.ServiceItem[subengines.DemoKey]] {
	demoFSM, ok := fsm.(*demoFSM)
	if !ok {
		logger.Panic("fsm is not *demoFSM", zap.String("actual type", reflect.TypeOf(fsm).String()))
	}

	return &demoEngine{
		usecases: usecases,
		services: services,
		fsm:      demoFSM,
	}
}
