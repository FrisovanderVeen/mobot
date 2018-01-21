package config

import (
	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Discord DiscordConf
	Youtube YoutubeConf
}

type DiscordConf struct {
	Token  string
	Prefix string
}

type YoutubeConf struct {
	Key string
}

// GetConfig decodes the config file
func GetConfig(configFile string) (*TomlConfig, error) {
	var config TomlConfig
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
