package whois

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
