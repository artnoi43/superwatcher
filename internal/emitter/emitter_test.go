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
	"github.com/artnoi43/superwatcher/pkg/datagateway"
	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

type testConfig struct {
	Param     reorgsim.BaseParam    `json:"baseParam"`
	Events    []reorgsim.ReorgEvent `json:"reorgEvents"`
	FromBlock uint64                `json:"fromBlock"`
	ToBlock   uint64                `json:"toBlock"`
	LogsFiles []string              `json:"logs"`
}

var (
	logsPath  = "../../test_logs"
	testCases = []testConfig{
		{
			FromBlock: 15944400,
			ToBlock:   15944500,
			LogsFiles: []string{
				logsPath + "/logs_poolfactory.json",
				logsPath + "/logs_lp.json",
			},
			Param: reorgsim.BaseParam{
				StartBlock: 15944390,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgBlock: 15944411,
					MovedLogs:  nil,
				},
			},
		},
		{
			FromBlock: 15965717,
			ToBlock:   15965748,
			LogsFiles: []string{
				logsPath + "/logs_lp_2_1.json",
				logsPath + "/logs_lp_2_2.json",
			},
			Param: reorgsim.BaseParam{
				StartBlock: 15965710,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgBlock: 15965730,
					MovedLogs:  nil,
				},
			},
		},
		{
			FromBlock: 15965802,
			ToBlock:   15965835,
			LogsFiles: []string{
				logsPath + "/logs_lp_3_1.json",
				logsPath + "/logs_lp_3_2.json",
			},
			Param: reorgsim.BaseParam{
				StartBlock: 15965800,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgBlock: 15965811,
					MovedLogs:  nil,
				},
			},
		},
		{
			FromBlock: 15966460,
			ToBlock:   15966479,
			LogsFiles: []string{
				logsPath + "/logs_lp_4.json",
			},
			Param: reorgsim.BaseParam{
				StartBlock: 15966455,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgBlock: 15966475,
					MovedLogs:  nil,
				},
			},
		},
		{
			FromBlock: 15966500,
			ToBlock:   15966536,
			LogsFiles: []string{
				logsPath + "/logs_lp_5.json",
			},
			Param: reorgsim.BaseParam{
				StartBlock: 15966490,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgBlock: 15966536,
					MovedLogs:  nil,
				},
			},
		},
		{
			FromBlock: 15966500,
			ToBlock:   15966536,
			LogsFiles: []string{
				logsPath + "/logs_lp_5.json",
			},
			Param: reorgsim.BaseParam{
				StartBlock: 15966490,
			},
			Events: []reorgsim.ReorgEvent{
				{
					ReorgBlock: 15966512, // 0xf3a130
					// Move logs of 1 txHash to new block
					MovedLogs: map[uint64][]reorgsim.MoveLogs{
						15966522: { // 0xf3a13a
							{
								NewBlock: 15966527,
								TxHashes: []common.Hash{
									common.HexToHash("0x53f6b4200c700208fe7bb8cb806b0ce962a75e7a31d8a523fbc4affdc22ffc44"),
								},
							},
						},
						15966525: { // 0xf3a13d
							{
								NewBlock: 15966527, // 0xf3a13f
								TxHashes: []common.Hash{
									common.HexToHash("0xa46b7e3264f2c32789c4af8f58cb11293ac9a608fb335e9eb6f0fb08be370211"),
								},
							},
						},
					},
				},
			},
		},
	}
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
var allCasesAlreadyRun bool

func TestEmitterAllCases(t *testing.T) {
	allCasesAlreadyRun = true

	_, verbose := getFlagValues()
	for i := range testCases {
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
var (
	flagCase    = flag.Int("case", -1, "Emitter test case")
	verboseFlag = flag.Bool("v", false, "Verbose emitter output")
)

func TestEmitterByCase(t *testing.T) {
	if allCasesAlreadyRun {
		t.Skip("all cases were tested before -- skipping")
	}

	caseNumber, verbose := getFlagValues()
	if caseNumber < 0 {
		TestEmitterAllCases(t)
		return
	}

	if len(testCases)+1 > caseNumber {
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

	fakeRedis := datagateway.NewMock(tc.FromBlock-1, true)

	param := reorgsim.ParamV1{
		BaseParam: reorgsim.BaseParam{
			StartBlock:    tc.FromBlock,
			BlockProgress: 20,
			Debug:         true,
			ExitBlock:     tc.ToBlock + 200,
		},
		ReorgEvent: tc.Events[0],
	}

	sim, err := reorgsim.NewReorgSimV2FromLogsFiles(param.BaseParam, []reorgsim.ReorgEvent{param.ReorgEvent}, tc.LogsFiles, 2)
	if err != nil {
		t.Fatal("error creating ReorgSimV2", err.Error())
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
	latestGoodBlocks := make(map[uint64]*superwatcher.BlockInfo)
	movedToCount := make(map[common.Hash]bool)

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
		syncChan <- struct{}{}
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
