package ensengine

import (
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/artnoi43/superwatcher/pkg/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/contracts"
)

type ensEngine struct {
	ensRegistrar  contracts.BasicContract
	ensController contracts.BasicContract
}

func NewEnsEngine(
	contractABI abi.ABI,
	contractEvents []abi.Event,
) superwatcher.ServiceEngine {
	return nil
	// return &limitOrderEngine{
	// 	poolFactoryContract: contracts.BasicContract{
	// 		ContractABI:    contractABI,
	// 		ContractEvents: contractEvents,
	// 	},
	// }
}
