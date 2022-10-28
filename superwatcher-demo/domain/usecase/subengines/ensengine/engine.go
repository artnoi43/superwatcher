package ensengine

import (
	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ensEngine struct {
	// https://etherscan.io/address/0x3ef51736315f52d568d6d2cf289419b9cfffe782
	limitOrderContract contracts.BasicContract
}

func NewEnsEngine(
	contractABI abi.ABI,
	contractEvents []abi.Event,
) engine.ServiceEngine {
	return nil
	// return &limitOrderEngine{
	// 	poolFactoryContract: contracts.BasicContract{
	// 		ContractABI:    contractABI,
	// 		ContractEvents: contractEvents,
	// 	},
	// }
}
