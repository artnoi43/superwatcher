package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/components"
	"github.com/artnoi43/superwatcher/pkg/components/mock"
)

type service struct {
	ctx         context.Context
	dataGateway superwatcher.SetStateDataGateway
}

func (s *service) HandleFilterResult(result *superwatcher.PollerResult) error {
	fmt.Printf("Got results for %d - %d\n", result.FromBlock, result.LastGoodBlock)
	fmt.Printf("Good blocks: %d\n", len(result.GoodBlocks))
	fmt.Printf("Reorged blocks: %d\n", len(result.ReorgedBlocks))
	fmt.Println("========================")

	if err := s.dataGateway.SetLastRecordedBlock(context.Background(), result.LastGoodBlock); err != nil {
		return errors.Wrapf(err, "failed to set lastRecordedBlock: %s", err.Error())
	}

	fmt.Println("Wrote lastRecordedBlock", result.LastGoodBlock)
	// Syncs with Emitter

	return nil
}

func (s *service) HandleEmitterError(err error) error {
	fmt.Printf("Got superwatcher Emitter error: %s", err.Error())

	return err
}

func main() {
	nodeURL := os.Args[1]
	if len(nodeURL) == 0 {
		panic("empty eth node URL")
	}

	conf := &superwatcher.Config{
		NodeURL:          nodeURL,
		StartBlock:       16720532,
		FilterRange:      10,
		DoReorg:          true,
		DoHeader:         true,
		MaxGoBackRetries: 2,
		LoopInterval:     2,
		LogLevel:         2,
		Policy:           superwatcher.PolicyExpensive,
	}

	dataGateway := mock.NewDataGatewayMem(conf.StartBlock, true)
	uniswapV3Addr := []common.Address{common.HexToAddress("0x5777d92f208679DB4b9778590Fa3CAB3aC9e2168")}

	ctx := context.Background()
	s := &service{ctx: ctx, dataGateway: dataGateway}

	emitter, engine := components.NewThinEngineWithEmitter(
		conf,
		dataGateway, dataGateway,
		uniswapV3Addr, nil,
		superwatcher.NewEthClient(ctx, nodeURL),
		conf.Policy,
		s,
	)

	go func() {
		defer emitter.Shutdown()

		if err := emitter.Loop(ctx); err != nil {
			log.Println("error from emitter:", err.Error())
		}
	}()

	if err := engine.Loop(ctx); err != nil {
		log.Println("error from engine:", err.Error())
	}
}
