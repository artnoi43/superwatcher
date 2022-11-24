package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type EmitterConfig struct {
	// External dependencies
	NodeURL       string `yaml:"node_url" json:"nodeURL"`
	RedisConnAddr string `yaml:"redis_conn_addr" json:"redisConnAddr"`

	// Event log filtering config
	Chain         string `yaml:"chain" json:"chain"`
	StartBlock    uint64 `yaml:"start_block" json:"startBlock"`
	FilterRange   uint64 `yaml:"filter_range" json:"filterRange"`
	GoBackRetries uint64 `yaml:"go_back_retries" json:"goBackRetries"`
	LoopInterval  uint64 `yaml:"loop_interval" json:"loopInterval"`
}

func ConfigYAML(filename string) (*EmitterConfig, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	var conf EmitterConfig
	if err := yaml.Unmarshal(b, &conf); err != nil {
		return nil, errors.Wrapf(err, "failed to parse YAML config from file %s", filename)
	}

	return &conf, nil
}
