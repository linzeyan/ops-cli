package whois

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ApiNinjas struct {
	DomainName     []string    `json:"domain_name"`
	Registrar      string      `json:"registrar"`
	WhoisServer    string      `json:"whois_server"`
	UpdatedDate    interface{} `json:"updated_date"`
	CreationDate   interface{} `json:"creation_date"`
	ExpirationDate interface{} `json:"expiration_date"`
	NameServers    []string    `json:"name_servers"`
	Emails         []string    `json:"emails"`
	Dnssec         interface{} `json:"dnssec"`
	Name           string      `json:"name"`
	Org            string      `json:"org"`
	Address        string      `json:"address"`
	City           string      `json:"city"`
	State          string      `json:"state"`
	Zipcode        string      `json:"zipcode"`
	Country        string      `json:"country"`
}

func (w ApiNinjas) Request(domain string) (*Response, error) {
	apiUrl := fmt.Sprintf("https://api.api-ninjas.com/v1/whois?domain=%s", domain)
	var client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ua)
	if ApiNinjasKey != "" {
		Key = ApiNinjasKey
	}
	req.Header.Set("X-Api-Key", Key)

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var w ApiNinjas
		err = json.Unmarshal(content, &w)
		if err != nil {
			return nil, err
		}
		// 	"CreatedDate": time.Unix(data.CreationDate, 0).Format(time.RFC3339),
		// 	"ExpiresDate": time.Unix(data.ExpirationDate, 0).Format(time.RFC3339),
		// 	"UpdatedDate": time.Unix(data.UpdatedDate, 0).Format(time.RFC3339),
		var r = Response{
			NameServers: w.NameServers,
			Registrar:   w.Registrar,
		}
		return &r, nil
	}
	return nil, err
}
