package ensengine

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/lib/logger"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/demoengine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/subengines"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts/ens/enscontroller"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts/ens/ensregistrar"
)

func TestHandleLogs(t *testing.T) {
	var logs []*types.Log
	if err := json.Unmarshal([]byte(logsJson), &logs); err != nil {
		t.Errorf("error unmarshaling json logs: %s", err.Error())
	}

	registrarContract := newContract(
		ensregistrar.EnsRegistrarABI,
		"0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85",
		nameRegistered,
	)
	controllerContract := newContract(
		enscontroller.EnsControllerABI,
		"0x283af0b28c62c092c9727f1ee09c02ca627eb7f5",
		nameRegistered,
	)

	ensEngine := &ensEngine{
		ensRegistrar:  registrarContract,
		ensController: controllerContract,
	}

	artifacts, err := ensEngine.HandleGoodLogs(logs)
	if err != nil {
		t.Errorf("HandleGoodLogs error: %s", err.Error())
	}

	// TODO: Assert
	// https://etherscan.io/tx/0x07fff3cd11172e3878900dd22e8e905674651aa5f91f04ff35926150d2db9671#eventlog
	expectedENS := entity.ENS{
		ID:      "0x05768d5da4db7b041a733407418205278f29329dde9153be3247cac968509d14",
		Name:    "onchainalpha",
		Expires: time.Unix(1730090099, 0),
		Owner:   common.HexToAddress("0x8AD703901c3FcDECD20D2B9349F8183d0d14FDDF"),
	}

	for _, artifact := range artifacts {
		ensArtifacts, ok := artifact.([]ENSArtifact)
		if !ok {
			t.Fatal("artifact is not []ENSArtifact")
		}
		for i, ensArtifact := range ensArtifacts {
			if ensArtifact.LastEvent != Registered {
				t.Fatal("invalid last event")
			}
			// The last one should map
			if i == 1 {
				if !reflect.DeepEqual(ensArtifact.ENS, expectedENS) {
					logger.Debug("expected", zap.Any("ens", expectedENS))
					logger.Debug("actual", zap.Any("ens", ensArtifact.ENS))
					t.Fatal("unexpected ENS result")
				}
			}
		}
	}

	demoSubEngines := map[common.Address]subengines.SubEngineEnum{
		registrarContract.Address:  subengines.SubEngineENS,
		controllerContract.Address: subengines.SubEngineENS,
	}
	demoServices := map[subengines.SubEngineEnum]engine.ServiceEngine{
		subengines.SubEngineENS: ensEngine,
	}
	demoEngine := demoengine.New(demoSubEngines, demoServices)

	artifacts, err = demoEngine.HandleGoodLogs(logs)
	if err != nil {
		t.Fatalf("demoEngine.HandleGoodLogs failed: %s\n", err.Error())
	}
	for i, artifact := range artifacts {
		t.Logf("%d %s %+v\n", i, reflect.TypeOf(artifact), artifact)
	}
}

func newContract(contractJsonABI string, addrString string, eventKeys ...string) contracts.BasicContract {
	contractABI, err := abi.JSON(strings.NewReader(contractJsonABI))
	if err != nil {
		panic("invalid json abi")
	}
	contractTopics := accrueEvents(contractABI, eventKeys...)
	basicContract := contracts.BasicContract{
		Address:        common.HexToAddress(addrString),
		ContractABI:    contractABI,
		ContractEvents: contractTopics,
	}

	return basicContract
}

func accrueEvents(contractABI abi.ABI, eventKeys ...string) []abi.Event {
	var events []abi.Event
	for key, event := range contractABI.Events {
		for _, eventKey := range eventKeys {
			if key == eventKey {
				events = append(events, event)
			}
		}
	}

	return events
}

const nameRegistered = "NameRegistered"
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
