package watcherstate

import (
	"context"
	"errors"

	"github.com/artnoi43/superwatcher/internal/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/pkg/enums"
)

// StateDataGateway is used by the default emitter to get LastRecordedBlock,
// and by the default engine to set LastRecordedBlock
type StateDataGateway interface {
	GetLastRecordedBlock(context.Context) (uint64, error)
	SetLastRecordedBlock(context.Context, uint64) error

	Shutdown() error
}

// NewRedisStateDataGateway returns default implementation of StateDataGateway.
// It uses |serviceName| to compose a Redis key to independently store multiple
// superwatcher-derived services on the same Redis database.
// If you only use default `superwatcher.WatcherEmitter` implementation for your service,
// then **your own code is responsible for calling `SetLastRecordedBlock`**.
func NewRedisStateDataGateway(
	chain enums.ChainType,
	serviceName string,
	redisClient datagateway.RedisClient,
) (
	StateDataGateway,
	error,
) {
	if redisClient == nil {
		return nil, errors.New("nil redisClient")
	}
	if serviceName == "" {
		return nil, errors.New("empty serviceName")
	}

	return watcherstate.NewRedisWatcherStateDataGateway(chain, serviceName, redisClient), nil
}
