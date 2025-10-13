package command

import (
	"cerber/internal/style"
	"cerber/internal/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var commandAdminFinder utils.AdminFindeType
type StringSlice map[string]bool

var findAdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Поиск админ панелей",
	Long:  `Поиск админ панелей`,
	Run: FindAdminPanels,
}

func init() {
	findAdminCmd.Flags().StringVarP(
		&commandAdminFinder.WorldList,
		"worldlis", 
		"w", 
		"", 
		"Файл со списком",
	)
	findAdminCmd.Flags().StringArrayVarP(
		&commandAdminFinder.Exclude,
		"exclude", 
		"e", 
		[]string{}, 
		"статусы для исключения",
	)
}

func FindAdminPanels(cmd *cobra.Command, args []string) {
	if(len(args) == 0){
		fmt.Println(style.NotFoundStyle.Render("Не указан домен"))
		return
	}
	if(commandAdminFinder.WorldList == ""){
		fmt.Println(style.NotFoundStyle.Render("Файл со списком не найден"))
		return
	}
	domain := strings.TrimSuffix(args[0], "/")
	worldList := utils.ReadFile(commandAdminFinder.WorldList)


	for _, path := range worldList {
		get(domain + "/" + path, path)
	}
}

func get(url, path string) {
		// Делаем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка запроса:", style.NotFoundStyle.Render(err.Error()))
		return
	}
	defer resp.Body.Close() // обязательно закрываем тело ответа
	printRespStatus(resp, path)
}

func printRespStatus(resp *http.Response, path string) {
	var strResp string

	if arrayToMap(commandAdminFinder.Exclude)[strconv.Itoa(resp.StatusCode)]{
		return
	}

	if resp.StatusCode > 400 {
		strResp = path + " : " + style.NotFoundStyle.Render(strconv.Itoa(resp.StatusCode))
	}else{
		strResp = path + " : " + style.SuccessStyle.Render(strconv.Itoa(resp.StatusCode))
	}

	fmt.Println(strResp) // выводим статус ответа
}

func arrayToMap(exclude []string) StringSlice{
	if(exclude == nil){
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