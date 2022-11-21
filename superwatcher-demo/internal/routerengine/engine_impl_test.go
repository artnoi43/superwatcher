package routerengine

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/hardcode"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/ensengine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/subengines/uniswapv3factoryengine"
)

const (
	assetsPath            = "../assets"
	assetsPathENS         = assetsPath + "/ens"
	assetsPathPoolFactory = assetsPath + "/poolfactory"
)

func TestHandleGoodLogs(t *testing.T) {
	ensLogs, err := readSampleLogs(assetsPathENS + "/logs_multi_names.json")
	if err != nil {
		t.Skip("bad or missing ENS logs file:", err.Error())
	}
	poolFactoryLogs, err := readSampleLogs(assetsPathPoolFactory + "/log_poolcreated.json")
	if err != nil {
		t.Skip("bad or missing PoolCreated logs file:", err.Error())
	}

	demoContracts := hardcode.DemoContracts(hardcode.Uniswapv3Factory, hardcode.ENSRegistrar, hardcode.ENSController)
	poolFactoryContract := demoContracts[hardcode.Uniswapv3Factory]
	poolFactoryHashes := contracts.CollectEventHashes(poolFactoryContract.ContractEvents)
	poolFactoryEngine := uniswapv3factoryengine.New(poolFactoryContract)
	ensRegistrarContract := demoContracts[hardcode.ENSRegistrar]
	ensRegistrarHashes := contracts.CollectEventHashes(ensRegistrarContract.ContractEvents)
	ensControllerContract := demoContracts[hardcode.ENSController]
	ensControllerHashes := contracts.CollectEventHashes(ensControllerContract.ContractEvents)
	ensEngine := ensengine.New(ensRegistrarContract, ensControllerContract)

	routes := map[subengines.SubEngineEnum]map[common.Address][]common.Hash{
		subengines.SubEngineUniswapv3Factory: {
			common.HexToAddress(hardcode.Uniswapv3FactoryAddr): poolFactoryHashes,
		},
		subengines.SubEngineENS: {
			common.HexToAddress(hardcode.ENSRegistrar):      ensRegistrarHashes,
			common.HexToAddress(hardcode.ENSControllerAddr): ensControllerHashes,
		},
	}

	services := map[subengines.SubEngineEnum]superwatcher.ServiceEngine{
		subengines.SubEngineUniswapv3Factory: poolFactoryEngine,
		subengines.SubEngineENS:              ensEngine,
	}

	routerEngine := New(routes, services)

	logs := append(ensLogs, poolFactoryLogs...)
	testHandleGoodLogs(t, routerEngine, logs, 2)
}

func testHandleGoodLogs(
	t *testing.T,
	routerEngine superwatcher.ServiceEngine,
	logs []*types.Log,
	numSubEngines int, // Number of subEngines within the router
) {

	// Should have len == 1, since this is just a single call to HandleGoodLogs
	// and the demoEngine only has 1 sub-engine.
	artifacts, err := routerEngine.HandleGoodLogs(logs, nil)
	if err != nil {
		t.Errorf("error in demoEngine.HandleGoodLogs: %s", err.Error())
	}
	t.Logf("len artifacts: %d\n", len(artifacts))
	if l := len(artifacts); l != numSubEngines {
		t.Errorf("unexpected artifacts len - expected 1, got %d", l)
	}

	for i, artifact := range artifacts {
		engineArtifacts := artifact.([]superwatcher.Artifact)
		for j, engineArtifact := range engineArtifacts {
			switch engineArtifact.(type) {
			case ensengine.ENSArtifact:
				err = assertType[ensengine.ENSArtifact](engineArtifact)
			case uniswapv3factoryengine.PoolFactoryArtifact:
				err = assertType[uniswapv3factoryengine.PoolFactoryArtifact](engineArtifact)
			}

			if err != nil {
				t.Fatalf("%d - %d: %s\n", i, j, err.Error())
			}
		}
	}
}

func readSampleLogs(filename string) ([]*types.Log, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read log file %s", filename)
	}

	var logs []*types.Log
	if err := json.Unmarshal(fileBytes, &logs); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal log from file %s", filename)
	}

	return logs, nil
}

func assertType[T any](artifact superwatcher.Artifact) error {
	var t T
	if _, ok := artifact.(T); !ok {
		return fmt.Errorf("artifact is not of type %s", reflect.TypeOf(t).String())
	}

	return nil
}
