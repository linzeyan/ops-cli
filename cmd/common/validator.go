/*
Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"net/netip"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
)

/* If i is a domain return true. */
func IsDomain(i any) bool {
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
func IsFile(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

/* If i is a ipv address return true. */
func IsIP(i string) bool {
	ip, err := netip.ParseAddr(i)
	if err != nil {
		return false
	}
	return ip.IsValid()
}

func IsCIDR(i string) bool {
	ip, err := netip.ParsePrefix(i)
	if err != nil {
		return false
	}
	return ip.IsValid()
}

/* If i is a ipv4 address return true. */
func IsIPv4(i string) bool {
	ip, err := netip.ParseAddr(i)
	if err != nil {
		return false
	}
	return ip.Is4()
}

func IsIPv4CIDR(i string) bool {
	ip, err := netip.ParsePrefix(i)
	if err != nil {
		return false
	}
	return ip.IsValid() && ip.Addr().Is4()
}

/* If i is a ipv6 address return true. */
func IsIPv6(i string) bool {
	ip, err := netip.ParseAddr(i)
	if err != nil {
		return false
	}
	return ip.Is6()
}

func IsIPv6CIDR(i string) bool {
	ip, err := netip.ParsePrefix(i)
	if err != nil {
		return false
	}
	return ip.IsValid() && ip.Addr().Is6()
}

/* If u is a valid url return true. */
func IsURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}

func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}