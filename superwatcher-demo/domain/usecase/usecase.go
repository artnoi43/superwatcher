package usecase

type (
	UseCase uint8
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
