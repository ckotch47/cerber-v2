package i18n

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
)

var currentLang atomic.Value

func init() {
	currentLang.Store("ru")
}

var messages = map[string]map[string]string{
	"ru": {
		"lang_flag_desc":                      "Язык сообщений: auto|ru|en",
		"err_lang_not_supported":              "неподдерживаемый язык: %s (доступно: ru,en)",
		"err_domain_required":                 "домен не указан",
		"err_domain_empty_after_normalize":    "домен пуст после нормализации",
		"err_input_domain_or_ip":              "введите домен или IP-адрес",
		"msg_not_found":                       "Не найдено",
		"err_wordlist_required":               "файл со списком не указан",
		"err_read_wordlist":                   "не удалось прочитать wordlist: %w",
		"err_host_empty_after_normalize":      "host пуст после нормализации",
		"err_show_exclude_mutually_exclusive": "флаги --show и --exclude нельзя использовать одновременно",
		"err_invalid_show":                    "невалидный --show: %w",
		"err_invalid_exclude":                 "невалидный --exclude: %w",
		"msg_request_error_prefix":            "Ошибка запроса:",
		"msg_version":                         "Версия: %s",
		"err_file_empty":                      "файл пустой",
		"err_api_key_requires_header":         "для --api-key требуется --api-key-header",
		"err_api_key_header_requires_key":     "для --api-key-header требуется --api-key",
		"err_invalid_openapi_json":            "невалидный openapi json: %w",
		"err_openapi_paths_missing":           "в спецификации отсутствует paths",
		"err_invalid_status_code":             "невалидный статус-код %q",
		"err_status_code_out_of_range":        "статус-код вне диапазона 100..999: %d",
		"err_spec_load_status":                "не удалось загрузить спецификацию: status %d",
		"err_invalid_google_mode":             "невалидный mode %q: ожидается число или all",
		"err_google_mode_out_of_range":        "mode %d вне диапазона 1..%d",
		"cmd_short_root":                      "CLI для DNS recon и поиска скрытых путей",
		"cmd_long_root":                       "Cerber — инструмент для DNS lookup, поиска поддоменов и скрытых путей.",
		"cmd_short_api":                       "Инструменты сканирования API",
		"cmd_short_google":                    "Инструменты для генерации Google dork ссылок",
		"cmd_short_find":                      "Выполняет поиск поддоменов по списку из файла",
		"cmd_short_find_path":                 "Поиск админ панелей",
		"cmd_short_look":                      "Найти IP по домену или доменные имена по IP",
		"cmd_short_version":                   "Показать версию приложения",
		"cmd_long_look":                       "Примеры:\n  cerber look http://example.com — найти IP по домену\n  cerber look 8.8.8.8 — найти доменные имена по IP",
		"cmd_short_google_links":              "Сгенерировать dork-ссылки для домена",
		"cmd_short_api_scan":                  "Сканировать OpenAPI-спецификацию на доступность эндпоинтов",
		"flag_wordlist_desc":                  "Файл со списком",
		"flag_worldlis_desc":                  "Устаревший алиас для --wordlist",
		"flag_recurse_desc":                   "Включить рекурсию для брутфорса",
		"flag_max_depth_desc":                 "Максимальная глубина рекурсии для поиска поддоменов",
		"flag_concurrency_desc":               "Количество параллельных DNS-запросов",
		"flag_exclude_desc":                   "Статусы для исключения",
		"flag_timeout_desc":                   "Время задержки между запросами в секундах (по умолчанию 5 сек)",
		"flag_request_timeout_desc":           "Таймаут HTTP запроса в секундах",
		"flag_google_mode_desc":               "Режимы через запятую (например: 1,5,12) или all",
		"flag_api_spec_desc":                  "Путь или URL до openapi.json",
		"flag_api_host_desc":                  "Базовый URL API, куда отправлять запросы",
		"flag_api_show_desc":                  "Показывать только эти коды статуса (например: 200,201)",
		"flag_api_exclude_desc":               "Исключить эти коды статуса (например: 401,403)",
		"flag_api_jwt_desc":                   "JWT токен для Authorization: Bearer <token>",
		"flag_api_key_header_desc":            "Имя заголовка для API key (например: X-API-Key)",
		"flag_api_key_desc":                   "Значение API key",
		"flag_api_spec_auth_desc":             "Передавать auth-заголовки при загрузке --spec по URL",
		"flag_api_show_errors_desc":           "Показывать ошибки сетевых запросов",
	},
	"en": {
		"lang_flag_desc":                      "Language for messages: auto|ru|en",
		"err_lang_not_supported":              "unsupported language: %s (available: ru,en)",
		"err_domain_required":                 "domain is required",
		"err_domain_empty_after_normalize":    "domain is empty after normalization",
		"err_input_domain_or_ip":              "enter a domain or IP address",
		"msg_not_found":                       "Not found",
		"err_wordlist_required":               "wordlist file is required",
		"err_read_wordlist":                   "failed to read wordlist: %w",
		"err_host_empty_after_normalize":      "host is empty after normalization",
		"err_show_exclude_mutually_exclusive": "--show and --exclude cannot be used together",
		"err_invalid_show":                    "invalid --show: %w",
		"err_invalid_exclude":                 "invalid --exclude: %w",
		"msg_request_error_prefix":            "Request error:",
		"msg_version":                         "Version: %s",
		"err_file_empty":                      "file is empty",
		"err_api_key_requires_header":         "--api-key requires --api-key-header",
		"err_api_key_header_requires_key":     "--api-key-header requires --api-key",
		"err_invalid_openapi_json":            "invalid openapi json: %w",
		"err_openapi_paths_missing":           "openapi spec has no paths",
		"err_invalid_status_code":             "invalid status code %q",
		"err_status_code_out_of_range":        "status code out of range 100..999: %d",
		"err_spec_load_status":                "failed to load spec: status %d",
		"err_invalid_google_mode":             "invalid mode %q: expected number or all",
		"err_google_mode_out_of_range":        "mode %d out of range 1..%d",
		"cmd_short_root":                      "CLI for DNS recon and hidden path discovery",
		"cmd_long_root":                       "Cerber is a CLI tool for DNS lookup, subdomain discovery, and hidden path scanning.",
		"cmd_short_api":                       "API scanning tools",
		"cmd_short_google":                    "Google dork helpers",
		"cmd_short_find":                      "Find subdomains from a wordlist",
		"cmd_short_find_path":                 "Find hidden/admin paths",
		"cmd_short_look":                      "Resolve domain IPs or reverse lookup by IP",
		"cmd_short_version":                   "Show application version",
		"cmd_long_look":                       "Examples:\n  cerber look http://example.com - resolve domain to IP\n  cerber look 8.8.8.8 - reverse lookup domain names by IP",
		"cmd_short_google_links":              "Generate dork links for a domain",
		"cmd_short_api_scan":                  "Scan OpenAPI spec for endpoint availability",
		"flag_wordlist_desc":                  "Wordlist file path",
		"flag_worldlis_desc":                  "Deprecated alias for --wordlist",
		"flag_recurse_desc":                   "Enable recursive brute-force",
		"flag_max_depth_desc":                 "Maximum recursion depth for subdomain discovery",
		"flag_concurrency_desc":               "Number of parallel DNS requests",
		"flag_exclude_desc":                   "Status codes to exclude",
		"flag_timeout_desc":                   "Delay between requests in seconds (default 5)",
		"flag_request_timeout_desc":           "HTTP request timeout in seconds",
		"flag_google_mode_desc":               "Comma-separated modes (e.g. 1,5,12) or all",
		"flag_api_spec_desc":                  "Path or URL to openapi.json",
		"flag_api_host_desc":                  "Base API URL for outgoing requests",
		"flag_api_show_desc":                  "Show only these status codes (e.g. 200,201)",
		"flag_api_exclude_desc":               "Exclude these status codes (e.g. 401,403)",
		"flag_api_jwt_desc":                   "JWT token for Authorization: Bearer <token>",
		"flag_api_key_header_desc":            "Header name for API key (e.g. X-API-Key)",
		"flag_api_key_desc":                   "API key value",
		"flag_api_spec_auth_desc":             "Send auth headers when loading --spec from URL",
		"flag_api_show_errors_desc":           "Show network request errors",
	},
}

