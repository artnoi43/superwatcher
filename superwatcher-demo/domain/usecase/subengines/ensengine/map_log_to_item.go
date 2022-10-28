package ensengine

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/logutils"
)

// Multiple event logs from ENS have owner or new owner address in log.Data, not in topics.
func extractOwnerAddressFromUnpacked(unpacked map[string]interface{}) (common.Address, error) {
	return logutils.ExtractFieldFromUnpacked[common.Address](unpacked, "owner")
}

func extractTTLFromUnpacked(unpacked map[string]interface{}) (uint64, error) {
	return logutils.ExtractFieldFromUnpacked[uint64](unpacked, "ttl")
}

// TODO: this is broken. To create a new ENS name, we need logs from 2 contracts,
// but 1 log can only come from 1 contract
func (e *ensEngine) newName(log *types.Log, logEvent string) (*entity.ENS, error) {
	if len(log.Topics) < 2 {
		panic("bad log with < 2 topics")
	}
	unpackedRegistrar, err := logutils.UnpackLogDataIntoMap(e.ensRegistrar.ContractABI, logEvent, log.Data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unpack ENSRegistrar log event %s", logEvent)
	}
	unpackedController, err := logutils.UnpackLogDataIntoMap(e.ensController.ContractABI, logEvent, log.Data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unpack ENSController log event %s", logEvent)
	}

	var name entity.ENS
	switch logEvent {
	case "NameRegistered":
		var err error
		name.Name, err = logutils.ExtractFieldFromUnpacked[string](unpackedRegistrar, "name")
		if err != nil {
			return nil, errors.Wrap(err, "failed to get domain name from controller log")
		}
		expire, err := logutils.ExtractFieldFromUnpacked[*big.Int](unpackedController, "expire")
		if err != nil {
			return nil, errors.Wrap(err, "failed to get expiration from registrar log")
		}
		name.ID = common.HexToHash(log.Topics[1].Hex()).String()
		name.Owner = common.HexToAddress(log.Topics[2].Hex())
		name.Expire = time.Unix(expire.Int64(), 0)
	}

	return &name, nil
}
