package subengines

import "github.com/artnoi43/superwatcher/domain/usecase/engine"

type (
	SubEngine uint8
	// DemoKey is used to track various states of various items from different contracts.
	DemoKey interface {
		engine.ItemKey
		ForSubEngine() SubEngine
	}
)

const (
	SubEngineInvalid SubEngine = iota
	SubEngineUniswapv3Factory
	SubEngineUniswapv3Pool
	SubEngineOneInchLimitOrder
)

func (se SubEngine) String() string {
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
