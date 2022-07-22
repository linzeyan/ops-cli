package whois

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

const ua string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"

// go:embed key_whoisxmlapi
var WhoisXMLAPIKey string

// go:embed key_ip2whois
var IP2WhoisKey string

// go:embed key_whoapi
var WhoApiKey string

// go:embed key_apininjas
var ApiNinjasKey string

var Key string

type Server interface {
	Request(string) (*Response, error)
}

func Request(s Server, domain string) *Response {
	resp, err := s.Request(domain)
	if err != nil {
		log.Println(err)
		return nil
	}
	return resp
}

type Response struct {
	Registrar   string   `json:"Registrar" yaml:"Registrar"`
	CreatedDate string   `json:"CreatedDate" yaml:"CreatedDate"`
	ExpiresDate string   `json:"ExpiresDate" yaml:"ExpiresDate"`
	UpdatedDate string   `json:"UpdatedDate" yaml:"UpdatedDate"`
	NameServers []string `json:"NameServers" yaml:"NameServers"`
}

func (r Response) String() {
	f := reflect.ValueOf(&r).Elem()
	t := f.Type()
	for i := 0; i < f.NumField(); i++ {
		fmt.Printf("%s\t%v\n", t.Field(i).Name, f.Field(i).Interface())
		//f.Field(i).Type()
	}
}

func (r Response) Json() {
	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (r Response) Yaml() {
	out, err := yaml.Marshal(r)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

type Verisign struct{}

func (w Verisign) Request(domain string) (*Response, error) {
	conn, err := net.Dial("tcp", "whois.verisign-grs.com:43")
	if err != nil {
		return nil, err
	}
	if conn != nil {
		defer conn.Close()
	}
	_, err = conn.Write([]byte(domain + "\n"))
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	replace := strings.ReplaceAll(string(result), ": ", ";")
	replace1 := strings.ReplaceAll(replace, "\r\n", ",")
	split := strings.Split(replace1, ",")
	var ns []string
	var r Response
	for i := range split {
		if strings.Contains(split[i], "Updated Date") {
			v := strings.Split(split[i], ";")
			r.UpdatedDate = v[1]
		}
		if strings.Contains(split[i], "Creation Date") {
			v := strings.Split(split[i], ";")
			r.CreatedDate = v[1]
		}
		if strings.Contains(split[i], "Registry Expiry Date") {
			v := strings.Split(split[i], ";")
			r.ExpiresDate = v[1]
		}
		if strings.Contains(split[i], "Registrar") {
			v := strings.Split(split[i], ";")
			if strings.TrimSpace(v[0]) == "Registrar" {
				r.Registrar = v[1]
			}
		}
		if strings.Contains(split[i], "Name Server") {
			v := strings.Split(split[i], ";")
			ns = append(ns, v[1])
		}
	}
	r.NameServers = ns
	return &r, nil
}

func RequestIana(domain string) (string, error) {
	conn, err := net.Dial("tcp", "whois.iana.org:43")
	if err != nil {
		return "", err
	}
	if conn != nil {
		defer conn.Close()
	}
	_, err = conn.Write([]byte(domain + "\n"))
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return "", err
	}
	return string(result), nil

}

func ParseIana(data string) map[string]string {
	var result = make(map[string]string)
	replace := strings.ReplaceAll(data, ": ", ";")
	replace1 := strings.ReplaceAll(replace, "\n", ",")
	split := strings.Split(replace1, ",")
	for i := range split {
		if strings.Contains(split[i], "organisation") {
			v := strings.Split(split[i], ";")
			result["Organization"] = v[1]
			break
		}
	}
	return result
}
