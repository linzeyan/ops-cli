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
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

/* Read config by Viper and return key/value to map struct. */
type readConfig struct {
	/* Specify the config path. */
	path string
	/* Specify the config field. */
	table string
	/* Return the key/value to mapping. */
	value map[string]any
}

func (r *readConfig) get() (map[string]any, error) {
	typ := []string{JSONFormat, TomlFormat, YamlFormat}
	ext := strings.Replace(filepath.Ext(r.path), ".", "", 1)
	for _, v := range typ {
		if ext == v {
			viper.SetConfigType(ext)
			break
		}
		viper.SetConfigType(typ[1])
	}
	viper.SetConfigFile(r.path)
	if err := viper.ReadInConfig(); err != nil {
		stdLogger.Log.Debug(err.Error())
		return nil, err
	}
	all := viper.AllSettings()
	values, ok := all[r.table]
	if !ok {
		stdLogger.Log.Debug(ErrConfigTable.Error(), DefaultField(r.table))
		return nil, ErrConfigTable
	}
	r.value, ok = values.(map[string]any)
	if !ok {
		stdLogger.Log.Debug(ErrConfigContent.Error())
		return nil, ErrConfigContent
	}
	return r.value, nil
}

/* Get secret token and other settings from config. */
func Config(path, table string) map[string]any {
	r := &readConfig{path: path, table: table}
	v, err := r.get()
	if err != nil {
		stdLogger.Log.Fatal(err.Error(), NewField("path", path), NewField("table", table))
	}
	return v
}
