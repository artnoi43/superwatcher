package entity

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// ENS represents the Ethereum domain names.
// For example, domain "foo.eth" has Name "foo" and TLD "eth"
type ENS struct {
	Name  string
	TLD   common.Address
	Owner common.Address

	TTL uint64
}

func (e *ENS) DomainString() string {
	return fmt.Sprintf("%s.%s", e.TLD, e.Name)
}
