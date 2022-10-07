package data

//import (
//	"time"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/core/types"
//
//	"github.com/artnoi43/superwatcher/domain/datagateway"
//	"github.com/artnoi43/superwatcher/lib/enums"
//)
//
//// logData is what gets saved to database,
//// and it should implements datagateway.LogData
//type logData struct {
//	chain       enums.ChainType
//	address     common.Address
//	blockHash   common.Hash
//	txHash      common.Hash
//	topics      []common.Hash
//	blockNumber uint64
//	logData     []byte
//	blockTime   time.Time
//	reorged     bool
//}
//
//func fromLog(l *types.Log, b *types.Block, chain enums.ChainType) datagateway.LogData {
//	return &logData{
//		chain:       chain,
//		address:     l.Address,
//		blockHash:   l.BlockHash,
//		txHash:      l.TxHash,
//		topics:      l.Topics,
//		blockNumber: l.BlockNumber,
//		logData:     l.Data,
//		blockTime:   time.Unix(int64(b.Header().Time), 0).UTC(),
//		reorged:     l.Removed,
//	}
//}
//
//// Implements datagateway.LogData
//func (d *logData) Chain() enums.ChainType {
//	return d.chain
//}
//
//// Implements datagateway.LogData
//func (d *logData) ContractAddress() common.Address {
//	return d.address
//}
//
//// Implements datagateway.LogData
//func (d *logData) BlockHash() common.Hash {
//	return d.blockHash
//}
//
//// Implements datagateway.LogData
//func (d *logData) TxHash() common.Hash {
//	return d.txHash
//}
//
//// Implements datagateway.LogData
//func (d *logData) BlockNumber() uint64 {
//	return d.blockNumber
//}
//
//// Implements datagateway.LogData
//func (d *logData) BlockTime() time.Time {
//	return d.blockTime
//}
//
//// Implements datagateway.LogData
//func (d *logData) Reorged() bool {
//	return d.reorged
//}
//
//// Implements datagateway.LogData
//func (d *logData) LogData() []byte {
//	return d.logData
//}
