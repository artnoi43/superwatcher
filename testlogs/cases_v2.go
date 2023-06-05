package testlogs

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/soyart/superwatcher/pkg/reorgsim"
)

var TestCasesV2 = []*TestConfig{
	{
		FromBlock: 15944400,
		ToBlock:   15944500,
		LogsFiles: []string{
			"/logs_poolfactory.json",
			"/logs_lp.json",
		},
		Param: reorgsim.Param{
			StartBlock: 15944390,
		},
		Events: []reorgsim.ReorgEvent{
			{
				ReorgTrigger: 15944430,
				ReorgBlock:   15944419,
				MovedLogs: map[uint64][]reorgsim.MoveLogs{
					15944419: {
						{
							NewBlock: 15944455,
							TxHashes: []common.Hash{
								common.HexToHash("0x41e29790ca68bd08062b3ed1670216d49ef458dbe398fafda3918f2322c18068"),
							},
						},
					},
					15944455: {
						{
							NewBlock: 15944498,
							TxHashes: []common.Hash{
								common.HexToHash("0x620be69b041f986127322985854d3bc785abe1dc9f4df49173409f15b7515164"),
							},
						},
					},
				},
			},
		},
	},
	{
		FromBlock: 15966500,
		ToBlock:   15966540,
		LogsFiles: []string{
			"/logs_lp_5.json",
		},
		Param: reorgsim.Param{
			StartBlock: 15966490,
		},
		Events: []reorgsim.ReorgEvent{
			{
				ReorgTrigger: 15966515,
				ReorgBlock:   15966512, // 0xf3a130
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
				},
			},
			{
				ReorgTrigger: 15966525,
				ReorgBlock:   15966524,
				MovedLogs: map[uint64][]reorgsim.MoveLogs{
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
	{
		FromBlock: 15966400,
		ToBlock:   15966500,
		LogsFiles: []string{
			"/logs_lp_4.json",
		},
		Param: reorgsim.Param{
			StartBlock: 15966465,
		},
		Events: []reorgsim.ReorgEvent{
			{
				ReorgTrigger: 15966466,
				ReorgBlock:   15966464,
				MovedLogs: map[uint64][]reorgsim.MoveLogs{
					15966464: {
						{
							NewBlock: 15966477,
							TxHashes: []common.Hash{
								common.HexToHash("0x6e23009a3f85e85fa9602382a29f74f4ca3027f4ed5c48fee3229b53c4c51e7d"),
							},
						},
					},
				},
			},
			{
				ReorgBlock: 15966473,
				MovedLogs: map[uint64][]reorgsim.MoveLogs{
					15966473: {
						{
							NewBlock: 15966475,
							TxHashes: []common.Hash{
								common.HexToHash("0x0d6ae93dc766ae8be6ae975a5a42cffe7ad6327e5cf1dc20c50d60a07f849f14"),
							},
						},
					},
				},
			},
		},
	},
}
