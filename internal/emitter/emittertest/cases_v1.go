package emittertest

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/artnoi43/superwatcher/pkg/reorgsim"
)

var TestCasesV1 = []TestConfig{
	{
		FromBlock: 15944400,
		ToBlock:   15944500,
		LogsFiles: []string{
			logsPath + "/logs_poolfactory.json",
			logsPath + "/logs_lp.json",
		},
		Param: reorgsim.Param{
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
		Param: reorgsim.Param{
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
		Param: reorgsim.Param{
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
		Param: reorgsim.Param{
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
		Param: reorgsim.Param{
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
		ToBlock:   15966540,
		LogsFiles: []string{
			logsPath + "/logs_lp_5.json",
		},
		Param: reorgsim.Param{
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
