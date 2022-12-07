package datagateway

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
)

// Use Redis Hash Map to store entity.ENS, with ID as field (sub-key)
const EnsRedisKey = "demo:ens"

type DataGatewayENS interface {
	SetENS(context.Context, *entity.ENS) error
	GetENS(context.Context, string) (*entity.ENS, error)
	GetENSes(context.Context) ([]*entity.ENS, error)
	DelENS(context.Context, *entity.ENS) error
}

type dataGatewayENS struct {
	redisCli *redis.Client
}

func NewEnsDataGateway(redisCli *redis.Client) *dataGatewayENS {
	return &dataGatewayENS{
		redisCli: redisCli,
	}
}

// handleRedisErr checks if err is redis.Nil
func handleRedisErr(err error, action, key string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		err = datagateway.WrapErrRecordNotFound(err, key)
	}

	return errors.Wrapf(err, "action: %s", action)
}

func (s *dataGatewayENS) SetENS(
	ctx context.Context,
	ens *entity.ENS,
) error {
	ensJSON, err := json.Marshal(ens)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal ENS for ensID %s", ens.ID)
	}

	id := ens.ID
	err = s.redisCli.HSet(ctx, EnsRedisKey, id, ensJSON).Err()

	return handleRedisErr(err, "HSET ens", id)
}

func (s *dataGatewayENS) GetENS(
	ctx context.Context,
	ensID string,
) (
	*entity.ENS,
	error,
) {
	stringData, err := s.redisCli.HGet(ctx, EnsRedisKey, ensID).Result()
	if err != nil {
		return nil, handleRedisErr(err, "HGET ens", ensID)
	}

	var ens entity.ENS
	err = json.Unmarshal([]byte(stringData), &ens)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal for ensID %s", ensID)
	}
	return &ens, nil
}

func (s *dataGatewayENS) GetENSes(
	ctx context.Context,
) (
	[]*entity.ENS,
	error,
) {
	resultMap, err := s.redisCli.HGetAll(ctx, EnsRedisKey).Result()
	if err != nil {
		return nil, handleRedisErr(err, "HGETALL ens", "null")
	}

	var enses []*entity.ENS
	for ensID, ensJSON := range resultMap {
		var ens entity.ENS
		if err := json.Unmarshal([]byte(ensJSON), &ens); err != nil {
			return nil, errors.Wrapf(err, "failec to unmarshal for ensID %s", ensID)
		}

		enses = append(enses, &ens)
	}

	return enses, nil
}

func (s *dataGatewayENS) DelENS(
	ctx context.Context,
	ens *entity.ENS,
) error {
	id := ens.ID
	_, err := s.redisCli.HDel(ctx, EnsRedisKey, id).Result()
	if err != nil {
		return handleRedisErr(err, "del RecordedENS", EnsRedisKey)
	}

	return nil
}

func (s *dataGatewayENS) Shutdown() error {
	return s.redisCli.Close()
}
