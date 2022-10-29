package entity

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// ENS represents the Ethereum domain names.
// For example, domain "foo.eth" has Name "foo" and TLD "eth"
type ENS struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	TLD   common.Address `json:"tld"`
	Owner common.Address `json:"owner"`

	TTL     uint64    `json:"ttl"`
	Expires time.Time `json:"expires"`
}

func (e *ENS) DomainString() string {
	return fmt.Sprintf("%s.%s", e.TLD, e.Name)
}
