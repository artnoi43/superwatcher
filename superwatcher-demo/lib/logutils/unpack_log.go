package logutils

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

// EventUnpackInputsIntoMap use given event to parse input from binary format to Go map.
func EventUnpackInputsIntoMap(
	event abi.Event,
	log *types.Log,
) (
	map[string]interface{},
	error,
) {
	return nil, errors.New("not implemented")
}

// UnpackLogDataIntoMap uses contract's ABI to parse the data bytes into map of key-value pair
func UnpackLogDataIntoMap(
	contractABI abi.ABI,
	eventKey string,
	logData []byte,
) (
	map[string]interface{},
	error,
) {
	unpacked := make(map[string]interface{})
	if err := contractABI.UnpackIntoMap(unpacked, eventKey, logData); err != nil {
		return nil, errors.Wrapf(err, "failed to unpack eventKey %s with log data %s", eventKey, common.Bytes2Hex(logData))
	}

	return unpacked, nil
}
