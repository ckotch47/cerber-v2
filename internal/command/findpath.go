package command

import (
	"cerber/internal/style"
	"cerber/internal/utils"
	"fmt"
	"io"
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
		"worldlis",
		"w",
		"",
		"Файл со списком",
	)
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
	domain := strings.TrimSuffix(args[0], "/")
	worldList := utils.ReadFile(commandPathFinder.WorldList)

	for _, path := range worldList {
		get(domain+"/"+path, path)
	}
}

func get(url, path string) {
	// Делаем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка запроса:", style.NotFoundStyle.Render(err.Error()))
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body) // обязательно закрываем тело ответа
	printRespStatus(resp, path)
	time.Sleep(time.Duration(commandPathFinder.Timeout) * time.Second)
}

func printRespStatus(resp *http.Response, path string) {
	var strResp string

	if arrayToMap(commandPathFinder.Exclude)[strconv.Itoa(resp.StatusCode)] {
		return
	}

	if resp.StatusCode > 400 {
		strResp = path + " : " + style.NotFoundStyle.Render(strconv.Itoa(resp.StatusCode))
	} else {
		strResp = path + " : " + style.SuccessStyle.Render(strconv.Itoa(resp.StatusCode))
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
