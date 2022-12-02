package config

import (
	superwatcherConfig "github.com/artnoi43/superwatcher/config"
)

type Config struct {
	SuperWatcherConfig *superwatcherConfig.EmitterConfig `yaml:"superwatcher_config" json:"superwatcherConfig"`
	Chain              string                            `yaml:"chain" json:"chain"`
}
