package datagateway

import (
	"fmt"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
)

func TestMockDataGatewayENS(t *testing.T) {
	dgw := NewMockDataGatewayENS()

	l := 5
	var enses []*entity.ENS
	var names []string
	var ids []string

	for i := 1; i <= l; i++ {
		id := fmt.Sprintf("%d", i)
		s100 := fmt.Sprintf("%d", i+100)
		name := "ens" + id
		ens := &entity.ENS{
			ID:        id,
			Name:      name,
			TxHash:    gslutils.StringerToLowerString(common.HexToHash("0x" + id)),
			BlockHash: gslutils.StringerToLowerString(common.HexToHash("0x" + s100)),
		}

		enses = append(enses, ens)
		names = append(names, name)
		ids = append(ids, id)
	}

	for i, ens := range enses {
		if err := dgw.SetENS(nil, ens); err != nil {
			t.Error("failed to SetENS:", err.Error())
		}

		key := ids[i]
		out, err := dgw.GetENS(nil, key)
		if err != nil {
			t.Fatal("failed to GetENS:", err.Error())
		}
		if out != ens {
			t.Fatal("got different ENS:")
		}

		err = dgw.DelENS(nil, ens)
		if err != nil {
			t.Fatal("failed to DelENS:", err.Error())
		}

		out, err = dgw.GetENS(nil, key)
		if err != nil {
			if errors.Is(err, ErrRecordNotFound) {
				if out != nil {
					t.Fatal("got non-nil ens after call to DelENS")
				}
				continue
			}

			t.Fatal("unexpected error", err.Error())
		}

		if out != nil {
			t.Fatal("got non-nil ens after call to DelENS")
		}
	}
}
