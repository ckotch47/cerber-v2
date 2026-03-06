package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"cerber/internal/recon"
	"cerber/internal/style"
	"cerber/internal/utils"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Выполняет поиск поддоменов по списку из файла",
	RunE:  FindHost,
}

var commandBruteForce utils.BruteForceType

func init() {
	findCmd.Flags().StringVarP(
		&commandBruteForce.WorldList,
		"wordlist",
		"w",
		"",
		"Файл со списком",
	)
	findCmd.Flags().StringVar(
		&commandBruteForce.WorldList,
		"worldlis",
		"",
		"Устаревший алиас для --wordlist",
	)
	_ = findCmd.Flags().MarkDeprecated("worldlis", "use --wordlist instead")
	findCmd.Flags().BoolVarP(
		&commandBruteForce.Recurse,
		"recurse",
		"r",
		false,
		"Включить рекурсию для брутфорса",
	)
	findCmd.Flags().IntVar(
		&commandBruteForce.MaxDepth,
		"max-depth",
		2,
		"Максимальная глубина рекурсии для поиска поддоменов",
	)
	findCmd.Flags().IntVarP(
		&commandBruteForce.Concurrency,
		"concurrency",
		"c",
		20,
		"Количество параллельных DNS-запросов",
	)
}

func FindHost(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("не указан домен")
	}
	if commandBruteForce.WorldList == "" {
		return fmt.Errorf("файл со списком не указан")
	}

	domain, err := cleanDomain(args[0])
	if err != nil {
		return err
	}
	domainList := utils.ReadFile(commandBruteForce.WorldList)
	if len(domainList) == 0 {
		return fmt.Errorf("файл со списком пустой или не удалось прочитать")
	}

	scanner := recon.NewSubdomainScanner(nil)
	found := scanner.Collect(
		domain,
		domainList,
		commandBruteForce.Recurse,
		commandBruteForce.MaxDepth,
		commandBruteForce.Concurrency,
	)
	for _, subdomain := range found {
		fmt.Println(style.SuccessStyle.Render(subdomain))
	}
	return nil
}
