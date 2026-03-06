# Architecture

Текущая структура разделена на слои:

- `internal/command`  
  CLI-команды, парсинг флагов, валидация ввода, вывод результата.

- `internal/recon`  
  Core-логика сканирования:
  - `subdomain_service.go`
  - `path_service.go`
  - `apiscan_service.go`
  - `google_links_service.go`

- `internal/dns`  
  DNS helper-функции.

- `internal/i18n`  
  Локализация (`ru/en`) и выбор языка (`auto`).

- `internal/utils`  
  Вспомогательные функции (чтение wordlist и пр.).

## Принципы

- Команды не содержат heavy-логики, только orchestration.
- Сканеры и сервисы покрываются unit-тестами.
- Ошибки пробрасываются через `RunE` и обрабатываются в root.
