package logutils

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

// EventUnpackInputsIntoMap use given event to parse input from binary format to Go map.
func EventUnpackInputsIntoMap(event abi.Event, log *types.Log) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func UnpackIntoMap(contractABI abi.ABI, eventKey string, log *types.Log) (map[string]interface{}, error) {
	var unpacked map[string]interface{}
	if err := contractABI.UnpackIntoMap(unpacked, eventKey, log.Data); err != nil {
		return nil, errors.Wrapf(err, "failed to unpack eventKey %s with log data %s", eventKey, common.Bytes2Hex(log.Data))
	}

	return unpacked, nil
}
