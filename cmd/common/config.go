package common

import (
	"github.com/spf13/viper"
)

/* Read config.toml. */
type ConfigBlock string

const (
	configType ConfigBlock = "toml"
	Discord    ConfigBlock = "discord"
	ICP        ConfigBlock = "west"
	LINE       ConfigBlock = "line"
	Slack      ConfigBlock = "slack"
	Telegram   ConfigBlock = "telegram"
)

func (c ConfigBlock) String() string {
	return string(c)
}

type readConfig struct {
	table ConfigBlock
	value map[string]interface{}
}

func (r *readConfig) get(config string) (map[string]interface{}, error) {
	viper.SetConfigFile(config)
	viper.SetConfigType(configType.String())
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	all := viper.AllSettings()
	values, ok := all[r.table.String()]
	if !ok {
		return nil, ErrConfigTable
	}
	r.value, ok = values.(map[string]interface{})
	if !ok {
		return nil, ErrConfigContent
	}
	return r.value, nil
}

/* Get secret token from config. */
func Config(config string, t ConfigBlock) (map[string]interface{}, error) {
	r := readConfig{table: t}
	return r.get(config)
}
