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

// LookupIPReverse выполняет обратный DNS поиск и возвращает имена, связанные с IP
func LookupIPReverse(ip string) []string {
	hosts, err := net.LookupAddr(ip)
	if err != nil {
		return nil
	}
	return hosts
}
