package datagateway

import (
	"context"

	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
)

type mockENS struct {
	ens     *entity.ENS
	removed bool
}

type mockDataGatewayENS struct {
	m map[string]*mockENS
}

func NewMockDataGatewayENS() DataGatewayENS {
	return &mockDataGatewayENS{
		m: make(map[string]*mockENS),
	}
}

func (m *mockDataGatewayENS) SetENS(ctx context.Context, ens *entity.ENS) error {
	m.m[ens.Name] = &mockENS{
		ens:     ens,
		removed: false,
	}

	return nil
}

func (m *mockDataGatewayENS) GetENS(ctx context.Context, domainName string) (*entity.ENS, error) {
	saved, ok := m.m[domainName]
	if !ok || saved == nil {
		return nil, errors.Wrapf(ErrRecordNotFound, "not found for key %s", domainName)
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
	m.m[ens.Name] = nil
	return nil
}
