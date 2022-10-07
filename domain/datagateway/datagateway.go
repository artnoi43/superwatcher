package datagateway

import (
	"github.com/ethereum/go-ethereum/common"
)

// DataGateway defines methods used by superwatcher to save data to databases.
type DataGateway interface {
	SaveLogData(LogKey, LogData) error
	GetLogData(LogKey) (LogData, error)

	LogReorg(k LogKey, newHash common.Hash) error
}
