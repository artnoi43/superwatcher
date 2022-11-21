package uniswapv3factoryengine

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts/uniswapv3factory"
)

func TestMapLogToItem(t *testing.T) {
	// The JSON was copied from this log:
	// https://etherscan.io/tx/0x3431dc2e3b6fd996e9d7672b6cd71eaae33394f03539e285f599bf3275da61f2#eventlog

	logFile := "../../assets/poolfactory/log_poolcreated.json"
	logJsonBytes, err := os.ReadFile(logFile)
	if err != nil {
		t.Errorf("failed to read poolCreated.json: %s", err.Error())
	}

	var logs []*types.Log
	if err := json.Unmarshal(logJsonBytes, &logs); err != nil {
		t.Errorf("failed to unmarshal %s: %s", logFile, err.Error())
	}

	uniswapv3factoryABI, err := abi.JSON(strings.NewReader(uniswapv3factory.Uniswapv3FactoryABI))
	if err != nil {
		t.Errorf("failed to parse contract ABI: %s", err.Error())
	}
	poolCreated, err := mapLogToPoolCreated(uniswapv3factoryABI, "PoolCreated", logs[0])
	if err != nil {
		t.Fatalf("mapLogToItem failed: %s", err.Error())
	}

	expected := entity.Uniswapv3PoolCreated{
		Address:      common.HexToAddress("0x2555E089B5EDceF0457533cDdAC12af27CE3926a"),
		Token0:       common.HexToAddress("0x4b13006980aCB09645131b91D259eaA111eaF5Ba"),
		Token1:       common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
		Fee:          500,
		BlockCreated: 15766355,
	}

	if !reflect.DeepEqual(*poolCreated, expected) {
		t.Logf("expected: %+v\n", expected)
		t.Logf("actual: %+v\n", poolCreated)
		t.Fatal("unexpected result")
	}
}
