package emittertest

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/gsl/soyutils"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitter"
	"github.com/artnoi43/superwatcher/internal/poller"
	"github.com/artnoi43/superwatcher/pkg/components/mock"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

var (
	serviceConfigPath  = "../../../examples/demoservice/config/config.yaml"
	allCasesAlreadyRun = false

	flagCase    = flag.Int("case", -1, "Emitter test case")
	verboseFlag = flag.Bool("v", false, "Verbose emitter output")
)

// TODO: verbose does not work
func getFlagValues() (caseNumber int, verbose bool) {
	caseNumber = 1 // default case 1
	if flagCase != nil {
		caseNumber = *flagCase
	}
	if verboseFlag != nil {
		verbose = *verboseFlag
	}

	return caseNumber, verbose
}

// TestEmitterV2 uses TestCasesV2 to call emitterTestTemplateV2.
// V2 means that there are > 1 ReorgEvent for the test case.
func TestEmitterV2(t *testing.T) {
	for i := range TestCasesV2 {
		t.Run("TestEmitterV2", func(t *testing.T) {
			emitterTestTemplateV2(t, i+1)
		})
	}

	allCasesAlreadyRun = true
}

// Run this from the root of the repo with:
// go test -v ./internal/emitter -run TestEmitterByCase -case 69
// Go test binary already called `flag.Parse`, so we just simply
// need to name our flag so that the flag package knows to parse it too.
func TestEmitterByCase(t *testing.T) {
	if allCasesAlreadyRun {
		t.Skip("all cases were tested before -- skipping")
	}

	caseNumber, _ := getFlagValues()
	if caseNumber < 0 {
		TestEmitterV2(t)
		return
	}

	if len(TestCasesV2)+1 > caseNumber {
		testName := fmt.Sprintf("Case:%d", caseNumber)
		t.Run(testName, func(t *testing.T) {
			emitterTestTemplateV2(t, caseNumber)
		})

		return
	}

	t.Skipf("no such test case: %d", caseNumber)
}

func emitterTestTemplateV2(t *testing.T, caseNumber int) {
	tc := TestCasesV2[caseNumber-1]
	b, _ := json.Marshal(tc)
	t.Logf("testConfig for case %d: %s", caseNumber, b)

	type serviceConfig struct {
		SuperWatcherConfig *config.Config `yaml:"superwatcher_config" json:"superwatcherConfig"`
	}

	serviceConf, err := soyutils.ReadFileYAMLPointer[serviceConfig](serviceConfigPath)
	if err != nil {
		t.Fatal("bad config", err.Error())
	}
	// Override LoopInterval
	conf := serviceConf.SuperWatcherConfig
	conf.LoopInterval = 0
	conf.FilterRange = 10

	param := reorgsim.Param{
		StartBlock:    tc.FromBlock,
		ExitBlock:     tc.ToBlock + 200,
		BlockProgress: 10,
		Debug:         true,
	}

	var logsFiles = make([]string, len(tc.LogsFiles))
	for i, logsFile := range tc.LogsFiles {
		logsFiles[i] = "../" + logsFile
	}

	sim, err := reorgsim.NewReorgSimFromLogsFiles(param, tc.Events, logsFiles, "EmitterTestV2", 4)
	if err != nil {
		t.Fatal("error creating ReorgSim", err.Error())
	}

	errChan := make(chan error, 5)
	syncChan := make(chan struct{})
	filterResultChan := make(chan *superwatcher.FilterResult)

	fakeRedis := mock.NewDataGatewayMem(tc.FromBlock-1, true)
	testPoller := poller.New(nil, nil, conf.DoReorg, conf.FilterRange, sim, conf.LogLevel)
	testEmitter := emitter.New(conf, sim, fakeRedis, testPoller, syncChan, filterResultChan, errChan)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		testEmitter.Loop(ctx)
	}()

	go func() {
		if err := <-errChan; errors.Is(err, reorgsim.ErrExitBlockReached) {
			// This triggers shutdown on testEmitter, causing result from channels to be nil
			cancel()
		}
	}()

	movedHashes, logsDest := reorgsim.LogsFinalDst(tc.Events)

	var reached = make(map[common.Hash]bool)
	for {
		result := <-filterResultChan
		if result == nil {
			break
		}

		for _, b := range result.GoodBlocks {
			for _, log := range b.Logs {
				if !gslutils.Contains(movedHashes, log.TxHash) {
					continue
				}

				dest := logsDest[log.TxHash]

				if b.Number != dest {
					// Explicitly marked as false here
					reached[log.TxHash] = false
					continue
				}

				reached[log.TxHash] = true
			}
		}

		fakeRedis.SetLastRecordedBlock(nil, result.LastGoodBlock)
		syncChan <- struct{}{}
	}

	t.Log("reached", reached)
	for txHash := range reached {
		reachedDest, ok := reached[txHash]
		if !ok {
			continue
		}

		if !reachedDest {
			t.Errorf("log hash %s was not moved to %d", txHash.String(), logsDest[txHash])
		}
	}
}
