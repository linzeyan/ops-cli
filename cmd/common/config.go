package common

import (
	"errors"

	"github.com/spf13/viper"
)

/* Read config.toml. */
type ConfigBlock string

const (
	Discord  ConfigBlock = "discord"
	ICP      ConfigBlock = "west"
	LINE     ConfigBlock = "line"
	Slack    ConfigBlock = "slack"
	Telegram ConfigBlock = "telegram"
)

func (c ConfigBlock) String() string {
	return string(c)
}

type readConfig struct {
	table ConfigBlock
}

func (c readConfig) get(config string) (map[string]interface{}, error) {
	viper.SetConfigFile(config)
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	all := viper.AllSettings()
	values, ok := all[c.table.String()]
	if !ok {
		return nil, errors.New("table not found in config")
	}
	v, ok := values.(map[string]interface{})
	if !ok {
		return nil, errors.New("config content is incorrect")
	}
	return v, nil
}

/* Get secret token from config. */
func Config(config string, t ConfigBlock) (map[string]interface{}, error) {
	r := readConfig{table: t}
	return r.get(config)
}
