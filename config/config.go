package config

// Config is superwatcher-wide configuration
type Config struct {
	// External dependencies
	NodeURL string `mapstructure:"node_url" yaml:"node_url" json:"nodeURL"`

	// StartBlock is the shortest block height the emitter will consider as base, usually a contract's genesis block
	StartBlock uint64 `mapstructure:"start_block" yaml:"start_block" json:"startBlock"`

	// FilterRange is the forward range (number of new blocks) each call to emitter.poller.poll will perform
	FilterRange uint64 `mapstructure:"filter_range" yaml:"filter_range" json:"filterRange"`

	// DoReorg specifies whether superwatcher superwatcher.EmitterPoller will process chain reorg for superwatcher.PollResult
	DoReorg bool `mapstructure:"do_reorg" yaml:"do_reorg" json:"doReorg"`

	// DoHeader specifies whether superwatcher.EmitterPoller should fetch block headers too
	DoHeader bool `mapstructure:"do_header" yaml:"do_header" json:"doHeader"`

	// MaxGoBackRetries is the maximum number of blocks the emitter will go back for. Once this is reached,
	// the emitter exits on error ErrMaxRetriesReached
	MaxGoBackRetries uint64 `mapstructure:"max_go_back_retries" yaml:"max_go_back_retries" json:"maxGoBackRetries"`

	// LoopInterval is the number of seconds the emitter sleeps after each call to emitter.poller.poll
	LoopInterval uint64 `mapstructure:"loop_interval" yaml:"loop_interval" json:"loopInterval"`

	// LogLevel for debugger.Debugger, the higher the more verbose
	LogLevel uint8 `mapstructure:"log_level" yaml:"log_level" json:"logLevel"`
}
