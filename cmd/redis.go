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

package cmd

import (
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func init() {
	var redisFlag RedisFlag
	var redisCmd = &cobra.Command{
		Use:   CommandRedis,
		Short: "Opens a connection to a Redis server",
		RunE: func(_ *cobra.Command, args []string) error {
			if err := ReadConfig(CommandRedis, &redisFlag); err != nil {
				return err
			}
			return redisFlag.Do(args)
		},
	}

	rootCmd.AddCommand(redisCmd)
	redisCmd.Flags().StringVarP(&redisFlag.Username, "user", "u", "", common.Usage("Username to authenticate the current connection"))
	redisCmd.Flags().StringVarP(&redisFlag.Password, "auth", "a", "", common.Usage("Password to use when connecting to the server"))
	redisCmd.Flags().StringVarP(&redisFlag.Host, "hostname", "H", "127.0.0.1", common.Usage("Server hostname"))
	redisCmd.Flags().StringVarP(&redisFlag.Port, "port", "p", "6379", common.Usage("Server port"))
	redisCmd.Flags().IntVarP(&redisFlag.DB, "db", "n", 0, common.Usage("Database number"))
}

type RedisFlag struct {
	Username string `json:"user"`
	Password string `json:"auth"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DB       int    `json:"db"`
}

func (r *RedisFlag) Connection() *redis.Client {
	if validator.ValidIP(r.Host) || validator.ValidDomain(r.Host) {
		return redis.NewClient(&redis.Options{
			Username: r.Username,
			Password: r.Password,
			Addr:     r.Host + ":" + r.Port,
			DB:       r.DB,
		})
	}
	return nil
}

func (r *RedisFlag) Do(commands []string) error {
	var err error
	if len(commands) == 0 {
		return err
	}
	var args []string
	if len(commands) == 1 {
		args = strings.Fields(commands[0])
	} else {
		args = commands
	}
	var arg []any
	for _, v := range args {
		arg = append(arg, v)
	}
	rdb := r.Connection()
	if rdb == nil {
		return common.ErrInvalidArg
	}
	cmd := rdb.Do(common.Context, arg...)
	out, err := cmd.Result()
	if err != nil {
		return err
	}
	PrintJSON(out)
	return err
}
