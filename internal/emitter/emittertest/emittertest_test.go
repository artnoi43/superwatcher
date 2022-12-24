package emittertest

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/gsl/soyutils"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitter"
	"github.com/artnoi43/superwatcher/internal/poller"
	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

// TestEmitterV2 uses TestCasesV2 to call emitterTestTemplateV2.
// V2 means that there are > 1 ReorgEvent for the test case.
func TestEmitterV2(t *testing.T) {
	for _, tc := range TestCasesV2 {
		emitterTestTemplateV2(t, tc)
	}
}

func emitterTestTemplateV2(t *testing.T, tc TestConfig) {
	type serviceConfig struct {
		SuperWatcherConfig *config.Config `yaml:"superwatcher_config" json:"superwatcherConfig"`
	}

	serviceConf, err := soyutils.ReadFileYAMLPointer[serviceConfig]("../../../superwatcher-demo/config/config.yaml")
	if err != nil {
		t.Fatal("bad config", err.Error())
	}
	// Override LoopInterval
	conf := serviceConf.SuperWatcherConfig
	conf.LoopInterval = 0
	conf.FilterRange = 10

	param := reorgsim.BaseParam{
		StartBlock:    tc.FromBlock,
		ExitBlock:     tc.ToBlock + 100,
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

	fakeRedis := datagateway.NewMock(tc.FromBlock-1, true)
	testPoller := poller.New(nil, nil, conf.DoReorg, conf.FilterRange, sim.FilterLogs, conf.LogLevel)
	testEmitter := emitter.New(conf, sim, fakeRedis, testPoller, syncChan, filterResultChan, errChan)

	movedHashes, logsDest := reorgsim.LogsFinalDst(tc.Events)

	var wg sync.WaitGroup
	wg.Add(1)
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

	var reached = make(map[common.Hash]bool)
	go func() {
		defer wg.Done()

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
	}()

	wg.Wait()

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
