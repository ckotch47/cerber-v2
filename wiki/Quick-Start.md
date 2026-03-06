# Quick Start

## 1. Установка

```bash
cd /Users/blant/GoLangProject/lessons/cerber
go build -o cerber .
```

Или в `~/go/bin`:

```bash
cd /Users/blant/GoLangProject/lessons/cerber
go build -o /Users/blant/go/bin/cerber .
```

## 2. Базовые команды

```bash
cerber look qatoria.ru
cerber find qatoria.ru -w ./subdomains-500.txt
cerber find path app.dev.ecofamilyschool.ru -w ./hidden-path.txt -e 404
cerber google links qatoria.ru --mode all
cerber api scan --spec ./openapi.json --host https://api.example.com --show 200
```

## 3. Проверка версии

```bash
cerber version
```
