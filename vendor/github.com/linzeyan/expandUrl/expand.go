package expandUrl

import (
	"net/http"
	"net/url"
)

func Expand(uri string) (string, error) {
	const ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	}

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

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	if req.Host == "reurl.cc" && resp.Header.Get("Target") != "" {
		result = resp.Header.Get("Target")
	}

	return result, nil
}
