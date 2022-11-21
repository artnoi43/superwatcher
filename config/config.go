package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	// External dependencies
	NodeURL       string `yaml:"node_url" json:"nodeURL"`
	RedisConnAddr string `yaml:"redis_conn_addr" json:"redisConnAddr"`

	// Event log filtering config
	Chain           string `yaml:"chain" json:"chain"`
	StartBlock      uint64 `yaml:"start_block" json:"startBlock"`
	LookBackBlocks  uint64 `yaml:"lookback_blocks" json:"lookBackBlock"`
	LookBackRetries uint64 `yaml:"lookback_retries" json:"lookBackRetries"`
	LoopInterval    uint64 `yaml:"loop_interval" json:"loop_interval"`
}

func ConfigYAML(filename string) (*Config, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	var conf Config
	if err := yaml.Unmarshal(b, &conf); err != nil {
		return nil, errors.Wrapf(err, "failed to parse YAML config from file %s", filename)
	}

	return &conf, nil
}
