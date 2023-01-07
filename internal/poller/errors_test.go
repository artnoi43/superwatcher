package poller

import (
	"testing"

	"github.com/artnoi43/superwatcher"
	"github.com/pkg/errors"
)

func TestPollerErrors(t *testing.T) {
	err := getErrHashesDiffer()
	if !errors.Is(err, errHashesDiffer) {
		t.Errorf("err %s is not errHashesDiffer %s", err.Error(), errHashesDiffer.Error())
	}
	if !errors.Is(err, superwatcher.ErrChainIsReorging) {
		t.Errorf("err %s is not ErrChainIsReorging %s", err.Error(), superwatcher.ErrChainIsReorging)
	}
}

func getErrHashesDiffer() error {
	return errors.Wrap(errHashesDiffer, "foo")
}
