package main

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/alexflint/go-arg"
	"github.com/artnoi43/gsl/gslutils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/net/context"

	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/utils"
)

type config struct {
	ConfigFile string   `yaml:"config_file" arg:"-c,--config" placeholder:"FILE" help:"Config file to read"`
	NodeURL    string   `yaml:"node_url" arg:"-n,--node-url" placeholder:"NODE_URL" help:"HTTP or WS URL of an Ethereum node"`
	Addresses  []string `yaml:"addresses" arg:"-a,--addresses" placeholder:"ADDR [ADDR..]" help:"Contract address"`
	Topics     []string `yaml:"topics" arg:"-t,--topics" placeholder:"TOPICS [TOPIC..]" help:"Logs topics"`
	TxHashes   []string `yaml:"tx_hashes" arg:"-h,--tx-hashes" placeholder:"HASH [HASH..]" help:"Transaction hashes"`
	FromBlock  int64    `yaml:"from_block" arg:"-f,--from-block" placeholder:"FROM_BLOCK" help:"Filter from block"`
	ToBlock    int64    `yaml:"to_block" arg:"-t,--to-block" placeholder:"TO_BLOCK" help:"Filter to block"`
	Verbose    bool     `yaml:"-" arg:"-v,--verbose" help:"Verbose output (CLI arg only, placing this in config file won't work)"`
}

func main() {
	// Parse config
	argConf := new(config)
	arg.MustParse(argConf)

	// Read in config
	configFile := "./config.yaml"
	if len(argConf.ConfigFile) > 0 {
		configFile = argConf.ConfigFile
	}

	conf, err := utils.ReadFileYAML[config](configFile)
	if err != nil {
		panic("read config failed: " + err.Error())
	}

	// Overwrite config from file with CLI args
	conf = mergeConfig(conf, argConf)

	client, err := ethclient.Dial(conf.NodeURL)
	if err != nil {
		panic("new client failed: " + err.Error())
	}

	var addresses []common.Address
	for _, addrString := range conf.Addresses {
		address := common.HexToAddress(addrString)
		addresses = append(addresses, address)
	}

	var topics []common.Hash
	for _, topicsStr := range conf.Topics {
		topic := common.HexToHash(topicsStr)
		topics = append(topics, topic)
	}

	if conf.Verbose {
		fmt.Println("Filter addresses", addresses)
		fmt.Println("Filter topics", topics)
		fmt.Println("Filter txHashes", conf.TxHashes)
	}

	logs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(conf.FromBlock),
		ToBlock:   big.NewInt(conf.ToBlock),
		Addresses: addresses,
		Topics:    [][]common.Hash{topics},
	})

	// Collect logs ([]types.Log) into []*types.Log,
	// and filtering for TxHash if we got one from the user.
	var targetLogs []*types.Log
	if len(conf.TxHashes) > 0 {
		targetLogs = gslutils.CollectPointersIf(&logs, func(log types.Log) bool {
			return gslutils.Contains(conf.TxHashes, log.TxHash.String())
		})
	} else {
		targetLogs = gslutils.CollectPointers(&logs)
	}

	if err != nil {
		panic("failed to filterLogs: " + err.Error())
	}

	logsJson, err := json.Marshal(targetLogs)
	if err != nil {
		panic("failed to marshal event logs to json: " + err.Error())
	}

	fmt.Printf("%s\n", logsJson)
}

// mergeConfig returns a merged config from a and b
// mergeConfig only uses a field value from a ONLY if that field is null/zero-valued in b,
// i.e. it uses b to overwrite a.
func mergeConfig(a, b *config) *config {
	if b == nil {
		if a == nil {
			return nil
		}
		return a
	}

	out := new(config)

	// Verbosity can only be toggled in arg
	out.Verbose = b.Verbose

	if len(a.NodeURL) == 0 {
		if len(b.NodeURL) != 0 {
			out.NodeURL = b.NodeURL
		}
	} else {
		out.NodeURL = a.NodeURL
	}

	if len(b.Addresses) != 0 {
		out.Addresses = b.Addresses
	} else {
		out.Addresses = a.Addresses
	}

	if len(b.Topics) > 0 {
		out.Topics = b.Topics
	} else {
		out.Topics = a.Topics
	}

	if len(b.TxHashes) != 0 {
		out.TxHashes = b.TxHashes
	} else {
		out.TxHashes = a.TxHashes
	}

	if b.FromBlock != 0 {
		out.FromBlock = b.FromBlock
	} else {
		out.FromBlock = a.FromBlock
	}

	if b.ToBlock != 0 {
		out.ToBlock = b.ToBlock
	} else {
		out.ToBlock = a.ToBlock
	}

	return out
}
