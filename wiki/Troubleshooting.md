# Troubleshooting

## `which cerber` показывает не тот бинарник

Проблема:
- В системе может быть несколько `cerber` (`/usr/local/bin`, `~/go/bin`).

Проверка:

```bash
which cerber
cerber version
/usr/local/bin/cerber version
/Users/blant/go/bin/cerber version
```

Решение:
- Обновить бинарник именно по пути, который выдает `which cerber`.
- При необходимости поправить `PATH` порядок.

## Команда `find`/`find path` ругается на wordlist

Проверьте, что передан флаг:

```bash
-w ./file.txt
```

и файл реально существует.

## `api scan` не загружает `--spec` по URL

Варианты:
- URL недоступен
- Нужна авторизация

Проверьте:
- `--jwt` или `--api-key-*`
- `--spec-auth=true`
- корректность TLS/прокси в окружении

## Много ошибок сети в выводе `api scan`

Скрыть можно так:

```bash
--show-errors=false
```

Или увеличить timeout:

```bash
--request-timeout 20
```
