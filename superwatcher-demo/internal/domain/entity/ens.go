package entity

import (
	"fmt"
	"time"
)

// ENS represents the Ethereum domain names.
// For example, domain "foo.eth" has Name "foo" and TLD "eth"
type ENS struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	TLD     string    `json:"tld"`
	Owner   string    `json:"owner"`
	Expires time.Time `json:"expires"`

	TxHash      string `json:"txHash"`
	BlockHash   string `json:"blockHash"`
	BlockNumber uint64 `json:"blockNumber"`
}

func (e *ENS) DomainString() string {
	return fmt.Sprintf("%s.%s", e.TLD, e.Name)
}
