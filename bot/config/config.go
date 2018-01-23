package config

import (
	"github.com/BurntSushi/toml"
)

// TomlConfig is a decoding template for the config file
type TomlConfig struct {
	Discord DiscordConf
	Youtube YoutubeConf
	Plugins map[string]Plugin
}

// DiscordConf are all the discord related settings
type DiscordConf struct {
	Token  string
	Prefix string
}

// YoutubeConf are all the youtube related settings
type YoutubeConf struct {
	Key string
}

// Plugin are all the plugin related settings
type Plugin struct {
	Necessary bool
	Enabled   bool
}

// GetConfig decodes the config file
func GetConfig(configFile string) (*TomlConfig, error) {
	var config TomlConfig
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
