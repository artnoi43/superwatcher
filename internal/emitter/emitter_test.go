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
	"github.com/artnoi43/superwatcher/pkg/datagateway/watcherstate/mockwatcherstate"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

type testConfig struct {
	StartBlock uint64                         `json:"startBlock"`
	ReorgedAt  uint64                         `json:"reorgedAt"`
	FromBlock  uint64                         `json:"fromBlock"`
	ToBlock    uint64                         `json:"toBlock"`
	LogsFiles  []string                       `json:"logs"`
	MovedLogs  map[uint64][]reorgsim.MoveLogs `json:"movedLogs"`
}

var testCases = []testConfig{
	{
		StartBlock: 15944390,
		ReorgedAt:  15944411,
		FromBlock:  15944400,
		ToBlock:    15944500,
		LogsFiles: []string{
			"./assets/logs_poolfactory.json",
			"./assets/logs_lp.json",
		},
		MovedLogs: nil,
	},
	{
		StartBlock: 15965710,
		ReorgedAt:  15965730,
		FromBlock:  15965717,
		ToBlock:    15965748,
		LogsFiles: []string{
			"./assets/logs_lp_2_1.json",
			"./assets/logs_lp_2_2.json",
		},
		MovedLogs: nil,
	},
	{
		StartBlock: 15965800,
		ReorgedAt:  15965811,
		FromBlock:  15965802,
		ToBlock:    15965835,
		LogsFiles: []string{
			"./assets/logs_lp_3_1.json",
			"./assets/logs_lp_3_2.json",
		},
		MovedLogs: nil,
	},
	{
		StartBlock: 15966455,
		ReorgedAt:  15966475,
		FromBlock:  15966460,
		ToBlock:    15966479,
		LogsFiles: []string{
			"./assets/logs_lp_4.json",
		},
		MovedLogs: nil,
	},
	{
		StartBlock: 15966490,
		ReorgedAt:  15966536,
		FromBlock:  15966500,
		ToBlock:    15966536,
		LogsFiles: []string{
			"./assets/logs_lp_5.json",
		},
		MovedLogs: nil,
	},
	{
		StartBlock: 15966490,
		ReorgedAt:  15966530, // 0xf3a142
		FromBlock:  15966500,
		ToBlock:    15966540,
		LogsFiles: []string{
			"./assets/logs_lp_5.json",
		},
		// Move logs of 1 txHash to new block
		MovedLogs: map[uint64][]reorgsim.MoveLogs{
			15966532: { // 0xf3a144
				{
					NewBlock: 15966536,
					TxHashes: []common.Hash{
						common.HexToHash("0xf1ead2d704cd903038dfd75afd252b9b7928f5070e35550c8daa1a8c5c5941a7"),
						// common.HexToHash("0x41f48d4614c1e7333e545d3824b1ca6b19ef640fd335be990d50e4cf36b3a95d"),
					},
				},
			},
		},
	},
}

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
var allCasesAlreadyRun bool

func TestEmitterAllCases(t *testing.T) {
	allCasesAlreadyRun = true

	_, verbose := getFlagValues()
	for i := range testCases {
		testName := fmt.Sprintf("Case:%d", i+1)
		t.Run(testName, func(t *testing.T) {
			emitterTestTemplate(t, i+1, verbose)
		})
	}
}

// Run this from the root of the repo with:
// go test -v ./internal/emitter -run TestEmitterByCase -case 69
// Go test binary already called `flag.Parse`, so we just simply
// need to name our flag so that the flag package knows to parse it too.
var flagCase = flag.Int("case", -1, "Emitter test case")
var verboseFlag = flag.Bool("v", false, "Verbose emitter output")

func TestEmitterByCase(t *testing.T) {
	if allCasesAlreadyRun {
		t.Skip("all cases were tested before -- skipping")
	}

	caseNumber, verbose := getFlagValues()
	if caseNumber < 0 {
		TestEmitterAllCases(t)
		return
	}

	if len(testCases) > caseNumber {
		testName := fmt.Sprintf("Case:%d", caseNumber)
		t.Run(testName, func(t *testing.T) {
			emitterTestTemplate(t, caseNumber, verbose)
		})

		return
	}

	t.Skipf("no such test case: %d", caseNumber)
}

// emitterTestTemplate is designed to test emitter's full `Loop` with reorgsim mocked chain.
// The assumption test checks are only valid for logs filtered by reorgsim code,
// i.e. this test is NOT a proper test for REAL ethclient,
// as reorged logs may reappear on different block on a real chain.
func emitterTestTemplate(t *testing.T, caseNumber int, verbose bool) {
	tc := testCases[caseNumber-1]
	b, _ := json.Marshal(tc)
	t.Logf("testConfig for case %d: %s", caseNumber, b)

	type serviceConfig struct {
		SuperWatcherConfig *config.EmitterConfig `yaml:"superwatcher_config" json:"superwatcherConfig"`
	}

	serviceConf, err := soyutils.ReadFileYAMLPointer[serviceConfig]("../../superwatcher-demo/config/config.yaml")
	if err != nil {
		t.Fatal("bad config", err.Error())
	}
	// Override LoopInterval
	conf := serviceConf.SuperWatcherConfig
	conf.LoopInterval = 0

	fakeRedis := mockwatcherstate.New(tc.FromBlock - 1)

	param := reorgsim.Param{
		StartBlock:    tc.FromBlock,
		BlockProgress: 20,
		ReorgedBlock:  tc.ReorgedAt,
		ExitBlock:     tc.ToBlock + 200,
		Debug:         true,
	}

	sim := reorgsim.NewReorgSimFromLogsFiles(param, tc.LogsFiles, 2, tc.MovedLogs)

	// Buffered error channels, because if sim will die on ExitBlock, then it will die multiple times
	errChan := make(chan error, 5)
	syncChan := make(chan struct{})
	filterResultChan := make(chan *superwatcher.FilterResult)
	testEmitter := New(conf, sim, fakeRedis, nil, nil, syncChan, filterResultChan, errChan)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		testEmitter.(*emitter).Loop(ctx)
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
	var latestGoodBlocks = make(map[uint64]*superwatcher.BlockInfo)

	for {
		result := <-filterResultChan

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
			}

			_, ok := latestGoodBlocks[block.Number]
			if !ok {
				latestGoodBlocks[block.Number] = block
				continue
			}
		}

		fakeRedis.SetLastRecordedBlock(ctx, result.LastGoodBlock)
		syncChan <- struct{}{}
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