func SetLang(lang string) error {
	normalized := normalizeLang(lang)
	if normalized == "" {
		return fmt.Errorf(TWithLang("err_lang_not_supported", "ru", lang))
	}
	currentLang.Store(normalized)
	return nil
}

func DetectLang() string {
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if detected := detectFromLocale(os.Getenv(key)); detected != "" {
			return detected
		}
	}
	return "en"
}

func Lang() string {
	if v, ok := currentLang.Load().(string); ok {
		return v
	}
	return "ru"
}

func T(key string, args ...any) string {
	return TWithLang(key, Lang(), args...)
}

func TWithLang(key string, lang string, args ...any) string {
	bundle, ok := messages[lang]
	if !ok {
		bundle = messages["ru"]
	}
	msg, ok := bundle[key]
	if !ok {
		msg = messages["ru"][key]
	}
	if msg == "" {
		msg = key
	}
	if len(args) == 0 {
		return msg
	}
	return fmt.Sprintf(msg, args...)
}

func normalizeLang(lang string) string {
	value := strings.TrimSpace(strings.ToLower(lang))
	if value == "auto" {
		return DetectLang()
	}
	if strings.HasPrefix(value, "en") {
		return "en"
	}
	if strings.HasPrefix(value, "ru") {
		return "ru"
	}
	return ""
}

func detectFromLocale(locale string) string {
	value := strings.TrimSpace(strings.ToLower(locale))
	if strings.HasPrefix(value, "ru") {
		return "ru"
	}
	if strings.HasPrefix(value, "en") {
		return "en"
	}
	return ""
}
