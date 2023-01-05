package routerengine

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/gsl/soyutils"
	"github.com/artnoi43/superwatcher"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/hardcode"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/lib/contracts"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines/ensengine"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines/uniswapv3factoryengine"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

const (
	logsPath            = "../../../test_logs/"
	logsPathENS         = logsPath + "/ens"
	logsPathPoolFactory = logsPath + "/poolfactory"
)

func TestHandleGoodLogs(t *testing.T) {
	ensLogs, err := soyutils.ReadFileJSON[[]types.Log](logsPathENS + "/logs_multi_names.json")
	if err != nil {
		t.Skip("bad or missing ENS logs file:", err.Error())
	}
	poolFactoryLogs, err := soyutils.ReadFileJSON[[]types.Log](logsPathPoolFactory + "/log_poolcreated.json")
	if err != nil {
		t.Skip("bad or missing PoolCreated logs file:", err.Error())
	}

	logLevel := uint8(2)

	demoContracts := hardcode.DemoContracts(hardcode.Uniswapv3Factory, hardcode.ENSRegistrar, hardcode.ENSController)
	poolFactoryContract := demoContracts[hardcode.Uniswapv3Factory]
	poolFactoryHashes := contracts.CollectEventHashes(poolFactoryContract.ContractEvents)
	poolFactoryEngine := uniswapv3factoryengine.New(poolFactoryContract, datagateway.NewMockDataGatewayPoolFactory(), logLevel)
	ensRegistrarContract := demoContracts[hardcode.ENSRegistrar]
	ensRegistrarHashes := contracts.CollectEventHashes(ensRegistrarContract.ContractEvents)
	ensControllerContract := demoContracts[hardcode.ENSController]
	ensControllerHashes := contracts.CollectEventHashes(ensControllerContract.ContractEvents)
	ensEngine := ensengine.New(ensRegistrarContract, ensControllerContract, datagateway.NewMockDataGatewayENS(), logLevel)

	routes := map[subengines.SubEngineEnum]map[common.Address][]common.Hash{
		subengines.SubEngineUniswapv3Factory: {
			common.HexToAddress(hardcode.Uniswapv3FactoryAddr): poolFactoryHashes,
		},
		subengines.SubEngineENS: {
			common.HexToAddress(hardcode.ENSRegistrarAddr):  ensRegistrarHashes,
			common.HexToAddress(hardcode.ENSControllerAddr): ensControllerHashes,
		},
	}

	services := map[subengines.SubEngineEnum]superwatcher.ServiceEngine{
		subengines.SubEngineUniswapv3Factory: poolFactoryEngine,
		subengines.SubEngineENS:              ensEngine,
	}

	routerEngine := New(routes, services, 2)

	logs := append(ensLogs, poolFactoryLogs...)
	mappedLogs := reorgsim.MapLogsToNumber(logs)
	var blocks []*superwatcher.Block

	for number, blockLogs := range mappedLogs {
		if len(blockLogs) == 0 {
			continue
		}

		blocks = append(blocks, &superwatcher.Block{
			Number: number,
			Hash:   blockLogs[0].BlockHash,
			Logs:   gslutils.CollectPointers(blockLogs),
		})
	}

	testHandleGoodBlocks(t, routerEngine, blocks, 2)
}

func testHandleGoodBlocks(
	t *testing.T,
	routerEngine superwatcher.ServiceEngine,
	blocks []*superwatcher.Block,
	numSubEngines int, // Number of subEngines within the router
) {

	// Should have len == 1, since this is just a single call to HandleGoodLogs
	// and the demoEngine only has 1 sub-engine.
	mapArtifacts, err := routerEngine.HandleGoodBlocks(blocks, nil)
	if err != nil {
		t.Errorf("error in demoEngine.HandleGoodLogs: %s", err.Error())
	}

	for blockHash, artifacts := range mapArtifacts {
		t.Logf("len artifacts for blockHash %s: %d\n", blockHash, len(artifacts))

		for i, artifact := range artifacts {
			switch artifact.(type) {
			case ensengine.ENSArtifact:
				err = assertType[ensengine.ENSArtifact](artifact)
			case uniswapv3factoryengine.PoolFactoryArtifact:
				err = assertType[uniswapv3factoryengine.PoolFactoryArtifact](artifact)
			}

			if err != nil {
				t.Fatalf("%d: %s\n", i, err.Error())
			}
		}
	}
}

func assertType[T any](artifact superwatcher.Artifact) error {
	var t T
	if _, ok := artifact.(T); !ok {
		return fmt.Errorf("artifact is not of type %s", reflect.TypeOf(t).String())
	}

	return nil
}
