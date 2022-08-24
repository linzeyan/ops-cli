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
	"os"
	"regexp"

	"github.com/fatih/color"
)

func Dos2Unix(filename string) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	stat, err := os.Stat(filename)
	if err != nil {
		return err
	}
	eol := regexp.MustCompile(`\r\n`)
	f = eol.ReplaceAllLiteral(f, []byte{'\n'})
	return os.WriteFile(filename, f, stat.Mode())
}

/* Print string with color. */
func Examples(s string) string {
	c := color.New(color.FgYellow)
	return c.Sprintf(`%s`, s)
}

func Usage(s string) string {
	c := color.New(color.FgGreen)
	return c.Sprintf(`%s`, s)
}
