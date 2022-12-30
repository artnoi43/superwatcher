package superwatcher

import (
	"testing"

	"github.com/pkg/errors"
)

func TestWrapErrRecordNotFound(t *testing.T) {
	err := errors.New("some err")
	err = WrapErrRecordNotFound(err, "some key")

	if !errors.Is(err, ErrRecordNotFound) {
		t.Fatalf("error \"%s\" is not ErrRecordNotFound", err.Error())
	}
}
