package routerengine

import (
	"reflect"
	"testing"

	"github.com/artnoi43/gsl/gslutils"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"

	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines/ensengine"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/subengines/uniswapv3factoryengine"
)

func TestRouterArtifacts(t *testing.T) {
	logsPath := "../../../../test_logs/servicetest/logs_servicetest_16054000_16054100.json"
	logs := reorgsim.InitLogsFromFiles(logsPath)
	logsCasted := gslutils.CollectPointers(logs)

	dgwENS := datagateway.NewMockDataGatewayENS()
	dgwPoolFactory := datagateway.NewMockDataGatewayPoolFactory()
	router := NewMockRouter(4, dgwENS, dgwPoolFactory)

	mapArtifacts, err := router.HandleGoodLogs(logsCasted, nil)
	if err != nil {
		t.Error(err.Error())
	}

	if len(mapArtifacts) == 0 {
		t.Fatal("empty artifacts")
	}

	var artifacts []superwatcher.Artifact
	for _, outArtifacts := range mapArtifacts {
		artifacts = append(artifacts, outArtifacts...)
	}

	var artifactsENS []ensengine.ENSArtifact
	var artifactsPoolFactory []uniswapv3factoryengine.PoolFactoryArtifact
	for _, artifact := range artifacts {
		switch artifact.(type) {
		case ensengine.ENSArtifact:
			artifactsENS = append(artifactsENS, artifact.(ensengine.ENSArtifact))
		case uniswapv3factoryengine.PoolFactoryArtifact:
			artifactsPoolFactory = append(artifactsPoolFactory, artifact.(uniswapv3factoryengine.PoolFactoryArtifact))
		default:
			t.Fatalf("unexpected artifact type: %s", reflect.TypeOf(artifact).String())
		}
	}

	if len(artifactsENS) == 0 {
		t.Fatal("0 ENS artifacts")
	}
	if len(artifactsPoolFactory) == 0 {
		t.Fatal("0 PoolFactory artifacts")
	}
}
