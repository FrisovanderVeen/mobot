package bot

import (
	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Discord discordConf
}

type discordConf struct {
	Token  string
	Prefix string
}

// getConfig decodes the config file
func getConfig(configFile string) (*tomlConfig, error) {
	var config tomlConfig
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
