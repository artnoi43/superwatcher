package datagateway

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
)

type mockENS struct {
	ens     *entity.ENS
	removed bool
}

type mockDataGatewayENS struct {
	// m maps ID to *mockENS
	m map[string]*mockENS
}

func NewMockDataGatewayENS() DataGatewayENS {
	return &mockDataGatewayENS{
		m: make(map[string]*mockENS),
	}
}

func (m *mockDataGatewayENS) SetENS(ctx context.Context, ens *entity.ENS) error {
	m.m[ens.ID] = &mockENS{
		ens:     ens,
		removed: false,
	}

	return nil
}

func (m *mockDataGatewayENS) GetENS(ctx context.Context, ensID string) (*entity.ENS, error) {
	saved, ok := m.m[ensID]
	if !ok || saved == nil {
		return nil, errors.Wrapf(ErrRecordNotFound, "ens not found for key %s", ensID)
	}

	return saved.ens, nil
}
func (m *mockDataGatewayENS) GetENSes(context.Context) ([]*entity.ENS, error) {
	var enses []*entity.ENS
	for _, saved := range m.m {
		enses = append(enses, saved.ens)
	}

	return enses, nil
}

func (m *mockDataGatewayENS) DelENS(ctx context.Context, ens *entity.ENS) error {
	_, ok := m.m[ens.ID]
	if !ok {
		return errors.Wrapf(ErrRecordNotFound, "ens not found for key %s", ens.Name)
	}
	fmt.Println("# DEL ENS", ens.Name)
	m.m[ens.ID] = nil

	return nil
}
