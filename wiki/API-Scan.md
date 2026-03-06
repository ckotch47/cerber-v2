# API Scan

`api scan` проверяет endpoint-ы из OpenAPI спецификации и показывает HTTP статус по каждому методу/пути.

## Базовый запуск

```bash
cerber api scan --spec ./openapi.json --host https://api.example.com
```

## Источник спецификации

- Локальный файл: `--spec ./openapi.json`
- URL: `--spec https://api.example.com/openapi.json`

## Авторизация

JWT:

```bash
cerber api scan --spec ./openapi.json --host https://api.example.com --jwt "$JWT"
```

API key:

```bash
cerber api scan \
  --spec ./openapi.json \
  --host https://api.example.com \
  --api-key-header X-API-Key \
  --api-key "$API_KEY"
```

## Фильтрация статусов

Показать только:

```bash
cerber api scan --spec ./openapi.json --host https://api.example.com --show 200,201
```

Исключить:

```bash
cerber api scan --spec ./openapi.json --host https://api.example.com --exclude 401,403,404
```

`--show` и `--exclude` взаимоисключающие.

## Полезные флаги

- `--request-timeout 10` — timeout запроса
- `--show-errors=false` — скрыть сетевые ошибки
- `--spec-auth=false` — не отправлять auth заголовки при загрузке spec по URL
