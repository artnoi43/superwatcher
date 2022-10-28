package ensengine

import (
	"github.com/artnoi43/superwatcher/domain/usecase/engine"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ensEngine struct {
	ensContract contracts.BasicContract
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
