package ensengine

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/logger/debugger"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts/ens/enscontroller"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts/ens/ensregistrar"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines"
)

const (
	eventNameRegistered string = "NameRegistered"
)

var ensEngineEvents = []string{eventNameRegistered}

type ensEngine struct {
	ensRegistrar  contracts.BasicContract
	ensController contracts.BasicContract
	dataGateway   datagateway.RepositoryENS
	debugger      *debugger.Debugger
}

type TestSuiteENS struct {
	Engine   superwatcher.ServiceEngine // *ensEngine
	Routes   map[subengines.SubEngineEnum]map[common.Address][]common.Hash
	Services map[subengines.SubEngineEnum]superwatcher.ServiceEngine
}

func New(
	registrarContract contracts.BasicContract,
	controllerContract contracts.BasicContract,
	dgwENS datagateway.RepositoryENS,
	logLevel uint8,
) superwatcher.ServiceEngine {
	return &ensEngine{
		ensRegistrar:  registrarContract,
		ensController: controllerContract,
		dataGateway:   dgwENS,
		debugger:      debugger.NewDebugger("ensEngine", logLevel),
	}
}

// NewTestSuiteENS returns a convenient struct for injecting into routerengine.routerEngine
func NewTestSuiteENS(dgwENS datagateway.RepositoryENS, logLevel uint8) *TestSuiteENS {
	registrarContract := contracts.NewBasicContract(
		"ENSRegistrar",
		ensregistrar.EnsRegistrarABI,
		"0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85",
		ensEngineEvents...,
	)
	controllerContract := contracts.NewBasicContract(
		"ENSController",
		enscontroller.EnsControllerABI,
		"0x283af0b28c62c092c9727f1ee09c02ca627eb7f5",
		ensEngineEvents...,
	)

	ensEngine := New(registrarContract, controllerContract, dgwENS, logLevel)

	registrarTopics := contracts.CollectEventHashes(registrarContract.ContractEvents)
	controllerTopics := contracts.CollectEventHashes(controllerContract.ContractEvents)

	return &TestSuiteENS{
		Engine: ensEngine,
		Routes: map[subengines.SubEngineEnum]map[common.Address][]common.Hash{
			subengines.SubEngineENS: {
				registrarContract.Address:  registrarTopics,
				controllerContract.Address: controllerTopics,
			},
		},
		Services: map[subengines.SubEngineEnum]superwatcher.ServiceEngine{
			subengines.SubEngineENS: ensEngine,
		},
	}
}
