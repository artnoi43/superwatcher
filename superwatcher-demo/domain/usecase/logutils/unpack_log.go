package logutils

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func UnpackIntoMap(contractABI abi.ABI, eventKey string, logData []byte) (map[string]interface{}, error) {
	var unpacked map[string]interface{}
	if err := contractABI.UnpackIntoMap(unpacked, eventKey, logData); err != nil {
		return nil, errors.Wrapf(err, "failed to unpack eventKey %s with log data %s", eventKey, common.Bytes2Hex(logData))
	}

	return unpacked, nil
}
