package superwatcher

import (
	"testing"

	"github.com/pkg/errors"
)

func TestWrapErrBlockNumber(t *testing.T) {
	baseErr := errors.New("failed to filterLogs")
	err := wrapErrFetchError(baseErr)
	if !errors.Is(err, ErrFetchError) {
		t.Error("err is not errFetchError")
	}

	err = mockErrNoHash()
	if !errors.Is(err, ErrProcessReorg) {
		t.Error("err is not errProcessReorg")
	}

	err = mockErrFromBlockReorged()
	if !errors.Is(err, ErrFromBlockReorged) {
		t.Error("err is not errFromBlockReorged")
	}

	if !errors.Is(ErrFromBlockReorged, ErrChainIsReorging) {
		t.Error("errFromBlockReorged is not ErrChainIsReorging")
	}
}

func wrapErrFetchError(err error) error {
	return errors.Wrap(ErrFetchError, err.Error())
}

func mockErrNoHash() error {
	err := errors.Wrap(ErrProcessReorg, "no block hash")
	return errors.Wrapf(err, "blockNumber %d", 69)
}

func mockErrFromBlockReorged() error {
	return errors.Wrapf(ErrFromBlockReorged, "fromBlock %d was removed (chain reorganization)", 69)
}
