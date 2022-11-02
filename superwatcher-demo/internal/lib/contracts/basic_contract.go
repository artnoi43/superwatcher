package contracts

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// BasicContract is a utility struct that provides
// contract-specific information
type BasicContract struct {
	Name           string
	Address        common.Address
	ContractABI    abi.ABI
	ContractEvents []abi.Event
}

func NewBasicContract(name string, contractJsonABI string, addrString string, eventKeys ...string) BasicContract {
	contractABI, err := abi.JSON(strings.NewReader(contractJsonABI))
	if err != nil {
		panic("error parsing contract ABI: " + err.Error())
	}

	contractTopics := accrueEvents(contractABI, eventKeys...)
	basicContract := BasicContract{
		Name:           name,
		Address:        common.HexToAddress(addrString),
		ContractABI:    contractABI,
		ContractEvents: contractTopics,
	}

	return basicContract
}
