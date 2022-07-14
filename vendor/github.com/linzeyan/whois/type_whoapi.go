package whois

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
