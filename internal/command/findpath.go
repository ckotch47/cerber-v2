package command

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"cerber/internal/recon"
	"cerber/internal/style"
	"cerber/internal/utils"
)

var commandPathFinder utils.AdminFindeType

var findPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Поиск админ панелей",
	Long:  `Поиск админ панелей`,
	RunE:  FindHiddenPath,
}

func init() {
	findPathCmd.Flags().StringVarP(
		&commandPathFinder.WorldList,
		"wordlist",
		"w",
		"",
		"Файл со списком",
	)
	findPathCmd.Flags().StringVar(
		&commandPathFinder.WorldList,
		"worldlis",
		"",
		"Устаревший алиас для --wordlist",
	)
	_ = findPathCmd.Flags().MarkDeprecated("worldlis", "use --wordlist instead")
	findPathCmd.Flags().StringArrayVarP(
		&commandPathFinder.Exclude,
		"exclude",
		"e",
		[]string{},
		"статусы для исключения",
	)
	findPathCmd.Flags().IntVarP(
		&commandPathFinder.Timeout,
		"timeout",
		"t",
		5,
		"Время задержки между запросами в секндах (по умолчанию 5 сек)",
	)
	findPathCmd.Flags().IntVar(
		&commandPathFinder.RequestTimeout,
		"request-timeout",
		10,
		"Таймаут HTTP запроса в секундах",
	)
}

func FindHiddenPath(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("не указан домен")
	}
	if commandPathFinder.WorldList == "" {
		return fmt.Errorf("файл со списком не указан")
	}

	domain := recon.NormalizeBaseURL(args[0])
	if domain == "" {
		return fmt.Errorf("домен пуст после нормализации")
	}
	worldList := utils.ReadFile(commandPathFinder.WorldList)
	if len(worldList) == 0 {
		return fmt.Errorf("файл со списком пустой или не удалось прочитать")
	}

	scanner := recon.NewPathScanner(
		commandPathFinder.RequestTimeout,
		commandPathFinder.Timeout,
		commandPathFinder.Exclude,
		!recon.HasHTTPPrefix(args[0]),
	)
	results := scanner.Scan(domain, worldList)

	for _, result := range results {
		if result.Err != nil {
			fmt.Println("Ошибка запроса:", style.NotFoundStyle.Render(result.URL+" -> "+result.Err.Error()))
			continue
		}
		printRespStatus(result.StatusCode, result.Path, scanner.IsExcluded(result.StatusCode))
	}

	return nil
}

func printRespStatus(statusCode int, path string, excluded bool) {
	var strResp string

	if excluded {
		return
	}

	if statusCode >= 400 {
		strResp = path + " : " + style.NotFoundStyle.Render(strconv.Itoa(statusCode))
	} else {
		strResp = path + " : " + style.SuccessStyle.Render(strconv.Itoa(statusCode))
	}

	fmt.Println(strResp)
}
