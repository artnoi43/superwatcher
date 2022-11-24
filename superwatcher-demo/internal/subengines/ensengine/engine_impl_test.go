package ensengine

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/utils"
)

func TestHandleENSLogs(t *testing.T) {
	testHandleENSLogs(t, singleNameLogsJSON, "singleNameLogsJSON")
	testHandleENSLogs(t, multiNamesLogsJSON, "multiNamesLogsJSON")
}

func testHandleENSLogs(t *testing.T, logsJSON, logsName string) {
	var logs []*types.Log
	if err := json.Unmarshal([]byte(logsJSON), &logs); err != nil {
		t.Fatalf("error unmarshaling json logs: %s", err.Error())
	}

	bundle := NewEnsSubEngineSuite(datagateway.NewMockDataGatewayENS())
	ensEngine := bundle.Engine
	ensNamesExpected := expecteds[logsJSON]

	var artifacts []superwatcher.Artifact
	artifacts, err := ensEngine.HandleGoodLogs(logs, artifacts)
	if err != nil {
		t.Errorf("HandleGoodLogs error: %s", err.Error())
	}

	var c int
	for i, artifact := range artifacts {
		ensArtifact, ok := artifact.(ENSArtifact)
		if !ok {
			t.Fatalf("artifact is not ENSArtifact: %s", reflect.TypeOf(artifact).String())
		}

		// Each ENS name results in 2 ENSArtifacts. The 1st one is from the Registrar contract (missing domain name),
		// while the 2nd one is from the Controller contract (has domain name).
		// We will be checking the expected value with the 2nd artifacts of that name only.
		if i%2 == 0 {
			continue
		}

		expected := ensNamesExpected[c]
		c++

		if !reflect.DeepEqual(expected, ensArtifact.ENS) {
			t.Logf("expected: %s", utils.StringJSON(expected))
			t.Logf("actual: %s", utils.StringJSON(ensArtifact.ENS))
			t.Fatalf("unexpected ENS result")
		}
	}
}

func TestCountArtifacts(t *testing.T) {
	checkLen := func(logsJSON, logsName string) {
		var logs []*types.Log
		if err := json.Unmarshal([]byte(logsJSON), &logs); err != nil {
			t.Fatalf("error unmarshaling %s: %s", logsName, err.Error())
		}

		bundle := NewEnsSubEngineSuite(datagateway.NewMockDataGatewayENS())
		ensEngine := bundle.Engine

		artifacts, err := ensEngine.HandleGoodLogs(logs, []superwatcher.Artifact{})
		if err != nil {
			t.Errorf("failed to process %s", logsName)
		}

		t.Logf("len artifacts for %s: %d", logsName, len(artifacts))
		t.Logf("artifacts for %s: %s", logsName, utils.StringJSON(artifacts))
	}

	checkLen(singleNameLogsJSON, "singleName")
	checkLen(multiNamesLogsJSON, "multiNames")
}

var (
	expecteds = map[string][]entity.ENS{
		singleNameLogsJSON: {
			{
				ID:               "0x05768d5da4db7b041a733407418205278f29329dde9153be3247cac968509d14",
				Name:             "onchainalpha",
				Expires:          time.Unix(1730090099, 0),
				Owner:            "0x8ad703901c3fcdecd20d2b9349f8183d0d14fddf",
				TxHash:           "0x07fff3cd11172e3878900dd22e8e905674651aa5f91f04ff35926150d2db9671",
				BlockHashCreated: "0xded84b4fda57886404b68129be4141db6e4dcd95a1b298049f38ed398e676619",
			},
		},

		multiNamesLogsJSON: {
			{
				ID:               "0x7dfae123aedc1f4a53cc5bccd1c061c277a7653c5c739cf23289a09d729794d9",
				Name:             "0xpetrichor",
				Expires:          time.Unix(1763290007, 0),
				Owner:            "0xfe4030c3828c9f791d83fdb562c306cca26e22cd",
				TxHash:           "0x91bfd59c93dd026a1b20d145955a35dc86e22837778e3b83035d42d7a53222f5",
				BlockHashCreated: "0xaa7a7c33bbf9417cebdd5fab539a2bb0009d515708761972dcb06abcdee79004",
			},
			{
				ID:               "0x227678a3504e7fad9b65de53daa1f512ad86c42af3ed4aa702e5eade43ef1caa",
				Name:             "botbap",
				Expires:          time.Unix(1731733067, 0),
				Owner:            "0xe8657a903d511e1841c3d383fccdefe567814773",
				TxHash:           "0x1e3b99db4d5b102c609479fa38106ad0af4c0237f851f0d84a55cd9af40e4e84",
				BlockHashCreated: "0x270ce3c16d779ee2670a8676bcb70499e78d07e626b52ca3c0da2bee607981fa",
			},
		},
	}
)

