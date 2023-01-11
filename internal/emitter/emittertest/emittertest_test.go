package emittertest

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/artnoi43/gsl"
	"github.com/artnoi43/gsl/soyutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/internal/emitter"
	"github.com/artnoi43/superwatcher/internal/poller"
	"github.com/artnoi43/superwatcher/pkg/components/mock"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
	"github.com/artnoi43/superwatcher/pkg/testutils"
)

var (
	serviceConfigPath = "../../../examples/demoservice/config/config.yaml"
)

// TestEmitterV2 uses TestCasesV2 to call testEmitterV2.
// V2 means that there are > 1 ReorgEvent for the test case.
func TestEmitterV2(t *testing.T) {
	testutils.RunTestCase(t, "testEmitterV2", TestCasesV2, testEmitterV2)
}

func testEmitterV2(t *testing.T, caseNumber int) error {
	for _, policy := range []superwatcher.Policy{
		superwatcher.PolicyFast,
		superwatcher.PolicyNormal,
		// superwatcher.PolicyExpensive,
	} {
		tc := TestCasesV2[caseNumber-1]
		b, _ := json.Marshal(tc)
		t.Logf("testConfig for case %d: %s (policy %s)", caseNumber, b, policy.String())

		type serviceConfig struct {
			SuperWatcherConfig *superwatcher.Config `yaml:"superwatcher_config" json:"superwatcherConfig"`
		}

		serviceConf, err := soyutils.ReadFileYAMLPointer[serviceConfig](serviceConfigPath)
		if err != nil {
			return errors.Wrap(err, "bad YAML config")
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

		logsFiles := make([]string, len(tc.LogsFiles))
		for i, logsFile := range tc.LogsFiles {
			logsFiles[i] = "../" + logsFile
		}

		sim, err := reorgsim.NewReorgSimFromLogsFiles(param, tc.Events, logsFiles, "EmitterTestV2", 4)
		if err != nil {
			return errors.Wrap(err, "error creating ReorgSim")
		}

		errChan := make(chan error, 5)
		syncChan := make(chan struct{})
		pollResultChan := make(chan *superwatcher.PollerResult)

		fakeRedis := mock.NewDataGatewayMem(tc.FromBlock-1, true)
		testPoller := poller.New(nil, nil, conf.DoReorg, conf.DoHeader, conf.FilterRange, sim, conf.LogLevel, policy)
		testEmitter := emitter.New(conf, sim, fakeRedis, testPoller, syncChan, pollResultChan, errChan)

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

		movedHashes, _, logsDest := reorgsim.LogsReorgPaths(tc.Events)

		reached := make(map[common.Hash]bool)
		for {
			result := <-pollResultChan
			if result == nil {
				break
			}

			for _, b := range result.GoodBlocks {
				for _, log := range b.Logs {
					if !gsl.Contains(movedHashes, log.TxHash) {
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

	return nil
}
