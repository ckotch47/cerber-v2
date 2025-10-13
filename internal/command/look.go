package command

import (
	"fmt"


	"github.com/spf13/cobra"

	"cerber/internal/dns"
	"cerber/internal/style"

)

var LookCmd = &cobra.Command{
	Use:   "look",
	Short: "Найти ip адрес по домену \nexample: cerber look http://example.com",
	Run: lookupHostRun,
}


func lookupHostRun(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Введите домен")
		return
	}
	host := cleanDomain(args[0])
	res := dns.CheckDomain(host)
	if len(res) == 0 {
		fmt.Println(style.NotFoundStyle.Render("Not found"))
		return 
	}
	for _, ip := range res {
		fmt.Println(style.SuccessStyle.Render(ip))
	}
}

func lookupHost(host string) string {
	if res := dns.CheckDomain(host); len(res) > 0 {
		return res[0]
	} else {
		return ""
	}
}


