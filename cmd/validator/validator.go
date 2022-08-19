package validator

import (
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

/* If i is a domain return true. */
func ValidDomain(i any) bool {
	const elements = "~!@#$%^&*()_+`={}|[]\\:\"<>?,/"
	if val, ok := i.(string); ok {
		if strings.ContainsAny(val, elements) {
			return false
		}
		slice := strings.Split(val, ".")
		l := len(slice)
		if l > 1 {
			n, err := strconv.Atoi(slice[l-1])
			if err != nil {
				return true
			}
			s := strconv.Itoa(n)
			return slice[l-1] != s
		}
	}
	return false
}

/* If f is a valid path return true. */
func ValidFile(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

/* If i is a ipv address return true. */
func ValidIP(i string) bool {
	return net.ParseIP(i) != nil
}

/* If i is a ipv4 address return true. */
func ValidIPv4(i string) bool {
	return net.ParseIP(i).To4() != nil
}

/* If i is a ipv6 address return true. */
func ValidIPv6(i string) bool {
	return net.ParseIP(i).To4() == nil && net.ParseIP(i).To16() != nil
}

/* If u is a valid url return true. */
func ValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
