package watcherstate

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/pkg/enums"
)

const (
	// superwatcher:${chain}:${service}:${field..}
	redisKeyBase                   = "superwatcher:%s:%s"
	redisKeyState                  = redisKeyBase + ":state"
	redisKeyStateLastRecordedBlock = redisKeyState + ":lastRecordedBlock"
)

type watcherStateRedisCli struct {
	keyBase, keyState, keyLastBlock string
	redisCli                        *redis.Client
}

func NewWatcherStateRedisClient(
	chain enums.ChainType, // Each key for different chain
	serviceName string, // Each key for different service
	redisCli *redis.Client,
) *watcherStateRedisCli {
	return &watcherStateRedisCli{
		// Format strings now to save CPU costs later
		keyBase:      fmt.Sprintf(redisKeyBase, chain, serviceName),
		keyState:     fmt.Sprintf(redisKeyState, chain, serviceName),
		keyLastBlock: fmt.Sprintf(redisKeyStateLastRecordedBlock, chain, serviceName),
		redisCli:     redisCli,
	}
}

func handleRedisErr(err error, action string) error {
	if errors.Is(err, redis.Nil) {
		return errors.Wrap(datagateway.ErrRecordNotFound, err.Error())
	}
	return errors.Wrap(err, action)
}

func (s *watcherStateRedisCli) SetLastRecordedBlock(
	ctx context.Context,
	blockNumber uint64,
) error {
	return handleRedisErr(
		s.redisCli.Set(ctx, s.keyLastBlock, blockNumber, -1).Err(),
		"set lastRecordedBlock",
	)
}

func (s *watcherStateRedisCli) GetLastRecordedBlock(
	ctx context.Context,
) (uint64, error) {
	val, err := s.redisCli.Get(ctx, s.keyLastBlock).Result()
	if err != nil {
		return 0, handleRedisErr(err, "get lastRecordedBlock")
	}

	lastRecordedBlock, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse lastRecordedBlock string to uint: \"%s\"", val)
	}

	return lastRecordedBlock, nil
}

func (s *watcherStateRedisCli) Shutdown() error {
	return s.redisCli.Close()
}
