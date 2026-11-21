# 🐷 steamer

**steamer** is a high-speed, interactive CLI and TUI for managing your Porkbun domains. Built with Go, Cobra, and Charm's Bubble Tea, it's designed to make DNS management feel less like a chore and more like a breeze.

Whether you're delegating a subdomain with `NS` records, adding a `CNAME` for a new project, configuring email forwarding, or browsing your portfolio, **steamer** has you covered.

## ✨ Features

- **🚀 Quick DNS Management:** Create `A`, `AAAA`, `CNAME`, `NS` (subdomain delegation), and `TXT` records, or safely delete records (`rm`) with interactive confirmation (`--yes` for automation).
- **🔍 Domain Search & TLD Pricing:** Check domain availability across multiple TLDs (`search`) and browse cached TLD pricing (`list-tlds`).
- **📧 Email Forwarding:** List, add, and remove Porkbun email forwarding rules (`email list`, `email add`, `email rm`).
- **🤖 Agent-Friendly Output:** Supports `--json` flags across query and mutation commands for deterministic scripting and AI agent workflows.
- **🎮 Interactive TUI:** A terminal interface to browse your domains, DNS records, and TLD pricing tables.
- **🔐 Secure Config:** Supports XDG-style configuration, environment variables, and `.env` fallback with config permission warnings.
- **📖 API Reference:** Includes a local copy of the Porkbun v3 API documentation.

## 🛠️ Installation

### Homebrew (macOS & Linux)

```bash
brew install ghchinoy/tap/steamer
```

### Go Install

```bash
go install github.com/ghchinoy/steamer@latest
```

This builds the binary and places it in your `$GOPATH/bin` (typically `~/go/bin`). Ensure that directory is in your system's `PATH`.

### From Source

```bash
go install .
```

## ⚙️ Configuration

**steamer** looks for your Porkbun API credentials in several places. You'll need an API Key and Secret from [Porkbun](https://porkbun.com/account/api).

### 1. The Pro Way (Config File)
Create a config file at `~/.config/steamer/config.yaml` (recommended permissions `chmod 600 ~/.config/steamer/config.yaml`):

```yaml
apikey: pk1_your_api_key
secretapikey: sk1_your_secret_key
```

### 2. The Quick Way (Environment Variables)
```bash
export PORKBUN_APIKEY=pk1_...
export PORKBUN_SECRETAPIKEY=sk1_...
```

### 3. The Dev Way (.env Fallback)
If credentials are not provided by a config file or environment variables, **steamer** falls back to a `.env` file in the current directory:

```env
API_KEY=pk1_your_api_key
API_SECRET=sk1_your_secret_key
```

## 🚀 Usage

### The Terminal UI (TUI)
Just run `steamer tui` and enjoy the ride. Use `j`/`k` to navigate and `enter` to dive into records.

```bash
# Start the TUI
steamer tui

# Jump straight to a domain
steamer tui -d example.com
```

### Command Line Interface

#### Domain Discovery & Information
```bash
# Check domain availability and pricing
steamer search mynewidea.com

# Search a phrase against specific TLDs
steamer search mynewidea --tlds com,dev,app

# List all supported TLDs and their prices (cached locally for 7 days)
steamer list-tlds
steamer list-tlds --force

# List all your domains (add --json for scripting)
steamer list-domains
steamer list-domains --json

# See DNS records for a specific domain (add --json for scripting)
steamer list-records example.com
steamer list-records example.com --json
```

#### DNS Record Management
```bash
# Add an A record (use "" for the root domain)
steamer add-a example.com home 192.0.2.1

# Add an AAAA record
steamer add-aaaa example.com home 2001:db8::1

# Add a CNAME in a flash
steamer add-cname example.com blog ghs.google.com

# Add an NS record for subdomain delegation (supports --json)
steamer add-ns example.com int ns-cloud-b1.googledomains.com
steamer add-ns example.com int ns-cloud-b2.googledomains.com --json

# Add a TXT record
steamer add-txt example.com _dmarc "v=DMARC1; p=none;"

# Remove a DNS record by its record ID (prompts for confirmation; pass --yes to skip)
steamer rm example.com 123456789
steamer rm example.com 123456789 --yes
```

#### Email Forwarding
```bash
# List email forwards for a domain (supports --json)
steamer email list example.com

# Add an email forwarding rule
steamer email add example.com hello user@example.org

# Remove an email forwarding rule (prompts for confirmation; pass --yes to skip)
steamer email rm example.com hello user@example.org
steamer email rm example.com hello user@example.org --yes
```

### 📋 Command Reference

| Command | Description | Flags |
| :--- | :--- | :--- |
| `steamer tui` | Launch interactive TUI for domains and records | `-d, --domain <domain>` |
| `steamer list-domains` | List all domains in your Porkbun account | `--json` |
| `steamer list-records <domain>` | List all DNS records for a domain | `--json` |
| `steamer list-tlds` | Browse TLD registration, renewal, and transfer pricing | `--force`, `--json` |
| `steamer search <query>` | Check domain availability and first-year pricing | `--tlds <list>`, `--json` |
| `steamer add-a <domain> <sub\|""> <ip>` | Create an IPv4 `A` record | `--json` |
| `steamer add-aaaa <domain> <sub\|""> <ip>` | Create an IPv6 `AAAA` record | `--json` |
| `steamer add-cname <domain> <sub> <target>` | Create a `CNAME` alias record | `--json` |
| `steamer add-ns <domain> <sub> <nameserver>` | Create an `NS` record to delegate a subdomain | `--json` |
| `steamer add-txt <domain> <sub\|""> <text>` | Create a `TXT` record | `--json` |
| `steamer rm <domain> <record-id>` | Delete a DNS record by ID | `-y, --yes`, `--json` |
| `steamer email list <domain>` | List configured email forwards and quota usage | `--json` |
| `steamer email add <domain> <addr> <to>` | Create an email forwarding rule | `--json` |
| `steamer email rm <domain> <addr> <to>` | Remove an email forwarding rule | `-y, --yes`, `--json` |
| `steamer version` | Display steamer version and build metadata | `--json` |

## 📚 Documentation
Check out `docs/api_reference.md` for a quick look at the Porkbun V3 API endpoints supported by this tool.

## 🐷 Why 'steamer'?
Because we're cooking with Pork(bun) and we're moving fast. Oink oink! 💨

---
*Built with ❤️ and Go and GeminiCLI.*

This is not an official Google product.
