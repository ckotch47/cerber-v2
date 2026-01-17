package command

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"cerber/internal/dns"
	"cerber/internal/style"
)

var LookCmd = &cobra.Command{
	Use:   "look",
	Short: "Найти IP по домену или доменные имена по IP",
	Long:  "Примеры:\n  cerber look http://example.com — найти IP по домену\n  cerber look 8.8.8.8 — найти доменные имена по IP",
	Run:   lookupHostRun,
}

func lookupHostRun(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Введите домен или IP-адрес")
		return
	}

	host := cleanDomain(args[0])

	// Проверяем, является ли аргумент IP-адресом
	if net.ParseIP(host) != nil {
		// Если это IP, выполняем обратный DNS-поиск
		lookUpIP(host)
		return
	}

	// Иначе считаем, что это домен, ищем IP
	res := dns.CheckDomain(host)
	if len(res) == 0 {
		fmt.Println(style.NotFoundStyle.Render("Not found"))
		return
	}

	for _, ip := range res {
		fmt.Println(style.SuccessStyle.Render(ip))
	}

}

func lookUpIP(ip string) {
	res := dns.LookupIPReverse(ip)
	if len(res) == 0 {
		fmt.Println(style.NotFoundStyle.Render("Not found"))
	}
	for _, domain := range res {
		domain = strings.TrimSuffix(domain, ".") // Убираем точку в конце, если есть
		fmt.Println(style.SuccessStyle.Render(domain))
	}
}

func lookupHost(host string) string {
	if res := dns.CheckDomain(host); len(res) > 0 {
		return res[0]
	}
	return ""
}
