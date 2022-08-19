package common

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var (
	Context = context.Background()
	TimeNow = time.Now().Local()
)

type ByteSize float64

const (
	_           = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

/* Read config.toml. */
type ConfigBlock string

const (
	Discord  ConfigBlock = "discord"
	ICP      ConfigBlock = "west"
	LINE     ConfigBlock = "line"
	Slack    ConfigBlock = "slack"
	Telegram ConfigBlock = "telegram"
)

func (c ConfigBlock) String() string {
	return string(c)
}

type readConfig struct {
	table ConfigBlock
}

func (c readConfig) get(config string) (map[string]interface{}, error) {
	viper.SetConfigFile(config)
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	all := viper.AllSettings()
	values, ok := all[c.table.String()]
	if !ok {
		return nil, errors.New("table not found in config")
	}
	v, ok := values.(map[string]interface{})
	if !ok {
		return nil, errors.New("config content is incorrect")
	}
	return v, nil
}

/* Get secret token from config. */
func Config(config string, t ConfigBlock) (map[string]interface{}, error) {
	r := readConfig{table: t}
	return r.get(config)
}

/* Print examples with color. */
func Examples(s string) string {
	c := color.New(color.FgYellow)
	return c.Sprintf(`%s`, s)
}

const (
	HashMd5        = "md5"
	HashSha1       = "sha1"
	HashSha224     = "sha224"
	HashSha256     = "sha256"
	HashSha384     = "sha384"
	HashSha512     = "sha512"
	HashSha512_224 = "sha512_224"
	HashSha512_256 = "sha512_256"
)

func HashAlgorithm(alg string) hash.Hash {
	m := map[string]hash.Hash{
		HashMd5:        md5.New(),
		HashSha1:       sha1.New(),
		HashSha224:     sha256.New224(),
		HashSha256:     sha256.New(),
		HashSha384:     sha512.New384(),
		HashSha512:     sha512.New(),
		HashSha512_224: sha512.New512_224(),
		HashSha512_256: sha512.New512_256(),
	}
	return m[alg]
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
