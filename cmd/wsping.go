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
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func init() {
	var wspingCmd = &cobra.Command{
		Use:   CommandWsping + " host",
		Args:  cobra.ExactArgs(1),
		Short: "Connect to a websocket server",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"ws://", "wss://"}, cobra.ShellCompDirectiveNoSpace
		},
		Run: func(_ *cobra.Command, args []string) {
			if !validator.ValidURL(args[0]) {
				PrintString(common.ErrInvalidURL)
				return
			}
			d := websocket.Dialer{
				HandshakeTimeout: 2 * time.Second,
				ReadBufferSize:   1024,
				WriteBufferSize:  1024,
			}

			ws, resp, err := d.DialContext(common.Context, args[0], nil)
			if err != nil {
				PrintString(err)
				return
			}
			defer ws.Close()
			if resp != nil {
				defer resp.Body.Close()
			}
			if resp.StatusCode != http.StatusSwitchingProtocols {
				PrintString("Status is not " + http.StatusText(http.StatusSwitchingProtocols))
				return
			}
			// err = ws.WriteMessage(websocket.PingMessage, []byte{})
			// if err != nil {
			// 	return err
			// }
			// _, message, err := ws.ReadMessage()
			// if err != nil {
			// 	return err
			// }
			// message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
			// if string(message) == "" {
			// 	return common.ErrResponse
			// }
			PrintString("Connect success")
		},
	}
	RootCmd.AddCommand(wspingCmd)
}
