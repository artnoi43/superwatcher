package main

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/net/context"

	"github.com/artnoi43/superwatcher/superwatcher-demo/lib/utils"
)

type config struct {
	NodeURL string `yaml:"node_url"`

	Addresses []string `yaml:"addresses"`
	Topics    []string `yaml:"topics"`
	FromBlock int64    `yaml:"from_block"`
	ToBlock   int64    `yaml:"to_block"`
}

func main() {
	conf, err := utils.ReadFileYAML[config]("./target_logs.yaml")
	if err != nil {
		panic("read config failed: " + err.Error())
	}

	client, err := ethclient.Dial(conf.NodeURL)
	if err != nil {
		panic("new client failed: " + err.Error())
	}

	var addresses []common.Address
	for _, addrString := range conf.Addresses {
		address := common.HexToAddress(addrString)
		addresses = append(addresses, address)
	}
	fmt.Println("addresses", addresses)

	var topics []common.Hash
	for _, topicsStr := range conf.Topics {
		topic := common.HexToHash(topicsStr)
		topics = append(topics, topic)
	}
	fmt.Println("topics", topics)

	logs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(conf.FromBlock),
		ToBlock:   big.NewInt(conf.ToBlock),
		Addresses: addresses,
		Topics:    [][]common.Hash{topics},
	})
	if err != nil {
		panic("failed to filterLogs: " + err.Error())
	}

	logsJson, err := json.Marshal(logs)
	if err != nil {
		panic("failed to marshal event logs to json: " + err.Error())
	}

	fmt.Printf("%s\n", logsJson)
}
