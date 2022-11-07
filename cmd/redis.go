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
	"net"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initRedis() *cobra.Command {
	var flags struct {
		Username string `json:"user"`
		Password string `json:"auth"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		DB       int    `json:"db"`
	}
	var redisCmd = &cobra.Command{
		Use:   CommandRedis,
		Short: "Opens a connection to a Redis server",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(_ *cobra.Command, args []string) {
			var r Redis
			if rootConfig != "" {
				if err := ReadConfig(CommandRedis, &flags); err != nil {
					logger.Info(err.Error())
					return
				}
			}
			conn := r.Connection(flags.Host, flags.Port, flags.Username, flags.Password, flags.DB)
			if conn == nil {
				logger.Info(common.ErrFailedInitial.Error())
				return
			}
			if err := r.Do(conn, args); err != nil {
				logger.Info(err.Error())
			}
		},
	}

	redisCmd.Flags().StringVarP(&flags.Username, "user", "u", "", common.Usage("Username to authenticate the current connection"))
	redisCmd.Flags().StringVarP(&flags.Password, "auth", "a", "", common.Usage("Password to use when connecting to the server"))
	redisCmd.Flags().StringVarP(&flags.Host, "hostname", "h", "127.0.0.1", common.Usage("Server hostname"))
	redisCmd.Flags().StringVarP(&flags.Port, "port", "p", "6379", common.Usage("Server port"))
	redisCmd.Flags().IntVarP(&flags.DB, "db", "n", 0, common.Usage("Database number"))
	return redisCmd
}

type Redis struct{}

func (r *Redis) Connection(host, port, user, pass string, db int) *redis.Client {
	if common.IsIP(host) || common.IsDomain(host) {
		return redis.NewClient(&redis.Options{
			Username: user,
			Password: pass,
			Addr:     net.JoinHostPort(host, port),
			DB:       db,
		})
	}
	logger.Debug(common.ErrInvalidArg.Error(), common.NewField("host", host))
	return nil
}

func (r *Redis) Do(rdb *redis.Client, commands []string) error {
	if len(commands) == 0 {
		logger.Debug(common.ErrInvalidArg.Error(), common.NewField("commands", commands))
		return nil
	}
	var args []string
	if len(commands) == 1 {
		args = strings.Fields(commands[0])
	} else {
		args = commands
	}
	arg := common.SliceStringToInterface(args)

	cmd := rdb.Do(common.Context, arg...)
	out, err := cmd.Result()
	if err != nil {
		logger.Debug(err.Error())
		return err
	}

	switch data := out.(type) {
	case []any:
		for i := 0; i < len(data); i++ {
			printer.Printf("%d) %s\n", i+1, data[i])
		}
	default:
		printer.Printf(printer.SetJSONAsDefaultFormat(rootOutputFormat), out)
	}
	return err
}
