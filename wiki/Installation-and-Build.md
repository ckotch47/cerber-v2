# Installation and Build

## Требования

- Go 1.22+
- Доступ к DNS/HTTP из среды запуска

## Сборка

```bash
cd /Users/blant/GoLangProject/lessons/cerber
go build -o cerber .
```

## Установка в PATH

```bash
cd /Users/blant/GoLangProject/lessons/cerber
go build -o /Users/blant/go/bin/cerber .
```

Проверка:

```bash
which cerber
cerber version
```

## Локальные тесты

```bash
cd /Users/blant/GoLangProject/lessons/cerber
go test ./...
```
