package datagateway

import "github.com/pkg/errors"

// This error is checked for in emitter.loopFilterLogs.
var ErrRecordNotFound = errors.New("record not found")

func WrapErrRecordNotFound(err error, keyNotFound string) error {
	return errors.Wrapf(err, "key %s not found", keyNotFound)
}
