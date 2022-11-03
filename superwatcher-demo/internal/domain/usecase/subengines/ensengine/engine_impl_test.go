package ensengine

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
	"github.com/artnoi43/superwatcher/pkg/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
)

func TestHandleLogs(t *testing.T) {
	var logs []*types.Log
	if err := json.Unmarshal([]byte(logsJson), &logs); err != nil {
		t.Errorf("error unmarshaling json logs: %s", err.Error())
	}

	bundle := NewEnsSubEngineSuite()
	ensEngine := bundle.Engine

	var artifacts []superwatcher.Artifact
	artifacts, err := ensEngine.HandleGoodLogs(logs, artifacts)
	if err != nil {
		t.Errorf("HandleGoodLogs error: %s", err.Error())
	}
	for _, a := range artifacts {
		t.Log(a)
	}

	// TODO: Assert
	// https://etherscan.io/tx/0x07fff3cd11172e3878900dd22e8e905674651aa5f91f04ff35926150d2db9671#eventlog
	expectedENS := entity.ENS{
		ID:      "0x05768d5da4db7b041a733407418205278f29329dde9153be3247cac968509d14",
		Name:    "onchainalpha",
		Expires: time.Unix(1730090099, 0),
		Owner:   common.HexToAddress("0x8AD703901c3FcDECD20D2B9349F8183d0d14FDDF"),
	}

	for i, artifact := range artifacts {
		ensArtifact, ok := artifact.(ENSArtifact)
		if !ok {
			t.Fatalf("artifact is not ENSArtifact: %s", reflect.TypeOf(artifact).String())
		}

		switch i {
		case 0:
			if ensArtifact.LastEvent != RegisteredRegistrar {
				t.Fatalf("unexpected last event from log %d\n", i)
			}
		case 1:
			if ensArtifact.LastEvent != RegisteredController {
				t.Fatalf("unexpected last event from log %d\n", i)
			}
			if !reflect.DeepEqual(ensArtifact.ENS, expectedENS) {
				logger.Debug("expected", zap.Any("ens", expectedENS))
				logger.Debug("actual", zap.Any("ens", ensArtifact.ENS))
				t.Fatal("unexpected ENS result\n")
			}
		}
	}
}

const logsJson = `[{
  "address": "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85",
  "topics": [
    "0xb3d987963d01b2f68493b4bdb130988f157ea43070d4ad840fee0466ed9370d9",
    "0x05768d5da4db7b041a733407418205278f29329dde9153be3247cac968509d14",
    "0x0000000000000000000000008ad703901c3fcdecd20d2b9349f8183d0d14fddf"
  ],
  "data": "0x00000000000000000000000000000000000000000000000000000000671f1473",
  "blockNumber": "0xf1d1d6",
  "transactionHash": "0x07fff3cd11172e3878900dd22e8e905674651aa5f91f04ff35926150d2db9671",
  "transactionIndex": "0xc5",
  "blockHash": "0xded84b4fda57886404b68129be4141db6e4dcd95a1b298049f38ed398e676619",
  "logIndex": "0x229",
  "removed": false
},
{
  "address": "0x283af0b28c62c092c9727f1ee09c02ca627eb7f5",
  "topics": [
    "0xca6abbe9d7f11422cb6ca7629fbf6fe9efb1c621f71ce8f02b9f2a230097404f",
    "0x05768d5da4db7b041a733407418205278f29329dde9153be3247cac968509d14",
    "0x0000000000000000000000008ad703901c3fcdecd20d2b9349f8183d0d14fddf"
  ],
  "data": "0x00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000016b939d5cd42a900000000000000000000000000000000000000000000000000000000671f1473000000000000000000000000000000000000000000000000000000000000000c6f6e636861696e616c7068610000000000000000000000000000000000000000",
  "blockNumber": "0xf1d1d6",
  "transactionHash": "0x07fff3cd11172e3878900dd22e8e905674651aa5f91f04ff35926150d2db9671",
  "transactionIndex": "0xc5",
  "blockHash": "0xded84b4fda57886404b68129be4141db6e4dcd95a1b298049f38ed398e676619",
  "logIndex": "0x22a",
  "removed": false
}]`
