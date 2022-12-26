package entity

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMarshalENS(t *testing.T) {
	ens := ENS{
		ID:      "foo",
		Expires: time.Now(),
	}

	if _, err := json.Marshal(ens); err != nil {
		t.Error(err.Error())
	}
}
