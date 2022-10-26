package engine

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitter"
	"github.com/artnoi43/superwatcher/internal/emitter/enums"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ethClient interface {
	BlockNumber(context.Context) (uint64, error)
	BlockByNumber(context.Context, *big.Int) (*types.Block, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

type engine struct {
	emitter            emitter.Emitter
	logChan            chan *types.Log
	reorgChan          chan *enums.BlockInfo
	isSolvingReorgChan chan int
}
type Engine interface {
	Loop(ctx context.Context) error
	HandleLog(handleLog func(e *types.Log))
	HandleReorg()
}

var isSolvingReorgChan = make(chan int)

func NewEngine(conf *config.Config,
	client ethClient,
	addresses []common.Address,
	topics [][]common.Hash,

) Engine {

	logChan := make(chan *types.Log)
	reorgChan := make(chan *enums.BlockInfo)
	emitter := emitter.NewEmitter(
		conf,
		client,
		addresses,
		topics,
		logChan,
		reorgChan,
		isSolvingReorgChan)

	return &engine{
		emitter:            emitter,
		logChan:            logChan,
		reorgChan:          reorgChan,
		isSolvingReorgChan: isSolvingReorgChan,
	}
}

func (e *engine) Loop(ctx context.Context) error {

	e.emitter.Loop(ctx)
	return nil

}

func (e *engine) HandleLog(handleLog func(e *types.Log)) {

	go func() {

		for {

			log := <-e.logChan

			fmt.Println("logggg---->", log)
			handleLog(log)
		}
	}()

}

func (e *engine) HandleReorg() {

	go func() {
		for reorg := range e.reorgChan {
			fmt.Println("reorg", reorg)

			fmt.Println("len --> reorg", len(isSolvingReorgChan))
			time.Sleep(10 * time.Second)
			isSolvingReorgChan <- 0
		}
	}()

}
