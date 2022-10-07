package datagateway

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/lib/enums"
)

// LogKey is used as the primary key for superWatcher/watcherClient
// to access stored LogData on the database.
type LogKey struct {
	Chain      enums.ChainType
	ArkenTopic []common.Hash
}

// SuperWatcherLogData is a raw/minimal storage unit of an Ethereum event log.
// Its main purpose is to capture chain data snapshot for handling chain reorgs.
// It is saved directly by superwatcher service, and not used by external services.
type LogData interface {
	Chain() enums.ChainType
	ContractAddress() common.Address
	BlockHash() common.Hash
	TxHash() common.Hash
	BlockNumber() uint64
	LogData() []byte
	BlockTime() time.Time
	Reorged() bool
}
