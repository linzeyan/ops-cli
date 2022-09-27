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
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func init() {
	var wspingCmd = &cobra.Command{
		Use:   "wsping",
		Args:  cobra.ExactArgs(1),
		Short: "Connect to a websocket server",
		RunE: func(_ *cobra.Command, args []string) error {
			if !validator.ValidURL(args[0]) {
				return common.ErrInvalidURL
			}
			d := websocket.Dialer{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
			}
			ws, resp, err := d.Dial(args[0], nil)
			if err != nil {
				return err
			}
			defer ws.Close()
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusSwitchingProtocols {
				o := fmt.Sprintf("Status Code is not %d", http.StatusSwitchingProtocols)
				return errors.New(o)
			}
			// err = ws.WriteMessage(websocket.PingMessage, []byte{})
			// if err != nil {
			// 	return err
			// }
			// _, message, err := ws.ReadMessage()
			// if err != nil {
			// 	return err
			// }
			// if string(message) == "" {
			// 	return common.ErrResponse
			// }
			PrintString("Connect success")
			return err
		},
	}
	rootCmd.AddCommand(wspingCmd)
}
