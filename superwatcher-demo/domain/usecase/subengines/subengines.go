package subengines

import "github.com/artnoi43/superwatcher/domain/usecase/engine"

type (
	UseCase uint8
	// DemoKey is used to track various states of various items from different contracts.
	DemoKey interface {
		engine.ItemKey
		GetUseCase() UseCase
	}
)

const (
	UseCaseInvalid UseCase = iota
	UseCaseUniswapv3Factory
	UseCaseUniswapv3Pool
	UseCaseOneInchLimitOrder
)

func (uc UseCase) String() string {
	switch uc {
	case UseCaseInvalid:
		return "USECASE_INVALID"
	case UseCaseUniswapv3Factory:
		return "USECASE_UNISWAPV3POOLFACTORY"
	case UseCaseUniswapv3Pool:
		return "USECASE_UNISWAPV3POOl"
	case UseCaseOneInchLimitOrder:
		return "USECASE_ONEINCHLIMITORDER"
	}

	panic("unhandled demo usecase")
}
