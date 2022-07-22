package whois

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WhoApi struct {
	Status            string              `json:"status"`
	WhoisServer       string              `json:"whois_server"`
	StatusDesc        string              `json:"status_desc"`
	LimitHit          bool                `json:"limit_hit"`
	Registered        bool                `json:"registered"`
	WhoisRaw          string              `json:"whois_raw"`
	Disclaimer        string              `json:"disclaimer"`
	Premium           bool                `json:"premium"`
	GenericWhois      bool                `json:"generic_whois"`
	RegistryDomainID  string              `json:"registry_domain_id"`
	RegistrarIANAID   string              `json:"registrar_iana_id"`
	DateCreated       string              `json:"date_created"`
	DateExpires       string              `json:"date_expires"`
	DateUpdated       string              `json:"date_updated"`
	DomainStatus      []string            `json:"domain_status"`
	Nameservers       []string            `json:"nameservers"`
	Emails            string              `json:"emails"`
	WhoisRawParent    string              `json:"whois_raw_parent"`
	WhoisName         string              `json:"whois_name"`
	Contacts          []map[string]string `json:"contacts"`
	DomainName        string              `json:"domain_name"`
	Cached            bool                `json:"_cached"`
	RequestsAvailable int64               `json:"requests_available"`
}

func (w WhoApi) Request(domain string) (*Response, error) {
	apiUrl := fmt.Sprintf("http://api.whoapi.com/?r=whois&apikey=%s&domain=%s", WhoApiKey, domain)
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
		var w WhoApi
		err = json.Unmarshal(content, &w)
		if err != nil {
			return nil, err
		}
		var r = Response{
			CreatedDate: w.DateCreated,
			UpdatedDate: w.DateUpdated,
			ExpiresDate: w.DateExpires,
			NameServers: w.Nameservers,
			Registrar:   w.Contacts[0]["organization"],
		}
		return &r, nil
	}
	return nil, err
}
