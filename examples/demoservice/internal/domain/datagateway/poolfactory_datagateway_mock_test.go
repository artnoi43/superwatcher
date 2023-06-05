package datagateway

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/soyart/superwatcher"
	"github.com/soyart/superwatcher/examples/demoservice/internal/domain/entity"
)

func TestMockDataGatewayPoolFactory(t *testing.T) {
	dgw := NewMockDataGatewayPoolFactory()

	l := int64(5)
	var addrs []common.Address
	var pools []*entity.Uniswapv3PoolCreated

	for i := int64(1); i <= l; i++ {
		addr := common.BigToAddress(big.NewInt(i))
		pool := &entity.Uniswapv3PoolCreated{
			Address:      addr,
			Fee:          uint64(i + 69),
			BlockCreated: uint64(i + 10000),
		}

		addrs = append(addrs, addr)
		pools = append(pools, pool)
	}

	for i, pool := range pools {
		addr := addrs[i]

		if err := dgw.SetPool(nil, pool); err != nil {
			t.Error("failed to SetPool", err.Error())
		}
		resultPool, err := dgw.GetPool(nil, addr)
		if err != nil {
			t.Error("failed to GetPool", err.Error())
		}

		if !reflect.DeepEqual(*pool, *resultPool) {
			t.Fatalf("unexpected pool, expecting %+v, got %+v", pool, resultPool)
		}

		if err := dgw.DelPool(nil, pool); err != nil {
			t.Error("failed to DelPool", err.Error())
		}

		out, err := dgw.GetPool(nil, addr)
		if err != nil {
			if errors.Is(err, superwatcher.ErrRecordNotFound) {
				if out != nil {
					t.Fatalf("got non nil pool %s", addr)
				}
				continue
			}

			t.Fatalf("unexpected error from GetPool %s", addr)
		}

		if out != nil {
			t.Fatalf("got non nil pool %s", addr)
		}
	}
}
