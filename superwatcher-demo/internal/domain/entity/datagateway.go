package entity

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisClient interface {
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Scan(ctx context.Context, cursor uint64, keymatch string, count int64) *redis.ScanCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Close() error
}

const (
	redisEnsKeyFormat = "demo:ens:%s"
)

type EnsDataGateway struct {
	redisEnsKey string
	redisClient RedisClient
}

func NewEnsDataGateway(
	keyPrefix string,
	redisCli RedisClient,
) *EnsDataGateway {
	return &EnsDataGateway{
		redisEnsKey: fmt.Sprintf(redisEnsKeyFormat, keyPrefix),
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

func (s *EnsDataGateway) SetRecordedENS(
	ctx context.Context,
	ens *ENS,
) error {
	key := fmt.Sprintf(s.redisEnsKey+":%d", ens.ID)
	return handleRedisErr(
		s.redisClient.Set(ctx, key, ens, -1).Err(),
		"set RecordedENS",
		key,
	)
}

func (s *EnsDataGateway) GetRecordedENS(
	ctx context.Context,
	key string,
) (*ENS, error) {
	stringData, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, handleRedisErr(err, "get RecordedENS", key)
	}

	ens := &ENS{}
	err = json.Unmarshal([]byte(stringData), ens)
	if err != nil {
		return nil, errors.Wrapf(err, "action: %s", "unmarshal RecordedENS")
	}
	return ens, nil
}

func (s *EnsDataGateway) GetRecordedENSs(
	ctx context.Context,
) ([]*ENS, error) {
	ENSs := []*ENS{}
	var cursor uint64
	prefix := s.redisEnsKey + ":*"

	for {
		var keys []string
		var err error
		keys, cursor, err = s.redisClient.Scan(ctx, cursor, prefix, 0).Result()
		if err != nil {
			return nil, errors.Wrapf(err, "action: %s", "scan RecordedENS")
		}

		for _, key := range keys {
			ens, err := s.GetRecordedENS(ctx, key)
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

func (s *EnsDataGateway) DelRecordedENS(
	ctx context.Context,
	key string,
) error {
	_, err := s.redisClient.Del(ctx, key).Result()
	if err != nil {
		return handleRedisErr(err, "del RecordedENS", key)
	}
	return nil
}

func (s *EnsDataGateway) Shutdown() error {
	return s.redisClient.Close()
}
