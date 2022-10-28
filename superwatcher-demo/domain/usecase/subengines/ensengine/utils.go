package ensengine

import (
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/logutils"
	"github.com/ethereum/go-ethereum/common"
)

// Multiple event logs from ENS have owner or new owner address in log.Data, not in topics.
func extractOwnerAddressFromUnpacked(unpacked map[string]interface{}) (common.Address, error) {
	return logutils.ExtractFieldFromUnpacked[common.Address](unpacked, "owner")
}

func extractTTLFromUnpacked(unpacked map[string]interface{}) (uint64, error) {
	return logutils.ExtractFieldFromUnpacked[uint64](unpacked, "ttl")
}
