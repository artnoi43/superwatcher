package config

type EmitterConfig struct {
	// External dependencies
	NodeURL string `mapstructure:"node_url" yaml:"node_url" json:"nodeURL"`

	// StartBlock is the shortest block height the emitter will consider as base, usually a contract's genesis block
	StartBlock uint64 `mapstructure:"start_block" yaml:"start_block" json:"startBlock"`

	// FilterRange is the forward range (number of new blocks) each filterLogs loop will perform
	FilterRange uint64 `mapstructure:"filter_range" yaml:"filter_range" json:"filterRange"`

	// MaxGoBackRetries is the maximum number of blocks the emitter will go back for. Once this is reached,
	// the emitter exits on error ErrMaxRetriesReached
	MaxGoBackRetries uint64 `mapstructure:"max_go_back_retries" yaml:"max_go_back_retries" json:"maxGoBackRetries"`

	// LoopInterval is the number of seconds the emitter sleeps after each call to filterLogs
	LoopInterval uint64 `mapstructure:"loop_interval" yaml:"loop_interval" json:"loopInterval"`

	// LogLevel for debugger.Debugger, the higher the more verbose
	LogLevel uint8 `mapstructure:"log_level" yaml:"log_level" json:"logLevel"`
}
