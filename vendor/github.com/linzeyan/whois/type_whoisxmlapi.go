package whois

type WhoisXML struct {
	WhoisRecord WhoisRecord `json:"WhoisRecord"`
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
