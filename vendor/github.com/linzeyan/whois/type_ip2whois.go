package whois

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Ip2Whois struct {
	Domain      string    `json:"domain"`
	DomainID    string    `json:"domain_id"`
	Status      string    `json:"status"`
	CreateDate  string    `json:"create_date"`
	UpdateDate  string    `json:"update_date"`
	ExpireDate  string    `json:"expire_date"`
	DomainAge   int64     `json:"domain_age"`
	WhoisServer string    `json:"whois_server"`
	Registrar   Registrar `json:"registrar"`
	Registrant  Admin     `json:"registrant"`
	Admin       Admin     `json:"admin"`
	Tech        Admin     `json:"tech"`
	Billing     Admin     `json:"billing"`
	Nameservers []string  `json:"nameservers"`
}

func (w Ip2Whois) Request(domain string) (*Response, error) {
	if IP2WhoisKey != "" {
		Key = IP2WhoisKey
	}
	apiUrl := fmt.Sprintf("https://api.ip2whois.com/v2?key=%s&domain=%s", Key, domain)
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
		var w Ip2Whois
		err = json.Unmarshal(content, &w)
		if err != nil {
			return nil, err
		}
		var r = Response{
			CreatedDate: w.CreateDate,
			UpdatedDate: w.UpdateDate,
			ExpiresDate: w.ExpireDate,
			NameServers: w.Nameservers,
			Registrar:   w.Registrar.Name,
		}
		return &r, nil
	}
	return nil, err
}

type Admin struct {
	Name          string `json:"name"`
	Organization  string `json:"organization"`
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	Region        string `json:"region"`
	ZipCode       string `json:"zip_code"`
	Country       string `json:"country"`
	Phone         string `json:"phone"`
	Fax           string `json:"fax"`
	Email         string `json:"email"`
}

type Registrar struct {
	IANAID string `json:"iana_id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
}
