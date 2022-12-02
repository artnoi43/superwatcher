package watcherstate

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/pkg/datagateway"
)

const (
	// superwatcher:${chain}:${service}:${field..}
	redisKeyBase                   = "superwatcher:%s"
	redisKeyState                  = redisKeyBase + ":state"
	redisKeyStateLastRecordedBlock = redisKeyState + ":lastRecordedBlock"
)

// watcherStateRedisCli is the default Redis-based implementation for watcherstate.WatcherStateDataGateway
type watcherStateRedisCli struct {
	keyBase, keyState, keyLastBlock string
	redisClient                     datagateway.RedisClient
}

func NewRedisWatcherStateDataGateway(
	serviceKey string, // Each key for each different service and chain
	redisCli datagateway.RedisClient,
) *watcherStateRedisCli {
	return &watcherStateRedisCli{
		// Format strings now to save CPU costs later
		keyBase:      fmt.Sprintf(redisKeyBase, serviceKey),
		keyState:     fmt.Sprintf(redisKeyState, serviceKey),
		keyLastBlock: fmt.Sprintf(redisKeyStateLastRecordedBlock, serviceKey),
		redisClient:  redisCli,
	}
}

// handleRedisErr checks if err is redis.Nil, and if so,
// wraps it with datagateway.ErrRecordNotFound.
func handleRedisErr(err error, action, key string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		err = errors.Wrap(datagateway.ErrRecordNotFound, err.Error())
		err = datagateway.WrapErrRecordNotFound(err, key)
	}
	return errors.Wrapf(err, "action: %s", action)
}

func (s *watcherStateRedisCli) SetLastRecordedBlock(
	ctx context.Context,
	blockNumber uint64,
) error {
	return handleRedisErr(
		s.redisClient.Set(ctx, s.keyLastBlock, blockNumber, -1).Err(),
		"set lastRecordedBlock",
		s.keyLastBlock,
	)
}

func (s *watcherStateRedisCli) GetLastRecordedBlock(
	ctx context.Context,
) (uint64, error) {
	val, err := s.redisClient.Get(ctx, s.keyLastBlock).Result()
	if err != nil {
		return 0, handleRedisErr(err, "get lastRecordedBlock", s.keyLastBlock)
	}

	lastRecordedBlock, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse lastRecordedBlock string to uint64: \"%s\"", val)
	}

	return lastRecordedBlock, nil
}

func (s *watcherStateRedisCli) Shutdown() error {
	return s.redisClient.Close()
}
