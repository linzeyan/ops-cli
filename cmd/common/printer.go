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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

/* Print some formats easily. */
type printer struct {
	/* JSON and yaml args. */
	indent int
	/* Table args. */
	headers bool
	padding string
	align   int
}

func (p *printer) Printf(format string, a ...any) {
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
				p.Error(data)
			default:
				fmt.Fprintf(os.Stdout, "%v", data)
			}
		}
	case JSONFormat:
		p.json(a...)
	case NoneFormat:
	case TableFormat:
		if len(a) != 2 {
			p.Error(ErrInvalidArg)
			return
		}
		/* Assume a[0] is header. */
		h1, ok1 := a[0].([]string)
		d1, ok2 := a[1].([][]string)
		if ok1 && ok2 {
			p.table(h1, d1)
			return
		}

		/* Assume a[1] is header. */
		h2, ok3 := a[1].([]string)
		d2, ok4 := a[0].([][]string)
		if ok3 && ok4 {
			p.table(h2, d2)
			return
		}
		p.Error(ErrInvalidArg)
	case YamlFormat:
		p.yaml(a...)
	default:
		fmt.Fprintf(os.Stdout, format, a...)
	}
}

/* Print error to stderr. */
func (*printer) Error(err error) {
	fmt.Fprintln(os.Stderr, err)
}

func (p *printer) json(a ...any) {
	for _, i := range a {
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		indent := strings.Repeat(" ", p.indent)
		encoder.SetIndent("", indent)
		if err := encoder.Encode(i); err != nil {
			p.Error(ErrInvalidArg)
			return
		}
		fmt.Fprintf(os.Stdout, "%s", buf.String())
	}
}

func (p *printer) yaml(a ...any) {
	for _, i := range a {
		var buf bytes.Buffer
		encoder := yaml.NewEncoder(&buf)
		encoder.SetIndent(p.indent)
		if err := encoder.Encode(i); err != nil {
			p.Error(ErrInvalidArg)
			return
		}
		fmt.Fprintf(os.Stdout, "%s", buf.String())
	}
}

func (p *printer) table(header []string, data [][]string) {
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

func (p *printer) SetIdent(indent int)               { p.indent = indent }
func (p *printer) SetTableAlign(align int)           { p.align = align }
func (p *printer) SetTablePadding(padding string)    { p.padding = padding }
func (p *printer) SetTableFormatHeaders(format bool) { p.headers = format }

func (*printer) SetJSONAsDefaultFormat(format string) string {
	if format == "" {
		return JSONFormat
	}
	return format
}

func (*printer) SetNoneAsDefaultFormat(format string) string {
	if format == "" {
		return NoneFormat
	}
	return format
}

func (*printer) SetTableAsDefaultFormat(format string) string {
	if format == "" {
		return TableFormat
	}
	return format
}

func (*printer) SetYamlAsDefaultFormat(format string) string {
	if format == "" {
		return YamlFormat
	}
	return format
}

func NewPrinter() *printer {
	return &printer{
		indent: 2,
	}
}
