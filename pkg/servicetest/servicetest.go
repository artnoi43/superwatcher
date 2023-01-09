package servicetest

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/pkg/components"
	"github.com/artnoi43/superwatcher/pkg/components/mock"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

// TestCase will be converted into config.Config and reorgsim.Param to create TestComponents
type TestCase struct {
	Param  reorgsim.Param        `json:"param"`
	Events []reorgsim.ReorgEvent `json:"reorgEvents"`
	// JSON logs files for initializing reorgSim
	LogsFiles []string `json:"logFiles"`
	// If set to true, the emitter will go back due to superwatcher.ErrRecordNotFound
	DataGatewayFirstRun bool `json:"dataGatewayFirstRun"`
	// EmitterPoller's poll level
	Policy superwatcher.Policy `json:"policy"`
}

// TestComponents is used by RunServiceTestComponents to instantiate
// superwatcher.Emitter and superwatcher.Engine for RunService
type TestComponents struct {
	conf           *superwatcher.Config
	client         superwatcher.EthClient
	serviceEngine  superwatcher.ServiceEngine
	dataGatewayGet superwatcher.GetStateDataGateway
	dataGatewaySet superwatcher.SetStateDataGateway
}

func DefaultServiceTestConfig(
	startBlock uint64,
	logLevel uint8,
	policy superwatcher.Policy,
) *superwatcher.Config {
	return &superwatcher.Config{
		// We use fakeRedis and fakeEthClient, so no need for token strings.
		StartBlock:       startBlock,
		DoReorg:          true,
		DoHeader:         true,
		FilterRange:      10,
		MaxGoBackRetries: 2,
		LoopInterval:     0,
		LogLevel:         logLevel,
		Policy:           policy,
	}
}

func InitTestComponents(
	conf *superwatcher.Config,
	serviceEngine superwatcher.ServiceEngine,
	param reorgsim.Param,
	events []reorgsim.ReorgEvent,
	logsFullPaths []string,
	firstRun bool, // If true, then the mock datagateway will return `ErrRecordNotFound` until `SetLastRecordedBlock`` is called
) *TestComponents {
	sim, err := reorgsim.NewReorgSimFromLogsFiles(param, events, logsFullPaths, "ServiceTest", conf.LogLevel)
	if err != nil {
		panic("failed to create ReorgSim")
	}

	fakeRedis := mock.NewDataGatewayMem(conf.StartBlock, !firstRun)

	return &TestComponents{
		conf:           conf,
		client:         sim,
		serviceEngine:  serviceEngine,
		dataGatewayGet: fakeRedis,
		dataGatewaySet: fakeRedis,
	}
}

// RunServiceTestComponents runs the entire service using |components| and |param|.
// It does so by setting up superwatcher.Emitter and superwatcher.Engine
// and pass these objects to RunService.
// StateDataGateway is created within this function and will be returned to caller
func RunServiceTestComponents(tc *TestComponents) (
	superwatcher.GetStateDataGateway, // Return this out so test code can read from dataGateway
	error,
) {
	emitter, engine := components.NewDefault(
		tc.conf,
		tc.client,
		tc.dataGatewayGet,
		tc.dataGatewaySet,
		tc.serviceEngine,
		// Use nil addresses and topics
		nil,
		nil,
	)

	return tc.dataGatewayGet, RunService(emitter, engine)
}

// RunService executes the most basic emitter and engine logic, and returns an error from these components.
func RunService(
	emitter superwatcher.Emitter,
	engine superwatcher.Engine,
) error {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		emitter.Poller().SetDoHeader(false)
	}()

	var wg sync.WaitGroup
	var retErr error
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := emitter.Loop(ctx); err != nil {
			if errors.Is(err, reorgsim.ErrExitBlockReached) {
				emitter.Shutdown()
				cancel()

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
