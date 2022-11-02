package demoengine

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/pkg/superwatcher"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/usecase/subengines/ensengine"
)

const (
	subEnginePath                     = "../subengines"
	ensSubEnginePath                  = subEnginePath + "/ensengine"
	uniswapv3PoolFactorySubEnginePath = subEnginePath + "/uniswapv3factoryengine"
)

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

func TestSubEngineENS(t *testing.T) {
	ensLogs, err := readSampleLogs(ensSubEnginePath + "/sample_logs.json")
	if err != nil {
		t.Errorf("failed to read ENS sample logs: %s", err.Error())
	}

	ensSubEngineBundle := ensengine.NewEnsSubEngine()
	demoEngine := New(ensSubEngineBundle.DemoSubEngines, ensSubEngineBundle.DemoServices)

	// Should have len == 1, since this is just a single call to HandleGoodLogs
	// and the demoEngine only has 1 sub-engine.
	artifacts, err := demoEngine.HandleGoodLogs(ensLogs, nil)
	if err != nil {
		t.Errorf("error in demoEngine.HandleGoodLogs: %s", err.Error())
	}
	t.Logf("len artifacts: %d\n", len(artifacts))
	if l := len(artifacts); l != 1 {
		t.Errorf("unexpected artifacts len - expected 1, got %d", l)
	}

	for _, artifact := range artifacts {
		engineArtifacts := artifact.([]superwatcher.Artifact)
		for i, engineArtifact := range engineArtifacts {
			ensArtifact, ok := engineArtifact.(ensengine.ENSArtifact)
			if !ok {
				t.Logf(
					"%d seArtifact returned is not subEngineArtifact: %s\n",
					i, reflect.TypeOf(engineArtifact).String(),
				)
			}
			t.Log(ensArtifact)
		}
	}
}
