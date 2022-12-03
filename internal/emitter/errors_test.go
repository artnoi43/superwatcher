package emitter

import (
	"testing"

	"github.com/pkg/errors"
)

func TestWrapErrBlockNumber(t *testing.T) {
	baseErr := errors.New("failed to filterLogs")
	err := wrapErrFetchError(baseErr)
	if !errors.Is(err, errFetchError) {
		t.Error("err is not errFetchError")
	}

	err = mockErrNoHash()
	if !errors.Is(err, errProcessReorg) {
		t.Error("err is not errProcessReorg")
	}
}

func wrapErrFetchError(err error) error {
	return errors.Wrap(errFetchError, err.Error())
}

func mockErrNoHash() error {
	return errors.Wrapf(errNoHash, "blockNumber %d", 69)
}
