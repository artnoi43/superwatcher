package entity

import "github.com/ethereum/go-ethereum/common"

// ENS represents the Ethereum domain names
type ENS struct {
	DomainName string
	Owner      common.Address
	TTL        uint64
}
