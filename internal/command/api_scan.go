package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"cerber/internal/recon"
	"cerber/internal/style"
)

var (
	apiSpecSource string
	apiHost       string
	apiShow       string
	apiExclude    string
	apiTimeout    int
	apiJWT        string
	apiKeyHeader  string
	apiKey        string
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "API scanning tools",
}

var apiScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Сканировать OpenAPI-спецификацию на доступность эндпоинтов",
	RunE:  runAPIScan,
}

func init() {
	apiScanCmd.Flags().StringVar(
		&apiSpecSource,
		"spec",
		"",
		"Путь или URL до openapi.json",
	)
	apiScanCmd.Flags().StringVar(
		&apiHost,
		"host",
		"",
		"Базовый URL API, куда отправлять запросы",
	)
	apiScanCmd.Flags().StringVar(
		&apiShow,
		"show",
		"",
		"Показывать только эти коды статуса (например: 200,201)",
	)
	apiScanCmd.Flags().StringVar(
		&apiExclude,
		"exclude",
		"",
		"Исключить эти коды статуса (например: 401,403)",
	)
	apiScanCmd.Flags().IntVar(
		&apiTimeout,
		"request-timeout",
		10,
		"Таймаут HTTP запроса в секундах",
	)
	apiScanCmd.Flags().StringVar(
		&apiJWT,
		"jwt",
		"",
		"JWT токен для Authorization: Bearer <token>",
	)
	apiScanCmd.Flags().StringVar(
		&apiKeyHeader,
		"api-key-header",
		"",
		"Имя заголовка для API key (например: X-API-Key)",
	)
	apiScanCmd.Flags().StringVar(
		&apiKey,
		"api-key",
		"",
		"Значение API key",
	)
	_ = apiScanCmd.MarkFlagRequired("spec")
	_ = apiScanCmd.MarkFlagRequired("host")

	apiCmd.AddCommand(apiScanCmd)
}

func runAPIScan(_ *cobra.Command, _ []string) error {
	host := recon.NormalizeBaseURL(apiHost)
	if host == "" {
		return fmt.Errorf("host пуст после нормализации")
	}

	showSet, err := recon.ParseStatusCodes(apiShow)
	if err != nil {
		return fmt.Errorf("невалидный --show: %w", err)
	}
	excludeSet, err := recon.ParseStatusCodes(apiExclude)
	if err != nil {
		return fmt.Errorf("невалидный --exclude: %w", err)
	}

	headers, err := recon.BuildAuthHeaders(apiJWT, apiKeyHeader, apiKey)
	if err != nil {
		return err
	}

	scanner := recon.NewAPIScanner(apiTimeout, headers)
	spec, err := scanner.LoadSpec(apiSpecSource)
	if err != nil {
		return err
	}

	results := scanner.Scan(host, spec)
	for _, result := range results {
		if result.Err != nil {
			fmt.Println(style.NotFoundStyle.Render(fmt.Sprintf("[%s] %s : error (%v)", result.Method, result.Path, result.Err)))
			continue
		}
		if !recon.ShouldIncludeStatus(result.StatusCode, showSet, excludeSet) {
			continue
		}
		if result.StatusCode >= 400 {
			fmt.Println(style.NotFoundStyle.Render(fmt.Sprintf("[%s] %s : %d", result.Method, result.Path, result.StatusCode)))
		} else {
			fmt.Println(style.SuccessStyle.Render(fmt.Sprintf("[%s] %s : %d", result.Method, result.Path, result.StatusCode)))
		}
	}
	return nil
}
