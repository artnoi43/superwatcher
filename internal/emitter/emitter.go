package emitter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitter/enums"
)

type ethClient interface {
	BlockNumber(context.Context) (uint64, error)
	BlockByNumber(context.Context, *big.Int) (*types.Block, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

type Emitter interface {
	Loop(context.Context) error
	shutdown()
}

type emitter struct {
	config    *config.Config
	client    ethClient
	addresses []common.Address
	topics    [][]common.Hash

	logChan            chan<- *types.Log
	blockChan          chan<- *enums.BlockInfo
	reorgChan          chan<- *enums.BlockInfo
	errChan            chan<- error
	isSolvingReorgChan chan int

	debug bool
}

type Config struct {
	Chain           enums.ChainType `mapstructure:"chain" json:"chain"`
	Node            string          `mapstructure:"node_url" json:"node"`
	StartBlock      uint64          `mapstructure:"start_block" json:"startBlock"`
	LookBackBlocks  uint64          `mapstructure:"lookback_blocks" json:"lookBackBlock"`
	LookBackRetries uint64          `mapstructure:"lookback_retries" json:"lookBackRetries"`
	IntervalSecond  int             `mapstructure:"interval_second" json:"intervalSecond"`
}

func NewEmitter(
	conf *config.Config,
	client ethClient,
	addresses []common.Address,
	topics [][]common.Hash,
	logChan chan<- *types.Log,
	reorgChan chan<- *enums.BlockInfo,
	isSolvingReorgChan chan int,
) Emitter {
	return &emitter{
		config: conf,
		client: client,

		addresses: addresses,
		topics:    topics,
		logChan:   logChan,

		reorgChan:          reorgChan,
		isSolvingReorgChan: isSolvingReorgChan,
	}
}

func (e *emitter) Loop(ctx context.Context) error {
	fmt.Println("Loop")
	fmt.Println("ctx", ctx)
	for {

		select {
		case <-ctx.Done():
			fmt.Println("shutdown")
			e.shutdown()
			return ctx.Err()
		default:
			for i := 1; i < 10; i++ {
				if i%3 == 0 {
					fmt.Println("i reorg--->", i)
					e.reorgChan <- &enums.BlockInfo{
						Number: 10,
						Hash:   common.HexToHash("0x5ac9b37d571677b80957ca05693f371526c602fd08042b416a29fdab7efefa49"),
						Logs: []*types.Log{{
							Address:     common.HexToAddress("0x0000000000000000000000000000000000001003"),
							Topics:      []common.Hash{common.HexToHash("0x5ac9b37d571677b80957ca05693f371526c602fd08042b416a29fdab7efefa49")},
							Data:        common.Hex2Bytes("0x0000000000000000000000000000000000000000000000000000000006915167cedaf7bbf7df47d932fdda630527ee648562cf3e52c5e5f46156a3a971a4ceb4"),
							BlockNumber: hexutil.MustDecodeUint64("0x1"),
							TxHash:      common.HexToHash("0x9ebc5237eabb339a103a34daf280db3d9498142b49fa47f1af71f64a605acffa"),
							TxIndex:     uint(hexutil.MustDecodeUint64("0x2")),
							BlockHash:   common.HexToHash("0x04055304e432294a65ff31069c4d3092ff8b58f009cdb50eba5351e0332ad0f6"),
							Index:       uint(hexutil.MustDecodeUint64("0x0")),
							Removed:     false,
						}},
					}
					<-e.isSolvingReorgChan

				} else {
					fmt.Println("i normal--->", i)
					e.logChan <- &types.Log{
						Address:     common.HexToAddress("0x0000000000000000000000000000000000001003"),
						Topics:      []common.Hash{common.HexToHash("0x5ac9b37d571677b80957ca05693f371526c602fd08042b416a29fdab7efefa49")},
						Data:        common.Hex2Bytes("0x0000000000000000000000000000000000000000000000000000000006915167cedaf7bbf7df47d932fdda630527ee648562cf3e52c5e5f46156a3a971a4ceb4"),
						BlockNumber: hexutil.MustDecodeUint64("0x1"),
						TxHash:      common.HexToHash("0x9ebc5237eabb339a103a34daf280db3d9498142b49fa47f1af71f64a605acffa"),

						TxIndex:   uint(i),
						BlockHash: common.HexToHash("0x04055304e432294a65ff31069c4d3092ff8b58f009cdb50eba5351e0332ad0f6"),
						Index:     uint(hexutil.MustDecodeUint64("0x0")),
						Removed:   false,
					}

				}
				time.Sleep(2 * time.Second)
			}
		}

	}
}

func (e *emitter) shutdown() {
	close(e.logChan)
	close(e.blockChan)
	close(e.reorgChan)
	close(e.errChan)
}
