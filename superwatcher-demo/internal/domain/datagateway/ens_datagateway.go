package datagateway

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/entity"
)

func EnsRedisKey(domainName string) string {
	return fmt.Sprintf("demo:ens:%s", domainName)
}

type DataGatewayENS interface {
	SetENS(context.Context, *entity.ENS) error
	GetENS(context.Context, string) (*entity.ENS, error)
	GetENSes(context.Context) ([]*entity.ENS, error)
	DelENS(context.Context, *entity.ENS) error
}

type dataGatewayENS struct {
	redisClient *redis.Client
}

func NewEnsDataGateway(redisCli *redis.Client) *dataGatewayENS {
	return &dataGatewayENS{
		redisClient: redisCli,
	}
}

// handleRedisErr checks if err is redis.Nil
func handleRedisErr(err error, action, key string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		err = errors.Wrap(ErrRecordNotFound, err.Error())
		err = WrapErrRecordNotFound(err, key)
	}
	return errors.Wrapf(err, "action: %s", action)
}

func (s *dataGatewayENS) SetENS(
	ctx context.Context,
	ens *entity.ENS,
) error {
	key := EnsRedisKey(ens.ID)
	ensJSON, err := json.Marshal(ens)
	if err != nil {
		return errors.Wrap(err, "failed to marshal ENS")
	}

	return handleRedisErr(
		s.redisClient.Set(ctx, key, ensJSON, -1).Err(),
		"set RecordedENS",
		key,
	)
}

func (s *dataGatewayENS) GetENS(
	ctx context.Context,
	key string,
) (*entity.ENS, error) {
	stringData, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, handleRedisErr(err, "get RecordedENS", key)
	}

	ens := &entity.ENS{}
	err = json.Unmarshal([]byte(stringData), ens)
	if err != nil {
		return nil, errors.Wrapf(err, "action: %s", "unmarshal RecordedENS")
	}
	return ens, nil
}

func (s *dataGatewayENS) GetENSes(
	ctx context.Context,
) (
	[]*entity.ENS,
	error,
) {
	ENSs := []*entity.ENS{}
	var cursor uint64
	prefix := EnsRedisKey("*")

	for {
		var keys []string
		var err error
		keys, cursor, err = s.redisClient.Scan(ctx, cursor, prefix, 0).Result()
		if err != nil {
			return nil, errors.Wrapf(err, "action: %s", "scan RecordedENS")
		}

		for _, key := range keys {
			ens, err := s.GetENS(ctx, key)
			if err != nil {
				return nil, errors.Wrapf(err, "action: %s", "get RecordedENS")
			}

			ENSs = append(ENSs, ens)
		}

		if cursor == 0 {
			break
		}
	}

	return ENSs, nil
}

func (s *dataGatewayENS) DelENS(
	ctx context.Context,
	ens *entity.ENS,
) error {
	key := EnsRedisKey(ens.ID)
	_, err := s.redisClient.Del(ctx, key).Result()
	if err != nil {
		return handleRedisErr(err, "del RecordedENS", key)
	}
	return nil
}

func (s *dataGatewayENS) Shutdown() error {
	return s.redisClient.Close()
}
