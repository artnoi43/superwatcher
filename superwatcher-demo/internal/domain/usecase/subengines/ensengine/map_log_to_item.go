package ensengine

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/lib/logutils"
)

// Multiple event logs from ENS have owner or new owner internal/address in log.Data, not in topics.
func extractOwnerAddressFromUnpacked(unpacked map[string]interface{}) (common.Address, error) {
	return logutils.ExtractFieldFromUnpacked[common.Address](unpacked, "owner")
}

func extractTTLFromUnpacked(unpacked map[string]interface{}) (uint64, error) {
	return logutils.ExtractFieldFromUnpacked[uint64](unpacked, "ttl")
}

// unmarshalLogToENS populates ens with data from log
func (e *ensEngine) unmarshalLogToENS(
	logEvent string,
	log *types.Log,
	ens *entity.ENS,
) error {
	switch log.Address {
	case e.ensRegistrar.Address:
		switch logEvent {
		case nameRegistered:
			if len(log.Topics) < 3 {
				return errors.Wrap(ErrLogLen, "log topics len < 3")
			}
			unpacked, err := logutils.UnpackLogDataIntoMap(e.ensRegistrar.ContractABI, logEvent, log.Data)
			if err != nil {
				return errors.Wrap(err, ErrMapENS.Error())
			}
			// Extract data from Registrar contract log
			expire, err := logutils.ExtractFieldFromUnpacked[*big.Int](unpacked, "expires")
			if err != nil {
				return errors.Wrap(err, ErrMapENS.Error())
			}

			ens.ID = gslutils.ToLower(log.Topics[1].String())
			ens.Owner = common.HexToAddress(log.Topics[2].Hex())
			ens.Expires = time.Unix(expire.Int64(), 0)
		}

	case e.ensController.Address:
		switch logEvent {
		case nameRegistered:
			// Extract data from Controller contract log yopics and data
			var err error
			unpacked, err := logutils.UnpackLogDataIntoMap(e.ensController.ContractABI, logEvent, log.Data)
			if err != nil {
				return errors.Wrap(err, ErrMapENS.Error())
			}
			name, err := logutils.ExtractFieldFromUnpacked[string](unpacked, "name")
			if err != nil {
				return errors.Wrap(err, ErrMapENS.Error())
			}
			expire, err := logutils.ExtractFieldFromUnpacked[*big.Int](unpacked, "expires")
			if err != nil {
				return errors.Wrap(err, ErrMapENS.Error())
			}

			ens.Name = gslutils.ToLower(name)
			ens.Owner = common.HexToAddress(log.Topics[2].Hex())
			ens.Expires = time.Unix(expire.Int64(), 0)
		}
	}

	return nil
}
