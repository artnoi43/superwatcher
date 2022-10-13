package engine

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
)

type fooEngine struct {
	mapTopicABI map[common.Hash]abi.ABI
}

func NewFooEngine(mapTopicABI map[common.Hash]abi.ABI) *fooEngine {
	return &fooEngine{
		mapTopicABI: mapTopicABI,
	}
}

func (e *fooEngine) MapLogToItem(log *types.Log) (*entity.UniswapSwap, error) {
	return nil, errors.New("not implemented")
}
