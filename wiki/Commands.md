# Commands

## Глобально

- `--lang auto|ru|en` — выбор языка, по умолчанию `auto`

## `look`

Назначение:
- Для домена выводит IP.
- Для IP выводит reverse DNS имена.

Примеры:

```bash
cerber look qatoria.ru
cerber look 8.8.8.8
```

## `find`

Назначение:
- Поиск поддоменов из wordlist.
- Опциональная рекурсия.

Флаги:
- `-w, --wordlist` (required)
- `-r, --recurse`
- `--max-depth` (default: `2`)
- `-c, --concurrency` (default: `20`)
- `--worldlis` (deprecated alias)

Пример:

```bash
cerber find ecofamilyschool.ru -w ./subdomains-500.txt -r --max-depth 2 -c 20
```

## `find path`

Назначение:
- Проверка путей/директорий из wordlist.
- Если схема не указана: сначала `https://`, потом fallback на `http://`.

Флаги:
- `-w, --wordlist` (required)
- `-e, --exclude` (можно несколько)
- `-t, --timeout` (delay между запросами, default: `5`)
- `--request-timeout` (HTTP timeout, default: `10`)
- `--worldlis` (deprecated alias)

Пример:

```bash
cerber find path app.dev.ecofamilyschool.ru -w ./hidden-path.txt -e 404 -e 429 -t 1 --request-timeout 10
```

## `google links`

Назначение:
- Генерация ссылок Google dorks.

Флаги:
- `--mode all|1,2,3...` (default: `all`)

Пример:

```bash
cerber google links qatoria.ru --mode all
```

## `api scan`

Назначение:
- Чтение OpenAPI и проверка доступности endpoint-ов.

Required:
- `--spec <url|file>`
- `--host <url>`

Флаги:
- `--show <codes>`
- `--exclude <codes>`
- `--request-timeout <sec>`
- `--show-errors`
- `--jwt <token>`
- `--api-key-header <name>`
- `--api-key <value>`
- `--spec-auth` (default: `true`)

Ограничение:
- `--show` и `--exclude` нельзя использовать одновременно.

Пример:

```bash
cerber api scan --spec ./openapi.json --host https://api.example.com --jwt "$JWT" --show 200,201
```

## `version`

```bash
cerber version
```
