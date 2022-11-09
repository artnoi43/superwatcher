package watcherstate

import (
	"context"

	"github.com/artnoi43/superwatcher/internal/watcherstate"
	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/pkg/enums"
	"github.com/artnoi43/superwatcher/pkg/logger"
)

// StateDataGateway is used by emitter to get LastRecordedBlock,
// and by engine to set LastRecordedBlock
type StateDataGateway interface {
	GetLastRecordedBlock(context.Context) (uint64, error)
	SetLastRecordedBlock(context.Context, uint64) error

	Shutdown() error
}

func NewRedisStateDataGateway(
	chain enums.ChainType,
	serviceName string,
	redisClient datagateway.RedisClient,
) StateDataGateway {
	if redisClient == nil {
		logger.Panic("nil redisClient")
	}

	return watcherstate.NewRedisWatcherStateDataGateway(chain, serviceName, redisClient)
}
