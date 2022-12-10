package config

import (
	superwatcherConfig "github.com/artnoi43/superwatcher/config"
)

type Config struct {
	SuperWatcherConfig *superwatcherConfig.EmitterConfig `mapstructure:"superwatcher_config" yaml:"superwatcher_config" json:"superwatcherConfig"`
	RedisConnAddr      string                            `mapstructure:"redis_conn_addr" yaml:"redis_conn_addr" json:"redisConnAddr"`
	Chain              string                            `mapstructure:"chain" yaml:"chain" json:"chain"`
}
