package whois

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WhoisXML struct {
	WhoisRecord WhoisRecord `json:"WhoisRecord"`
}

func (w WhoisXML) Request(domain string) (*Response, error) {
	apiUrl := fmt.Sprintf("https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=%s&domainName=%s&outputFormat=JSON", WhoisXMLAPIKey, domain)

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
		err = json.Unmarshal(content, &w)
		if err != nil {
			return nil, err
		}
		var r = Response{
			CreatedDate: w.WhoisRecord.CreatedDate,
			UpdatedDate: w.WhoisRecord.UpdatedDate,
			ExpiresDate: w.WhoisRecord.ExpiresDate,
			NameServers: w.WhoisRecord.NameServers.HostNames,
			Registrar:   w.WhoisRecord.RegistrarName,
		}
		return &r, nil
	}
	return nil, err
}

type WhoisRecord struct {
	CreatedDate           string                `json:"createdDate"`
	UpdatedDate           string                `json:"updatedDate"`
	ExpiresDate           string                `json:"expiresDate"`
	Registrant            AdministrativeContact `json:"registrant"`
	AdministrativeContact AdministrativeContact `json:"administrativeContact"`
	TechnicalContact      AdministrativeContact `json:"technicalContact"`
	DomainName            string                `json:"domainName"`
	NameServers           NameServers           `json:"nameServers"`
	Status                string                `json:"status"`
	RawText               string                `json:"rawText"`
	ParseCode             int64                 `json:"parseCode"`
	Header                string                `json:"header"`
	StrippedText          string                `json:"strippedText"`
	Footer                string                `json:"footer"`
	Audit                 Audit                 `json:"audit"`
	RegistrarName         string                `json:"registrarName"`
	RegistrarIANAID       string                `json:"registrarIANAID"`
	CreatedDateNormalized string                `json:"createdDateNormalized"`
	UpdatedDateNormalized string                `json:"updatedDateNormalized"`
	ExpiresDateNormalized string                `json:"expiresDateNormalized"`
	RegistryData          RegistryData          `json:"registryData"`
	ContactEmail          string                `json:"contactEmail"`
	DomainNameEXT         string                `json:"domainNameExt"`
	EstimatedDomainAge    int64                 `json:"estimatedDomainAge"`
}

type AdministrativeContact struct {
	Organization string `json:"organization"`
	State        string `json:"state"`
	Country      string `json:"country"`
	CountryCode  string `json:"countryCode"`
	RawText      string `json:"rawText"`
}

type Audit struct {
	CreatedDate string `json:"createdDate"`
	UpdatedDate string `json:"updatedDate"`
}

type NameServers struct {
	RawText   string        `json:"rawText"`
	HostNames []string      `json:"hostNames"`
	IPS       []interface{} `json:"ips"`
}

type RegistryData struct {
	CreatedDate           string      `json:"createdDate"`
	UpdatedDate           string      `json:"updatedDate"`
	ExpiresDate           string      `json:"expiresDate"`
	DomainName            string      `json:"domainName"`
	NameServers           NameServers `json:"nameServers"`
	Status                string      `json:"status"`
	RawText               string      `json:"rawText"`
	ParseCode             int64       `json:"parseCode"`
	Header                string      `json:"header"`
	StrippedText          string      `json:"strippedText"`
	Footer                string      `json:"footer"`
	Audit                 Audit       `json:"audit"`
	RegistrarName         string      `json:"registrarName"`
	RegistrarIANAID       string      `json:"registrarIANAID"`
	CreatedDateNormalized string      `json:"createdDateNormalized"`
	UpdatedDateNormalized string      `json:"updatedDateNormalized"`
	ExpiresDateNormalized string      `json:"expiresDateNormalized"`
	WhoisServer           string      `json:"whoisServer"`
}
