package watcherstate

import (
	"errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/datagateway"
)

// NewRedisStateDataGateway returns default implementation of StateDataGateway.
// It uses |serviceName| to compose a Redis key to independently store multiple
// superwatcher-derived services on the same Redis database.
// If you only use default `superwatcher.WatcherEmitter` implementation for your service,
// then **your own code is responsible for calling `SetLastRecordedBlock`**.
func NewRedisStateDataGateway(
	serviceName string,
	redisClient datagateway.RedisClient,
) (
	superwatcher.StateDataGateway,
	error,
) {
	if redisClient == nil {
		return nil, errors.New("nil redisClient")
	}
	if serviceName == "" {
		return nil, errors.New("empty serviceName")
	}

	return watcherstate.NewRedisWatcherStateDataGateway(serviceName, redisClient), nil
}
