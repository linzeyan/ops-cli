/*
Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

/* Read config.toml. */
type ConfigBlock string

func (c ConfigBlock) String() string {
	return string(c)
}

/* Read config by Viper and return key/value to map struct. */
type readConfig struct {
	/* Specify the config path. */
	path string
	/* Specify the config type. */
	format string
	/* Specify the config field. */
	table ConfigBlock
	/* Return the key/value to mapping. */
	value map[string]any
}

func (r *readConfig) get() (map[string]any, error) {
	viper.SetConfigFile(r.path)
	viper.SetConfigType(r.format)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	all := viper.AllSettings()
	values, ok := all[r.table.String()]
	if !ok {
		return nil, ErrConfigTable
	}
	r.value, ok = values.(map[string]any)
	if !ok {
		return nil, ErrConfigContent
	}
	return r.value, nil
}

/* Get secret token and other settings from config. */
func Config(path string, table ConfigBlock) map[string]any {
	r := &readConfig{path: path, format: "toml", table: table}
	v, err := r.get()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return v
}
