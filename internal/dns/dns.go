package dns

import (
	
	"net"
)

func CheckDomain(domain string) []string {
	if hosts, err := net.LookupHost(domain); err != nil {
		return nil
	} else {
		return hosts
	}
}

func CheckIP(ip string) bool {
	if _, err := net.LookupIP(ip); err != nil {
		return false
	} else {
		return true
	}
}


