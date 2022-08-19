package common

import (
	"context"
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/fatih/color"
)

var (
	Context = context.Background()
	TimeNow = time.Now().Local()
)

var (
	ErrConfigContent = errors.New("config content is incorrect")
	ErrConfigTable   = errors.New("table not found in the config")
	ErrInvalidArg    = errors.New("invalid argument")
	ErrInvalidURL    = errors.New("invalid URL")
	ErrPicSize       = errors.New("picture size is too small")
	ErrStatusCode    = errors.New("status code is not 200")
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
