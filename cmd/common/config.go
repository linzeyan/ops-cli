package common

import (
	"log"
	"os"

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
	config string
	format string
	table  ConfigBlock
	value  map[string]interface{}
}

func (r *readConfig) get() (map[string]interface{}, error) {
	viper.SetConfigFile(r.config)
	viper.SetConfigType(r.format)
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
func Config(config string, table ConfigBlock) map[string]interface{} {
	r := &readConfig{config: config, format: "toml", table: table}
	v, err := r.get()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return v
}
