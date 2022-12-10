package watcherstate

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"

	"github.com/artnoi43/superwatcher/superwatcher-demo/internal/domain/datagateway"
)

const (
	// superwatcher:${serviceKey}:${field..}
	redisKeyBase                   = "superwatcher:%s"
	redisKeyState                  = redisKeyBase + ":state" // Unused for now
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
) superwatcher.StateDataGateway {
	return &watcherStateRedisCli{
		// Format strings now to save CPU costs later
		keyBase:      fmt.Sprintf(redisKeyBase, serviceKey),
		keyState:     fmt.Sprintf(redisKeyState, serviceKey),
		keyLastBlock: fmt.Sprintf(redisKeyStateLastRecordedBlock, serviceKey),
		redisClient:  redisCli,
	}
}

func (s *watcherStateRedisCli) SetLastRecordedBlock(
	ctx context.Context,
	blockNumber uint64,
) error {
	return datagateway.HandleRedisErr( //nolint:wrapcheck
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
		return 0, datagateway.HandleRedisErr(err, "get lastRecordedBlock", s.keyLastBlock) //nolint:wrapcheck
	}

	lastRecordedBlock, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse lastRecordedBlock string to uint64: \"%s\"", val)
	}

	return lastRecordedBlock, nil
}

func (s *watcherStateRedisCli) Shutdown() error {
	return errors.Wrap(s.redisClient.Close(), "error shutting down watcherStateRedisCli")
}
