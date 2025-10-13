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
   cerber find example.com -w wordlist.txt [-r]
   ```
   Flags:
   - `-w, --worldlis`: Path to the wordlist file (required)
   - `-r, --recurse`: Enable recursive subdomain enumeration

3. **find admin** - Search for admin panels
   ```bash
   cerber find admin example.com -w wordlist.txt [-e status_codes]
   ```
   Flags:
   - `-w, --worldlis`: Path to the wordlist file (required)
   - `-e, --exclude`: Status codes to exclude from results (can be specified multiple times)
   
   Example:
   ```bash
   cerber find admin example.com -w paths.txt -e 404 -e 500
   ```

4. **version** - Show application version
   ```bash
   cerber version
   ```

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
   cerber find admin example.com -w admin-paths.txt -e 404
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