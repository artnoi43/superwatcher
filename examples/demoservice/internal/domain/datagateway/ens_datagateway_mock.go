package datagateway

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/examples/demoservice/internal/domain/entity"
)

type mockDataGatewayENS struct {
	// m maps ID to *mockENS
	m map[string]*entity.ENS
}

func NewMockDataGatewayENS() RepositoryENS {
	return &mockDataGatewayENS{
		m: make(map[string]*entity.ENS),
	}
}

func (s *mockDataGatewayENS) SetENS(ctx context.Context, ens *entity.ENS) error {
	s.m[ens.ID] = ens

	return nil
}

func (s *mockDataGatewayENS) GetENS(ctx context.Context, ensID string) (*entity.ENS, error) {
	saved, ok := s.m[ensID]
	if !ok || saved == nil {
		return nil, errors.Wrapf(superwatcher.ErrRecordNotFound, "ens not found for key %s", ensID)
	}

	return saved, nil
}

func (s *mockDataGatewayENS) GetENSes(context.Context) ([]*entity.ENS, error) {
	var enses []*entity.ENS //nolint:prealloc
	for _, saved := range s.m {
		enses = append(enses, saved)
	}

	return enses, nil
}

func (s *mockDataGatewayENS) DelENS(ctx context.Context, ens *entity.ENS) error {
	_, ok := s.m[ens.ID]
	if !ok {
		return errors.Wrapf(superwatcher.ErrRecordNotFound, "ens not found for key %s", ens.Name)
	}

	fmt.Println("# DEL ENS", ens.Name, ens.ID)
	s.m[ens.ID] = nil

	return nil
}
