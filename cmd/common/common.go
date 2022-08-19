package common

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/fatih/color"
)

var (
	Context = context.Background()
	TimeNow = time.Now().Local()
)

/* Print string with color. */
func Examples(s string) string {
	c := color.New(color.FgYellow)
	return c.Sprintf(`%s`, s)
}

/* HttpRequestContent make a simple request to url, and return response body, default request method is get. */
func HTTPRequestContent(url string, body io.Reader, methods ...string) ([]byte, error) {
	var method string
	if len(methods) == 0 {
		method = http.MethodGet
	} else {
		method = methods[0]
	}
	var client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	req, err := http.NewRequestWithContext(Context, method, url, body)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return content, err
	}
	return nil, errors.New("status code is not 200")
}
