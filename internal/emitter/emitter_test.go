package emitter

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"testing"

	"github.com/artnoi43/gsl/gslutils"
	"github.com/artnoi43/gsl/soyutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/artnoi43/superwatcher"
	"github.com/artnoi43/superwatcher/config"
	"github.com/artnoi43/superwatcher/internal/emitter/emittertest"
	"github.com/artnoi43/superwatcher/internal/poller"
	"github.com/artnoi43/superwatcher/pkg/components/mock"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

var (
	flagCase           = flag.Int("case", -1, "Emitter test case")
	verboseFlag        = flag.Bool("v", false, "Verbose emitter output")
	testLogsPath       = "../../test_logs"
	serviceConfigPath  = "../../examples/demoservice/config/config.yaml"
	allCasesAlreadyRun = false
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

// allCasesAlreadyRun is used to skip TestEmitterByCase if TestEmitterAllCases were run.

func TestEmitterAllCases(t *testing.T) {
	allCasesAlreadyRun = true

	_, verbose := getFlagValues()
	for i := range emittertest.TestCasesV1 {
		testName := fmt.Sprintf("Case:%d", i+1)
		t.Run(testName, func(t *testing.T) {
			emitterTestTemplateV1(t, i+1, verbose)
		})
	}
}

// Run this from the root of the repo with:
// go test -v ./internal/emitter -run TestEmitterByCase -case 69
// Go test binary already called `flag.Parse`, so we just simply
// need to name our flag so that the flag package knows to parse it too.
func TestEmitterByCase(t *testing.T) {
	if allCasesAlreadyRun {
		t.Skip("all cases were tested before -- skipping")
	}

	caseNumber, verbose := getFlagValues()
	if caseNumber < 0 {
		TestEmitterAllCases(t)
		return
	}

	if len(emittertest.TestCasesV1)+1 > caseNumber {
		testName := fmt.Sprintf("Case:%d", caseNumber)
		t.Run(testName, func(t *testing.T) {
			emitterTestTemplateV1(t, caseNumber, verbose)
		})

		return
	}

	t.Skipf("no such test case: %d", caseNumber)
}

// emitterTestTemplateV1 is designed to test emitter's full `Loop` with ReorgSimV1 mocked chain.
// This means that the test chain will only have 1 reorg event.
func emitterTestTemplateV1(t *testing.T, caseNumber int, verbose bool) {
	tc := emittertest.TestCasesV1[caseNumber-1]
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

	fakeRedis := mock.NewDataGatewayMem(tc.FromBlock-1, true)

	param := reorgsim.BaseParam{
		StartBlock:    tc.FromBlock,
		BlockProgress: 20,
		Debug:         true,
		ExitBlock:     tc.ToBlock + 200,
	}
	events := []reorgsim.ReorgEvent{tc.Events[0]}

	sim, err := reorgsim.NewReorgSimFromLogsFiles(param, events, tc.LogsFiles, "EmitterTestV1", 4)
	if err != nil {
		t.Fatal("error creating ReorgSim", err.Error())
	}

	// Collect MovedLogs info
	var movedFromBlocks []uint64
	var movedToBlocks []uint64
	var movedTxHashes []common.Hash
	for movedFromBlock, moves := range tc.Events[0].MovedLogs {
		movedFromBlocks = append(movedFromBlocks, movedFromBlock)
		for _, move := range moves {
			movedToBlocks = append(movedToBlocks, move.NewBlock)
			movedTxHashes = append(movedTxHashes, move.TxHashes...)
		}
	}

	// testPoller got nil addresses and topics so it will poll logs from all addresses and topics
	testPoller := poller.New(nil, nil, conf.DoReorg, conf.FilterRange, sim, conf.LogLevel)

	// Buffered error channels, because if sim will die on ExitBlock, then it will die multiple times
	errChan := make(chan error, 5)
	syncChan := make(chan struct{})
	filterResultChan := make(chan *superwatcher.FilterResult)
	testEmitter := New(conf, sim, fakeRedis, testPoller, syncChan, filterResultChan, errChan)

	// Check if the emitter noticed a reorg
	var reorgedOnce bool

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		testEmitter.Loop(ctx)
	}()

	go func() {
		if err := <-errChan; err != nil {
			if errors.Is(err, reorgsim.ErrExitBlockReached) {
				// This triggers shutdown on testEmitter, causing result from channels to be nil
				cancel()
			}
		}
	}()

	var seenLogs []*types.Log
	latestGoodBlocks := make(map[uint64]*superwatcher.BlockInfo)
	movedToCount := make(map[common.Hash]bool)

	var prevResult *superwatcher.FilterResult
	for {
		result := <-filterResultChan

		if prevResult != nil {
			if result.LastGoodBlock <= prevResult.ToBlock {
				reorgedOnce = true
			}
		}

		if result == nil {
			break
		}
		if result.LastGoodBlock > tc.ToBlock {
			t.Logf("finished case %d", caseNumber)
			cancel()
			break
		}

		lastGoodBlock := result.LastGoodBlock

		for i, block := range result.ReorgedBlocks {
			reorgedOnce = true

			blockNumber := block.Number

			// Check LastGoodBlock
			if lastGoodBlock > blockNumber {
				t.Fatalf(
					"invalid LastGoodBlock: ReorgedBlocks[%d] blockNumber=%v, LastGoodBlock=%v",
					i, blockNumber, lastGoodBlock,
				)
			}

			b, ok := latestGoodBlocks[block.Number]
			if !ok {
				t.Fatalf("reorged block not in latestGoodBlocks - %d %s", block.Number, block.String())
			}
			if block.Hash != b.Hash {
				t.Fatalf("reorged block hash not seen before - %d %s", block.Number, block.String())
			}

			// Check that all the reorged logs were seen before in |seenLogs|
			for _, log := range block.Logs {
				if !gslutils.Contains(seenLogs, log) {
					fatalBadLog(t, "reorgedLog not seen before", log)
				}
			}
		}

		for _, block := range result.GoodBlocks {
			for _, log := range block.Logs {
				var ok bool

				// We should only see a good log once
				seenLogs, ok = appendUnique(seenLogs, log)
				if !ok {
					fatalBadLog(t, "duplicate good log in seenLogs", log)
				}

				// If the block is one of the movedFromBlocks, then it's not supposed to have any logs with movedTxHashes
				if gslutils.Contains(movedFromBlocks, block.Number) {
					if gslutils.Contains(movedTxHashes, log.TxHash) {
						// t.Log("movedBlock from", log.BlockNumber, log.BlockHash.String(), log.TxHash.String())
						fatalBadLog(t, "log was supposed to be removed from this block", log)
					}
				}

				// If the block is NOT one of the movedToBlocks, then it's not supposed to have any logs with movedTxHashes
				if !gslutils.Contains(movedToBlocks, block.Number) {
					if gslutils.Contains(movedTxHashes, log.TxHash) {
						// t.Log("movedBlock to", log.BlockNumber, log.BlockHash.String(), log.TxHash.String())
						fatalBadLog(t, "log was supposed to be moved to this block", log)
					}
				} else {
					movedToCount[log.TxHash] = true
				}
			}

			_, ok := latestGoodBlocks[block.Number]
			if !ok {
				latestGoodBlocks[block.Number] = block
				continue
			}
		}

		fakeRedis.SetLastRecordedBlock(ctx, result.LastGoodBlock)
		prevResult = result

		syncChan <- struct{}{}
	}

	if len(tc.Events) != 0 {
		if !reorgedOnce {
			t.Fatal("not reorgedOnce")
		}
	}

	for _, txHash := range movedTxHashes {
		if !movedToCount[txHash] {
			t.Errorf("movedToTxHash %s was not tagged true", txHash.String())
		}

		t.Log("movedLog", txHash.String())
	}
}

func fatalBadLog(t *testing.T, msg string, log *types.Log) {
	t.Fatalf(
		"%s: blockNumber: %d, blockHash %s, txHash %s, addr %s, topics[0]: %s",
		msg, log.BlockNumber, log.BlockHash.String(), log.TxHash.String(), log.Address.String(), log.Topics[0].String(),
	)
}

// appendUnique appends item to arr if arr does not contain item.
func appendUnique[T comparable](arr []T, item T) ([]T, bool) {
	if !gslutils.Contains(arr, item) {
		return append(arr, item), true
	}

	return arr, false
}
