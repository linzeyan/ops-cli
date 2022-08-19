package common

import "github.com/fatih/color"

/* Print examples with color. */
func Examples(s string) string {
	c := color.New(color.FgYellow)
	return c.Sprintf(`%s`, s)
}