const (
	singleNameLogsJSON = `
[
  {
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
  }
]`

	multiNamesLogsJSON = `
[
  {
    "address": "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85",
    "topics": [
      "0xb3d987963d01b2f68493b4bdb130988f157ea43070d4ad840fee0466ed9370d9",
      "0x7dfae123aedc1f4a53cc5bccd1c061c277a7653c5c739cf23289a09d729794d9",
      "0x000000000000000000000000283af0b28c62c092c9727f1ee09c02ca627eb7f5"
    ],
    "data": "0x000000000000000000000000000000000000000000000000000000006919ab97",
    "blockNumber": "0xf3e59b",
    "transactionHash": "0x91bfd59c93dd026a1b20d145955a35dc86e22837778e3b83035d42d7a53222f5",
    "transactionIndex": "0x7e",
    "blockHash": "0xaa7a7c33bbf9417cebdd5fab539a2bb0009d515708761972dcb06abcdee79004",
    "logIndex": "0xed",
    "removed": false
  },
  {
    "address": "0x283af0b28c62c092c9727f1ee09c02ca627eb7f5",
    "topics": [
      "0xca6abbe9d7f11422cb6ca7629fbf6fe9efb1c621f71ce8f02b9f2a230097404f",
      "0x7dfae123aedc1f4a53cc5bccd1c061c277a7653c5c739cf23289a09d729794d9",
      "0x000000000000000000000000fe4030c3828c9f791d83fdb562c306cca26e22cd"
    ],
    "data": "0x0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000002c16b55b7cedcf000000000000000000000000000000000000000000000000000000006919ab97000000000000000000000000000000000000000000000000000000000000000b3078706574726963686f72000000000000000000000000000000000000000000",
    "blockNumber": "0xf3e59b",
    "transactionHash": "0x91bfd59c93dd026a1b20d145955a35dc86e22837778e3b83035d42d7a53222f5",
    "transactionIndex": "0x7e",
    "blockHash": "0xaa7a7c33bbf9417cebdd5fab539a2bb0009d515708761972dcb06abcdee79004",
    "logIndex": "0xf3",
    "removed": false
  },
  {
    "address": "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85",
    "topics": [
      "0xb3d987963d01b2f68493b4bdb130988f157ea43070d4ad840fee0466ed9370d9",
      "0x227678a3504e7fad9b65de53daa1f512ad86c42af3ed4aa702e5eade43ef1caa",
      "0x000000000000000000000000283af0b28c62c092c9727f1ee09c02ca627eb7f5"
    ],
    "data": "0x000000000000000000000000000000000000000000000000000000006738264b",
    "blockNumber": "0xf3e59c",
    "transactionHash": "0x1e3b99db4d5b102c609479fa38106ad0af4c0237f851f0d84a55cd9af40e4e84",
    "transactionIndex": "0x86",
    "blockHash": "0x270ce3c16d779ee2670a8676bcb70499e78d07e626b52ca3c0da2bee607981fa",
    "logIndex": "0xeb",
    "removed": false
  },
  {
    "address": "0x283af0b28c62c092c9727f1ee09c02ca627eb7f5",
    "topics": [
      "0xca6abbe9d7f11422cb6ca7629fbf6fe9efb1c621f71ce8f02b9f2a230097404f",
      "0x227678a3504e7fad9b65de53daa1f512ad86c42af3ed4aa702e5eade43ef1caa",
      "0x000000000000000000000000e8657a903d511e1841c3d383fccdefe567814773"
    ],
    "data": "0x0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000001d6478e7a89e8a000000000000000000000000000000000000000000000000000000006738264b0000000000000000000000000000000000000000000000000000000000000006626f746261700000000000000000000000000000000000000000000000000000",
    "blockNumber": "0xf3e59c",
    "transactionHash": "0x1e3b99db4d5b102c609479fa38106ad0af4c0237f851f0d84a55cd9af40e4e84",
    "transactionIndex": "0x86",
    "blockHash": "0x270ce3c16d779ee2670a8676bcb70499e78d07e626b52ca3c0da2bee607981fa",
    "logIndex": "0xf1",
    "removed": false
  }
]`
)
