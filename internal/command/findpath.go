package command

import (
	"cerber/internal/style"
	"cerber/internal/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var commandPathFinder utils.AdminFindeType

type StringSlice map[string]bool

var findPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Поиск админ панелей",
	Long:  `Поиск админ панелей`,
	Run:   FindHiddenPath,
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

func FindHiddenPath(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println(style.NotFoundStyle.Render("Не указан домен"))
		return
	}
	if commandPathFinder.WorldList == "" {
		fmt.Println(style.NotFoundStyle.Render("Файл со списком не найден"))
		return
	}
	domain := normalizeBaseURL(args[0])
	allowFallback := !hasHTTPPrefix(args[0])
	worldList := utils.ReadFile(commandPathFinder.WorldList)
	client := &http.Client{}
	exclude := arrayToMap(commandPathFinder.Exclude)

	for _, path := range worldList {
		get(client, domain, path, allowFallback, exclude)
	}
}

func get(client *http.Client, baseURL, path string, allowFallback bool, exclude StringSlice) {
	target := joinURL(baseURL, path)
	resp, actualURL, err := requestWithFallback(client, target, commandPathFinder.RequestTimeout, allowFallback)
	if err != nil {
		fmt.Println("Ошибка запроса:", style.NotFoundStyle.Render(actualURL+" -> "+err.Error()))
		return
	}
	defer resp.Body.Close()
	printRespStatus(resp.StatusCode, path, exclude)
	time.Sleep(time.Duration(commandPathFinder.Timeout) * time.Second)
}

func printRespStatus(statusCode int, path string, exclude StringSlice) {
	var strResp string

	if exclude[strconv.Itoa(statusCode)] {
		return
	}

	if statusCode >= 400 {
		strResp = path + " : " + style.NotFoundStyle.Render(strconv.Itoa(statusCode))
	} else {
		strResp = path + " : " + style.SuccessStyle.Render(strconv.Itoa(statusCode))
	}

	fmt.Println(strResp) // выводим статус ответа
}

func arrayToMap(exclude []string) StringSlice {
	if exclude == nil {
		return StringSlice{}
	}
	res := make(StringSlice)

	for _, code := range exclude {
		if code != "" {
			res[code] = true
		}
	}
	return res
}

func normalizeBaseURL(input string) string {
	trimmed := strings.TrimSpace(strings.TrimSuffix(input, "/"))
	if trimmed == "" {
		return ""
	}
	if hasHTTPPrefix(trimmed) {
		return trimmed
	}
	return "https://" + trimmed
}

func hasHTTPPrefix(input string) bool {
	return strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")
}

func joinURL(baseURL, path string) string {
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func requestWithFallback(client *http.Client, target string, timeoutSeconds int, allowFallback bool) (*http.Response, string, error) {
	resp, err := requestOnce(client, target, timeoutSeconds)
	if err == nil {
		return resp, target, nil
	}
	if !allowFallback || !strings.HasPrefix(target, "https://") {
		return nil, target, err
	}

	fallbackTarget := "http://" + strings.TrimPrefix(target, "https://")
	fallbackResp, fallbackErr := requestOnce(client, fallbackTarget, timeoutSeconds)
	if fallbackErr != nil {
		return nil, fallbackTarget, fallbackErr
	}
	return fallbackResp, fallbackTarget, nil
}

func requestOnce(client *http.Client, target string, timeoutSeconds int) (*http.Response, error) {
	reqTimeout := time.Duration(timeoutSeconds) * time.Second
	if reqTimeout <= 0 {
		reqTimeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
