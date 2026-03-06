# Cerber

Cerber is a command-line tool for domain reconnaissance and admin panel discovery. It provides functionality for subdomain enumeration, DNS lookups, and admin panel scanning.

## Features

- Subdomain enumeration with recursive scanning
- DNS lookup for domain resolution
- Admin panel discovery with customizable status code filtering
- Support for wordlist-based scanning

## Installation

```bash
go install
```

## Usage

### Basic Commands

```bash
cerber [command] [flags]
```

### Available Commands

1. **look** - Find IP addresses for a domain
   ```bash
   cerber look example.com
   ```

2. **find** - Perform subdomain enumeration
   ```bash
   cerber find example.com -w wordlist.txt [-r] [--max-depth N] [-c N]
   ```
   Flags:
   - `-w, --wordlist`: Path to the wordlist file (required)
   - `-r, --recurse`: Enable recursive subdomain enumeration
   - `--max-depth`: Max recursion depth for recursive mode (default: 2)
   - `-c, --concurrency`: Number of parallel DNS lookups (default: 20)

3. **find path** - Search for hidden path 
   ```bash
   cerber find path example.com -w wordlist.txt [-e status_codes] [--request-timeout seconds]
   ```
   Flags:
   - `-w, --wordlist`: Path to the wordlist file (required)
   - `-e, --exclude`: Status codes to exclude from results (can be specified multiple times)
   - `--request-timeout`: HTTP timeout per request in seconds (default: 10)
   
   Example:
   ```bash
   cerber find path example.com -w paths.txt -e 404 -e 500
   ```

4. **version** - Show application version
   ```bash
   cerber version
   ```

5. **google links** - Generate Google dork links
   ```bash
   cerber google links example.com [--mode all|1,5,12]
   ```
   Flags:
   - `--mode`: Modes list (comma-separated) or `all` (default: `all`)

6. **api scan** - Scan OpenAPI endpoints
   ```bash
   cerber api scan --spec ./openapi.json --host https://api.example.com [--jwt TOKEN | --api-key-header X-API-Key --api-key VALUE] [--show 200,201] [--exclude 401,403]
   ```
   Flags:
   - `--spec`: Path or URL to `openapi.json` (required)
   - `--host`: Base API URL for requests (required)
   - `--jwt`: JWT token (`Authorization: Bearer ...`)
   - `--api-key-header`: API key header name
   - `--api-key`: API key value
   - `--show`: Show only these status codes
   - `--exclude`: Hide these status codes
   - `--request-timeout`: HTTP timeout per request in seconds (default: 10)

## Examples

1. Basic DNS lookup:
   ```bash
   cerber look example.com
   ```

2. Subdomain enumeration:
   ```bash
   cerber find example.com -w subdomains.txt
   ```

3. Recursive subdomain scanning:
   ```bash
   cerber find example.com -w subdomains.txt -r
   ```

4. Admin panel discovery (excluding 404 responses):
   ```bash
   cerber find path example.com -w admin-paths.txt -e 404
   ```

## Wordlist Format

- For subdomain enumeration: One subdomain prefix per line
- For admin panel discovery: One path per line

## Note

The tool automatically handles various domain formats:
- Removes "http://" and "https://" prefixes
- Removes "www." prefix
- Removes trailing slashes

## Version

Current version: v0.0.1a

## License

[Add your license information here]
