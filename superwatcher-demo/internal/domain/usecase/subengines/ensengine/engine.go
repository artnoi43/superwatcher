package ensengine

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/pkg/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/subengines"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts/ens/enscontroller"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts/ens/ensregistrar"
)

type ensEngine struct {
	ensRegistrar  contracts.BasicContract
	ensController contracts.BasicContract
}

type EnsSubEngineSuite struct {
	Engine       superwatcher.ServiceEngine // *ensEngine
	DemoRoutes   map[subengines.SubEngineEnum][]common.Address
	DemoServices map[subengines.SubEngineEnum]superwatcher.ServiceEngine
}

func New(registrarContract, controllerContract contracts.BasicContract) superwatcher.ServiceEngine {
	return &ensEngine{
		ensRegistrar:  registrarContract,
		ensController: controllerContract,
	}
}

// NewEnsSubEngineSuite returns a convenient struct for injecting into demoengine.demoEngine
func NewEnsSubEngineSuite() *EnsSubEngineSuite {
	registrarContract := contracts.NewBasicContract(
		"ENSRegistrar",
		ensregistrar.EnsRegistrarABI,
		"0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85",
		nameRegistered,
	)
	controllerContract := contracts.NewBasicContract(
		"ENSController",
		enscontroller.EnsControllerABI,
		"0x283af0b28c62c092c9727f1ee09c02ca627eb7f5",
		nameRegistered,
	)
	ensEngine := &ensEngine{
		ensRegistrar:  registrarContract,
		ensController: controllerContract,
	}

	return &EnsSubEngineSuite{
		Engine: ensEngine,
		DemoRoutes: map[subengines.SubEngineEnum][]common.Address{
			subengines.SubEngineENS: {
				registrarContract.Address,
				controllerContract.Address,
			},
		},
		DemoServices: map[subengines.SubEngineEnum]superwatcher.ServiceEngine{
			subengines.SubEngineENS: ensEngine,
		},
	}
}

func NewEnsEngine(
	registrarContract contracts.BasicContract,
	controllerContract contracts.BasicContract,
) superwatcher.ServiceEngine {
	return &ensEngine{
		ensRegistrar:  registrarContract,
		ensController: controllerContract,
	}
}
