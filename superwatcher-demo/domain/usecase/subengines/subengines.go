package subengines

import (
	"github.com/artnoi43/superwatcher/lib/logger"
)

type (
	SubEngineEnum uint8
	// DemoKey is used to track various states of various items from different contracts.
	DemoKey interface {
		ForSubEngine() SubEngineEnum
	}
)

func AssertDemoKey(itemKey any) DemoKey {
	demoKey, ok := itemKey.(DemoKey)
	if !ok {
		logger.Panic("type assertion failed - itemKey is not DemoKey")
	}
	return demoKey
}

const (
	SubEngineInvalid SubEngineEnum = iota
	SubEngineUniswapv3Factory
	SubEngineUniswapv3Pool
	SubEngineOneInchLimitOrder
)

func (se SubEngineEnum) String() string {
	switch se {
	case SubEngineInvalid:
		return "SUBENGINE_INVALID"
	case SubEngineUniswapv3Factory:
		return "SUBENGINE_UNISWAPV3POOLFACTORY"
	case SubEngineUniswapv3Pool:
		return "SUBENGINE_UNISWAPV3POOl"
	case SubEngineOneInchLimitOrder:
		return "SUBENGINE_ONEINCHLIMITORDER"
	}

	panic("unhandled demo usecase")
}
