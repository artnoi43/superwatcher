package config

type EmitterConfig struct {
	// External dependencies
	NodeURL string `mapstructure:"node_url" yaml:"node_url" json:"nodeURL"`

	// Event log filtering config
	StartBlock    uint64 `mapstructure:"start_block" yaml:"start_block" json:"startBlock"`
	FilterRange   uint64 `mapstructure:"filter_range" yaml:"filter_range" json:"filterRange"`
	GoBackRetries uint64 `mapstructure:"go_back_retries" yaml:"go_back_retries" json:"goBackRetries"`
	LoopInterval  uint64 `mapstructure:"loop_interval" yaml:"loop_interval" json:"loopInterval"`
	LogLevel      uint8  `mapstructure:"log_level" yaml:"log_level" json:"logLevel"`
}
