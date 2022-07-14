package whois

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
