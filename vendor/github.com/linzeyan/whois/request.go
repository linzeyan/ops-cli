package whois

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ua string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"

// go:embed key_whoisxmlapi
var WhoisXMLAPIKey string

func RequestWhoisXML(domain string) (*WhoisRecord, error) {
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
		var data WhoisXML
		err = json.Unmarshal(content, &data)
		if err != nil {
			return nil, err
		}
		return &data.WhoisRecord, nil
	}
	return nil, err
}

func ParserWhoisXML(data *WhoisRecord) map[string]interface{} {
	var result = make(map[string]interface{})
	result["Audit"] = map[string]string{
		"CreatedDate": data.CreatedDate,
		"ExpiresDate": data.ExpiresDate,
		"UpdatedDate": data.UpdatedDate,
	}
	result["NameServers"] = data.NameServers.HostNames
	result["Registrant"] = map[string]string{
		"Country":      data.Registrant.CountryCode,
		"Organization": data.Registrant.Organization,
		"State":        data.Registrant.State,
	}
	result["Registrar"] = data.RegistrarName
	return result
}

// go:embed key_ip2whois
var IP2WhoisKey string

func RequestIp2Whois(domain string) (*Ip2Whois, error) {
	apiUrl := fmt.Sprintf("https://api.ip2whois.com/v2?key=%s&domain=%s", IP2WhoisKey, domain)
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
		var data Ip2Whois
		err = json.Unmarshal(content, &data)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
	return nil, err
}

func ParserIp2Whois(data *Ip2Whois) map[string]interface{} {
	var result = make(map[string]interface{})
	result["Audit"] = map[string]string{
		"CreatedDate": data.CreateDate,
		"ExpiresDate": data.ExpireDate,
		"UpdatedDate": data.UpdateDate,
	}
	result["NameServers"] = data.Nameservers
	result["Registrant"] = map[string]string{
		"Country":      data.Registrant.Country,
		"Organization": data.Registrant.Organization,
		"State":        data.Registrant.Region,
	}
	result["Registrar"] = data.Registrar.Name
	return result
}

// go:embed key_whoapi
var WhoApiKey string

func RequestWhoApi(domain string) (*WhoApi, error) {
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
		var data WhoApi
		err = json.Unmarshal(content, &data)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
	return nil, err
}

func ParserWhoApi(data *WhoApi) map[string]interface{} {
	var result = make(map[string]interface{})
	result["Audit"] = map[string]string{
		"CreatedDate": data.DateCreated,
		"ExpiresDate": data.DateExpires,
		"UpdatedDate": data.DateUpdated,
	}
	result["NameServers"] = data.Nameservers
	result["Registrant"] = map[string]string{
		"Country":      data.Contacts[1]["country"],
		"Organization": data.Contacts[1]["organization"],
		"State":        data.Contacts[1]["state"],
	}
	result["Registrar"] = data.Contacts[0]["organization"]
	return result
}

// go:embed key_apininjas
var ApiNinjasKey string

func RequestApiNinjas(domain string) (*ApiNinjas, error) {
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
	req.Header.Set("X-Api-Key", ApiNinjasKey)

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
		var data ApiNinjas
		err = json.Unmarshal(content, &data)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
	return nil, err
}

func ParserApiNinjas(data *ApiNinjas) map[string]interface{} {
	var result = make(map[string]interface{})
	// result["Audit"] = map[string]string{
	// 	"CreatedDate": time.Unix(data.CreationDate, 0).Format(time.RFC3339),
	// 	"ExpiresDate": time.Unix(data.ExpirationDate, 0).Format(time.RFC3339),
	// 	"UpdatedDate": time.Unix(data.UpdatedDate, 0).Format(time.RFC3339),
	// }

	result["NameServers"] = data.NameServers
	result["Registrant"] = map[string]string{
		"Country":      data.Country,
		"Organization": data.Org,
		"State":        data.State,
	}
	result["Registrar"] = data.Registrar
	return result
}
