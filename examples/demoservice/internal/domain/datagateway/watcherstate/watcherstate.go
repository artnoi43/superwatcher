package watcherstate

import (
	"errors"

	"github.com/soyart/superwatcher"

	"github.com/soyart/superwatcher/examples/demoservice/internal/domain/datagateway"
	"github.com/soyart/superwatcher/examples/demoservice/internal/watcherstate"
)

// NewRedisStateDataGateway returns default implementation of StateDataGateway.
// It uses |serviceKey| to compose a Redis key to independently store multiple
// superwatcher-derived services on the same Redis database.
// If you only use default `superwatcher.Emitter` implementation for your service,
// then **your own code is responsible for calling `SetLastRecordedBlock`**.
func NewRedisStateDataGateway(
	serviceKey string,
	redisClient datagateway.RedisClient,
) (
	superwatcher.StateDataGateway,
	error,
) {
	if redisClient == nil {
		return nil, errors.New("nil redisClient")
	}
	if serviceKey == "" {
		return nil, errors.New("empty serviceName")
	}

	return watcherstate.NewRedisWatcherStateDataGateway(serviceKey, redisClient), nil
}
