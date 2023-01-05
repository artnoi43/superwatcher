package config

import "github.com/artnoi43/superwatcher"

type Config struct {
	SuperWatcherConfig *superwatcher.Config `mapstructure:"superwatcher_config" yaml:"superwatcher_config" json:"superwatcherConfig"`
	RedisConnAddr      string               `mapstructure:"redis_conn_addr" yaml:"redis_conn_addr" json:"redisConnAddr"`
	Chain              string               `mapstructure:"chain" yaml:"chain" json:"chain"`
}
