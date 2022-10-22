package uniswapv3factoryengine

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/superwatcher-demo/domain/entity"
	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/logutils"
)

// mapLogToPoolCreated maps *types.Log to entity.Uniswapv3PoolCreated.
// It can be updated to handle more events from this contract.
func mapLogToPoolCreated(
	contractABI abi.ABI,
	eventName string,
	l *types.Log,
) (
	*entity.Uniswapv3PoolCreated,
	error,
) {
	unpacked, err := logutils.UnpackLogDataIntoMap(contractABI, eventName, l.Data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unpack uniswapv3factory logs (event %s)", eventName)
	}
	poolCreated, err := parseLogDataToUniswapv3Factory(unpacked)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert unpacked log data to *entity.Uniswapv3PoolCreated")
	}

	poolCreated.Token0 = common.HexToAddress(l.Topics[1].String())
	poolCreated.Token1 = common.HexToAddress(l.Topics[2].String())
	poolCreated.Fee = l.Topics[3].Big().Uint64()
	poolCreated.BlockCreated = l.BlockNumber

	return poolCreated, nil
}
