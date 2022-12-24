package servicetest

import (
	"context"
	"errors"
	"sync"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/pkg/components"
	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

// TestCase will be converted into config.Config and reorgsim.BaseParam to create TestComponents
type TestCase struct {
	Param  reorgsim.BaseParam    `json:"baseParam"`
	Events []reorgsim.ReorgEvent `json:"reorgEvents"`
	// JSON logs files for initializing reorgSim
	LogsFiles []string `json:"logFiles"`
	// If set to true, the emitter will go back due to datagateway.ErrRecordNotFound
	DataGatewayFirstRun bool `json:"dataGatewayFirstRun"`
}

// TestComponents is used by RunServiceTestComponents to instantiate
// superwatcher.WatcherEmitter and superwatcher.WatcherEngine for RunService
type TestComponents struct {
	conf           *config.Config
	client         superwatcher.EthClient
	serviceEngine  superwatcher.ServiceEngine
	dataGatewayGet superwatcher.GetStateDataGateway
	dataGatewaySet superwatcher.SetStateDataGateway
}

func DefaultServiceTestConfig(startBlock uint64, logLevel uint8) *config.Config {
	return &config.Config{
		// We use fakeRedis and fakeEthClient, so no need for token strings.
		StartBlock:       startBlock,
		DoReorg:          true,
		FilterRange:      10,
		MaxGoBackRetries: 2,
		LoopInterval:     0,
		LogLevel:         logLevel,
	}
}

func InitTestComponents(
	conf *config.Config,
	serviceEngine superwatcher.ServiceEngine,
	param reorgsim.BaseParam,
	events []reorgsim.ReorgEvent,
	logsFullPaths []string,
	firstRun bool, // If true, then the mock datagateway will return `ErrRecordNotFound` until `SetLastRecordedBlock`` is called
) *TestComponents {
	sim, err := reorgsim.NewReorgSimFromLogsFiles(param, events, logsFullPaths, "ServiceTest", conf.LogLevel)
	if err != nil {
		panic("failed to create ReorgSim")
	}

	fakeRedis := datagateway.NewMock(conf.StartBlock, !firstRun)

	return &TestComponents{
		conf:           conf,
		client:         sim,
		serviceEngine:  serviceEngine,
		dataGatewayGet: fakeRedis,
		dataGatewaySet: fakeRedis,
	}
}

// RunServiceTestComponents runs the entire service using |components| and |param|.
// It does so by setting up superwatcher.WatcherEmitter and superwatcher.WatcherEngine
// and pass these objects to RunService.
// StateDataGateway is created within this function and will be returned to caller
func RunServiceTestComponents(testComponents *TestComponents) (
	superwatcher.GetStateDataGateway, // Return this out so test code can read from dataGateway
	error,
) {
	// Use nil addresses and topics
	emitter, engine := components.NewDefault(
		testComponents.conf,
		testComponents.client,
		testComponents.dataGatewayGet,
		testComponents.dataGatewaySet,
		testComponents.serviceEngine,
		nil,
		nil,
	)

	return testComponents.dataGatewayGet, RunService(emitter, engine)
}

// RunService executes the most basic emitter and engine logic, and returns an error from these components.
func RunService(emitter superwatcher.Emitter, engine superwatcher.Engine) error {
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
