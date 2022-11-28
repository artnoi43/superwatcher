package routerengine

import (
	"reflect"
	"testing"

	"github.com/artnoi43/gsl/gslutils"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
)

func TestRouterArtifacts(t *testing.T) {
	logsPath := "../../assets/servicetest/logs_servicetest_16054000_16054100.json"
	logs := reorgsim.InitLogsFromFiles(logsPath)
	logsCasted := gslutils.CollectPointers(logs)

	dgwENS := datagateway.NewMockDataGatewayENS()
	dgwPoolFactory := datagateway.NewMockDataGatewayPoolFactory()
	router := NewMockRouter(4, dgwENS, dgwPoolFactory)

	artifacts, err := router.HandleGoodLogs(logsCasted, nil)
	if err != nil {
		t.Error(err.Error())
	}

	typeArtifacts := reflect.TypeOf(artifacts).String()
	if typeArtifacts != "[]superwatcher.Artifact" {
		t.Error("typeOf(artifacts)", typeArtifacts)
	}
	for _, seArtifacts := range artifacts {
		typeSeArtifacts := reflect.TypeOf(seArtifacts).String()
		if typeSeArtifacts != "[]superwatcher.Artifact" {
			t.Error("typeOf(seArtifacts)", typeSeArtifacts)
		}

		for _, seArtifact := range seArtifacts.([]superwatcher.Artifact) {
			typeSeArtifact := reflect.TypeOf(seArtifact).String()
			if typeSeArtifact == "[]superwatcher.Artifact" {
				t.Error("typeOf(seArtifact)", typeSeArtifact)
			}
		}
	}
}
