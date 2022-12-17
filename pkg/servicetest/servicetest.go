package servicetest

import (
	"context"
	"errors"
	"sync"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/pkg/initsuperwatcher"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

// TestCase will be converted into config.EmitterConfig and reorgsim.Param
// to cerate TestComponents
type TestCase struct {
	StartBlock          uint64   `json:"startBlock"`          // Maps to EmitterConfig.StartBlock and is also used to init mock superwatcher.StateDataGateway
	ReorgBlock          uint64   `json:"reorgBlock"`          // Block that reorgSim will use as the reorging point
	ExitBlock           uint64   `json:"exitBlock"`           // Block for reorgSim to exit cleanly (in servicetest only)
	LogsFiles           []string `json:"logFiles"`            // JSON logs files for initializing reorgSim
	DataGatewayFirstRun bool     `json:"dataGatewayFirstRun"` // If set to true, the emitter will go back due to datagateway.ErrRecordNotFound
}

// TestComponents is used by RunServiceTestComponents to instantiate
// superwatcher.WatcherEmitter and superwatcher.WatcherEngine for RunService
type TestComponents struct {
	conf           *config.EmitterConfig
	client         superwatcher.EthClient
	serviceEngine  superwatcher.ServiceEngine
	dataGatewayGet superwatcher.GetStateDataGateway
	dataGatewaySet superwatcher.SetStateDataGateway
}

func InitTestComponents(
	conf *config.EmitterConfig,
	serviceEngine superwatcher.ServiceEngine,
	logsFullPaths []string,
	start uint64,
	reorgAt uint64,
	exit uint64,
	firstRun bool, // If true, then the mock datagateway will return `ErrRecordNotFound` until `SetLastRecordedBlock`` is called
) (
	*TestComponents,
	reorgsim.ParamV1, // For logging
) {
	param := reorgsim.ParamV1{
		BaseParam: reorgsim.BaseParam{
			StartBlock:    start,
			BlockProgress: 5,
			ExitBlock:     exit,
		},
		ReorgEvent: reorgsim.ReorgEvent{
			ReorgBlock: reorgAt,
		},
	}

	fakeRedis := datagateway.NewMock(conf.StartBlock, !firstRun)
	return &TestComponents{
		conf:           conf,
		client:         reorgsim.NewReorgSimFromLogsFiles(param, logsFullPaths, conf.LogLevel),
		serviceEngine:  serviceEngine,
		dataGatewayGet: fakeRedis,
		dataGatewaySet: fakeRedis,
	}, param
}

// RunServiceTestComponents runs the entire service using |components| and |param|.
// It does so by setting up superwatcher.WatcherEmitter and superwatcher.WatcherEngine
// and pass these objects to RunService.
// StateDataGateway is created within this function and will be returned to caller
func RunServiceTestComponents(components *TestComponents) (
	superwatcher.GetStateDataGateway, // Return this out so test code can read from dataGateway
	error,
) {
	// Use nil addresses and topics
	emitter, engine := initsuperwatcher.New(
		components.conf,
		components.client,
		components.dataGatewayGet,
		components.dataGatewaySet,
		nil,
		nil,
		components.serviceEngine,
	)

	return components.dataGatewayGet, RunService(emitter, engine)
}

// RunService executes the most basic emitter and engine logic, and returns an error from these components.
func RunService(emitter superwatcher.WatcherEmitter, engine superwatcher.WatcherEngine) error {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	var retErr error
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := emitter.Loop(ctx); err != nil {
			if errors.Is(err, reorgsim.ErrExitBlockReached) {
				cancel()
				emitter.Shutdown()
				return
			}

			retErr = err
		}
	}()

	if err := engine.Loop(ctx); err != nil {
		if errors.Is(err, reorgsim.ErrExitBlockReached) {
			return nil
		}

		retErr = err
	}

	wg.Wait()

	return retErr
}
