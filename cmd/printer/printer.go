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

package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

type Printer struct {
	/* Table args. */
	headers bool
	padding string
	align   int
}

func (p *Printer) Printf(format string, a ...any) {
	switch format {
	case "":
		for _, i := range a {
			switch data := i.(type) {
			case string, []byte:
				fmt.Fprintf(os.Stdout, "%s", data)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				fmt.Fprintf(os.Stdout, "%d", data)
			case []string, map[string]string:
				p.json(data)
			case error:
				fmt.Fprintln(os.Stderr, data)
			default:
				fmt.Fprintf(os.Stdout, "%v", data)
			}
		}
	case "json":
		p.json(a...)
	case "table":
		if len(a) != 2 {
			fmt.Fprintln(os.Stderr, common.ErrInvalidArg)
			return
		}
		/* assume a[0] is header. */
		h1, ok1 := a[0].([]string)
		d1, ok2 := a[1].([][]string)
		if ok1 && ok2 {
			p.table(h1, d1)
			return
		}

		/* assume a[1] is header. */
		h2, ok3 := a[1].([]string)
		d2, ok4 := a[0].([][]string)
		if ok3 && ok4 {
			p.table(h2, d2)
			return
		}
		fmt.Fprintln(os.Stderr, common.ErrInvalidArg)
	case "yaml":
		p.yaml(a...)
	default:
		fmt.Fprintf(os.Stdout, format, a...)
	}
}

func (*Printer) json(a ...any) {
	for _, i := range a {
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(i); err != nil {
			fmt.Fprintln(os.Stderr, common.ErrInvalidArg)
			return
		}
		fmt.Fprintln(os.Stdout, buf.String())
	}
}

func (*Printer) yaml(a ...any) {
	for _, i := range a {
		var buf bytes.Buffer
		encoder := yaml.NewEncoder(&buf)
		encoder.SetIndent(2)
		if err := encoder.Encode(i); err != nil {
			fmt.Fprintln(os.Stderr, common.ErrInvalidArg)
			return
		}
		fmt.Fprintln(os.Stdout, buf.String())
	}
}

func (p *Printer) table(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(data)

	table.SetAutoFormatHeaders(p.headers)
	table.SetHeaderAlignment(p.align)
	table.SetAlignment(p.align)
	table.SetTablePadding(p.padding)

	table.SetAutoWrapText(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetNoWhiteSpace(true)

	table.Render()
}

func (p *Printer) SetTableAlign(align int)           { p.align = align }
func (p *Printer) SetTablePadding(padding string)    { p.padding = padding }
func (p *Printer) SetTableFormatHeaders(format bool) { p.headers = format }

func NewPrinter() *Printer {
	return &Printer{}
}
