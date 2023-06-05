package datagateway

import (
	"context"
	"time"

	"github.com/soyart/superwatcher"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisClient interface {
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Close() error
}

// wraps it with superwather.ErrRecordNotFound.
func HandleRedisErr(err error, action, key string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		err = superwatcher.WrapErrRecordNotFound(err, key)
	}

	return errors.Wrapf(err, "action: %s", action)
}
