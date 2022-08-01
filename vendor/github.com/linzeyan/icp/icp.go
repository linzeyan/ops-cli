package icp

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var (
	ConfigFile, Domain      string
	WestAccount, WestApiKey string
)

type West struct {
	DomainName string `json:"domain"`
	ICPCode    string `json:"icp"`
	ICPStatus  string `json:"icpstatus"`
}

func md5encode(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func requestURI() (uri string) {
	/* MD5 Hash */
	var hash_data string = WestAccount + WestApiKey + "domainname"
	sig := md5encode(hash_data)
	rawCmd := fmt.Sprintf("domainname\r\ncheck\r\nentityname:icp\r\ndomains:%s\r\n.\r\n", Domain)
	/* URL Encoding */
	strCmd := url.QueryEscape(rawCmd)
	return fmt.Sprintf(`http://api.west263.com/api/?userid=%s&strCmd=%s&versig=%s`, WestAccount, strCmd, sig)
}

func httpPOST() (content []byte, err error) {
	var tr = &http.Transport{DisableKeepAlives: true}
	var client = &http.Client{Transport: tr}
	uri := requestURI()
	data := strings.NewReader(``)
	req, err := http.NewRequest(http.MethodPost, uri, data)
	if err != nil {
		fmt.Println("Resquest error.")
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		fmt.Println("Response error.")
		fmt.Println(err)
		return nil, err
	}

	/* Convert GBK to UTF-8 */
	reader := simplifiedchinese.GB18030.NewDecoder().Reader(resp.Body)
	content, err = ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("Content error.")
		fmt.Println(err)
		return nil, err
	}
	return
}

func Check() string {
	body, err := httpPOST()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	/* Find String */
	re, _ := regexp.Compile("{.*}")
	match := fmt.Sprintln(re.FindString(string(body)))
	/* Parse Json */
	var icp West
	json.Unmarshal([]byte(match), &icp)
	return icp.ICPStatus
}

func ReadConf() {
	if ConfigFile != "" {
		viper.SetConfigType("toml")
		viper.SetConfigFile(ConfigFile)
	} else {
		return
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	WestAccount = viper.GetString("west_api.account")
	WestApiKey = viper.GetString("west_api.key")
}
