package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"cerber/internal/style"
	"cerber/internal/utils"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Выполняет поиск поддоменов по списку из файла",
	Run:   FindHost,
}

var commandBruteForce utils.BruteForceType
var MaxDepth int = 2

func init() {
	findCmd.Flags().StringVarP(
		&commandBruteForce.WorldList,
		"worldlis",
		"w",
		"",
		"Файл со списком",
	)
	findCmd.Flags().BoolVarP(
		&commandBruteForce.Recurse,
		"recurse",
		"r",
		false,
		"Включить рекурсию для брутфорса",
	)
}

func FindHost(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println(style.NotFoundStyle.Render("Не указан домен"))
		return
	}
	if commandBruteForce.WorldList == "" {
		fmt.Println(style.NotFoundStyle.Render("Файл со списком не найден"))
		return
	}

	domain := cleanDomain(args[0])
	domainList := utils.ReadFile(commandBruteForce.WorldList)

	findSubDomains(domain, domainList, 0)
}

func findSubDomains(domain string, worldlist []string, depth int) {
	if depth > MaxDepth {
		return
	}

	for _, prefix := range worldlist {
		if prefix != "" && len(lookupHost(prefix+"."+domain)) > 0 {
			fmt.Println(style.SuccessStyle.Render(prefix + "." + domain))
			if commandBruteForce.Recurse {
				findSubDomains(prefix+"."+domain, worldlist, depth+1)
			}
		}
	}
}
