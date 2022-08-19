package common

import (
	"io"
	"net/http"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"

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
	if resp.StatusCode != http.StatusOK {
		return nil, ErrStatusCode
	}
	return io.ReadAll(resp.Body)
}

func HTTPRequestRedirectURL(uri string) (string, error) {
	var result = uri
	var client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		CheckRedirect: func(req *http.Request, _ []*http.Request) error {
			result = req.URL.String()
			return nil
		},
	}
	req, err := http.NewRequestWithContext(Context, http.MethodGet, uri, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return result, err
	}
	if req.Host == "reurl.cc" && resp.Header.Get("Target") != "" {
		result = resp.Header.Get("Target")
	}
	return result, err
}
