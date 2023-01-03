package logutils

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

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

func ExtractFieldFromUnpacked[T any](unpacked map[string]interface{}, field string) (T, error) {
	var t T

	v, found := unpacked[field]
	if !found {
		return t, fmt.Errorf("field %s not found in unpacked", field)
	}

	value, ok := v.(T)
	if !ok {
		return t, fmt.Errorf("type assertion failed: field %s is not %s", field, reflect.TypeOf(t))
	}

	return value, nil
}
