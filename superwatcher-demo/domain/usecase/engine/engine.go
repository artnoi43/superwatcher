package engine

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/usecase/logutils"
)

type uniswapv3Engine struct {
	mapAddrToABI    map[common.Address]abi.ABI
	mapAddrToTopics map[common.Address][]string
}

func NewUniswapV3Engine(mapTopicABI map[common.Address]abi.ABI) *uniswapv3Engine {
	return &uniswapv3Engine{
		mapAddrToABI: mapTopicABI,
	}
}

func (e *uniswapv3Engine) MapLogToItem(log *types.Log) (*entity.UniswapSwap, error) {
	contractAddr := log.Address
	contractABI, ok := e.mapAddrToABI[contractAddr]
	if !ok {
		return nil, fmt.Errorf("abi not found for address %s", contractAddr.String())
	}

	topicStrings := e.mapAddrToTopics[contractAddr]
	realUnpacked := make(map[string]interface{})
	for _, topicString := range topicStrings {
		unpacked, err := logutils.UnpackIntoMap(contractABI, topicString, log)
		if err != nil {
			continue
		}
		realUnpacked = unpacked
		break
	}

	swapEvent, err := parseLogToUniswapSwap(realUnpacked)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert unpacked log data to *entity.UniswapSwap")
	}

	return swapEvent, nil
}

func parseLogToUniswapSwap(unpacked map[string]interface{}) (*entity.UniswapSwap, error) {
	return nil, errors.New("not implemented")
}
