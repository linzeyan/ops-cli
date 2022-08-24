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
	"io"
	"net/http"

	"golang.org/x/text/encoding/simplifiedchinese"
)

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

/* Same as HttpRequestContent, but read content in GB18030 format. */
func HTTPRequestContentGB18030(url string, body io.Reader, methods ...string) ([]byte, error) {
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
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
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
	reader := simplifiedchinese.GB18030.NewDecoder().Reader(resp.Body)
	return io.ReadAll(reader)
}

/* Return redirect url. */
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
