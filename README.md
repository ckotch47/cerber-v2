# Cerber

Cerber is a CLI tool for recon tasks: DNS lookup, subdomain discovery, hidden path scan, Google dork links, and OpenAPI endpoint availability checks.

## Quick Start

```bash
cd /Users/blant/GoLangProject/lessons/cerber
go build -o /Users/blant/go/bin/cerber .
cerber version
```

Global flag:
- `--lang auto|ru|en` (default: `auto`)

## Main Commands

```bash
cerber look <domain-or-ip>
cerber find <domain> -w <wordlist> [-r] [--max-depth N] [-c N]
cerber find path <host-or-url> -w <wordlist> [-e CODE] [-t SEC] [--request-timeout SEC]
cerber google links <domain> [--mode all|1,5,12]
cerber api scan --spec <url-or-file> --host <base-url> [--jwt TOKEN | --api-key-header H --api-key V]
```

## Full Docs (Wiki)

- https://github.com/ckotch47/cerber-v2/wiki
- https://github.com/ckotch47/cerber-v2/wiki/Commands
- https://github.com/ckotch47/cerber-v2/wiki/API-Scan
- https://github.com/ckotch47/cerber-v2/wiki/Troubleshooting

## Version

Current version: `v0.0.1`

## Deprecations

- `--worldlis` is a deprecated alias for `--wordlist`.
- It stays available in `v0.0.1` for compatibility.
- Planned removal: next minor release after migration window.

## Release Notes

- `MVP_CHECKLIST.md`
- `CHANGELOG_MIGRATION.md`
- `VERSIONING.md`
- `RELEASE_CHECKLIST.md`
