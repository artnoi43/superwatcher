package config

type EmitterConfig struct {
	// External dependencies
	NodeURL       string `yaml:"node_url" json:"nodeURL"`
	RedisConnAddr string `yaml:"redis_conn_addr" json:"redisConnAddr"`

	// Event log filtering config
	StartBlock    uint64 `yaml:"start_block" json:"startBlock"`
	FilterRange   uint64 `yaml:"filter_range" json:"filterRange"`
	GoBackRetries uint64 `yaml:"go_back_retries" json:"goBackRetries"`
	LoopInterval  uint64 `yaml:"loop_interval" json:"loopInterval"`
	LogLevel      uint8  `yaml:"log_level" json:"logLevel"`
}
