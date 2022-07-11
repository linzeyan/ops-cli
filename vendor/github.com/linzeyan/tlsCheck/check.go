package tlsCheck

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"os"
)

func CheckByHost(host string) (*tls.Conn, error) {
	conn, err := tls.Dial("tcp", host, nil)
	if err != nil {
		return nil, err
	}
	if conn != nil {
		defer conn.Close()
	}
	return conn, err
}

func CheckByFile(fileName string) ([]*x509.Certificate, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	buf := make([]byte, 4096*3)
	var t int
	for {
		n, err := reader.Read(buf)
		if n == 0 {
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println(err)
				break
			}
		}
		t = n
	}
	buf = buf[0:t]
	crtPem, _ := pem.Decode(buf)
	crt, err := x509.ParseCertificates(crtPem.Bytes)
	if err != nil {
		return nil, err
	}
	return crt, nil
}
