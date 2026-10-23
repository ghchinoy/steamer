# 🐷 steamer

**steamer** is a high-speed, interactive CLI and TUI for managing your Porkbun domains. Built with Go, Cobra, and Charm's Bubble Tea, it's designed to make DNS management feel less like a chore and more like a breeze.

Whether you're delegating a subdomain with `NS` records, adding a `CNAME` for a new project, configuring email forwarding, or browsing your portfolio, **steamer** has you covered.

## ✨ Features

- **🚀 Quick DNS Management:** Create `A`, `AAAA`, `CNAME`, `NS` (subdomain delegation), and `TXT` records, or delete records (`rm`) in a single line.
- **🔍 Domain Search & TLD Pricing:** Check domain availability across multiple TLDs (`search`) and browse cached TLD pricing (`list-tlds`).
- **📧 Email Forwarding:** List, add, and remove Porkbun email forwarding rules (`email list`, `email add`, `email rm`).
- **🤖 Agent-Friendly Output:** Supports `--json` flags for deterministic scripting and AI agent workflows.
- **🎮 Interactive TUI:** A terminal interface to browse your domains, DNS records, and TLD pricing tables.
- **🔐 Secure Config:** Supports XDG-style configuration, `.env` files, and environment variables.
- **📖 API Reference:** Includes a local copy of the Porkbun v3 API documentation.

## 🛠️ Installation

The most idiomatic way to install **steamer** is via `go install`:

```bash
go install github.com/ghchinoy/steamer@latest
```

This will build the binary and place it in your `$GOPATH/bin` (typically `~/go/bin`). Ensure that directory is in your system's `PATH`.

Alternatively, if you've cloned the repository, you can install it locally:

```bash
go install .
```

## ⚙️ Configuration

**steamer** looks for your Porkbun API credentials in several places. You'll need an API Key and Secret from [Porkbun](https://porkbun.com/account/api).

### 1. The Pro Way (Config File)
Create a config file at `~/.config/steamer/config.yaml`:

```yaml
apikey: pk1_your_api_key
secretapikey: sk1_your_secret_key
```

### 2. The Dev Way (.env)
Create a `.env` file in your project root:

```env
API_KEY=pk1_your_api_key
API_SECRET=sk1_your_secret_key
```

### 3. The Quick Way (Environment Variables)
```bash
export PORKBUN_APIKEY=pk1_...
export PORKBUN_SECRETAPIKEY=sk1_...
```

## 🚀 Usage

### The Terminal UI (TUI)
Just run `steamer tui` and enjoy the ride. Use `j`/`k` to navigate and `enter` to dive into records.

```bash
# Start the TUI
steamer tui

# Jump straight to a domain
steamer tui -d aaie.cloud
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
steamer list-records aaie.cloud
steamer list-records aaie.cloud --json
```

#### DNS Record Management
```bash
# Add an A record (use "" for the root domain)
steamer add-a aaie.cloud home 1.2.3.4

# Add an AAAA record
steamer add-aaaa aaie.cloud home 2001:db8::1

# Add a CNAME in a flash
steamer add-cname aaie.cloud blog ghs.google.com

# Add an NS record for subdomain delegation (supports --json)
steamer add-ns aaie.cloud int ns-cloud-b1.googledomains.com
steamer add-ns aaie.cloud int ns-cloud-b2.googledomains.com --json

# Add a TXT record
steamer add-txt aaie.cloud _dmarc "v=DMARC1; p=none;"

# Remove a DNS record by its record ID
steamer rm aaie.cloud 123456789
```

#### Email Forwarding
```bash
# List email forwards for a domain (supports --json)
steamer email list chinoy.si

# Add an email forwarding rule
steamer email add chinoy.si hello ghchinoy@gmail.com

# Remove an email forwarding rule
steamer email rm chinoy.si hello ghchinoy@gmail.com
```

### 📋 Command Reference

| Command | Description | Flags |
| :--- | :--- | :--- |
| `steamer tui` | Launch interactive TUI for domains and records | `-d, --domain <domain>` |
| `steamer list-domains` | List all domains in your Porkbun account | `--json` |
| `steamer list-records <domain>` | List all DNS records for a domain | `--json` |
| `steamer list-tlds` | Browse TLD registration, renewal, and transfer pricing | `--force`, `--json` |
| `steamer search <query>` | Check domain availability and first-year pricing | `--tlds <list>`, `--json` |
| `steamer add-a <domain> <sub\|""> <ip>` | Create an IPv4 `A` record | |
| `steamer add-aaaa <domain> <sub\|""> <ip>` | Create an IPv6 `AAAA` record | |
| `steamer add-cname <domain> <sub> <target>` | Create a `CNAME` alias record | |
| `steamer add-ns <domain> <sub> <nameserver>` | Create an `NS` record to delegate a subdomain | `--json` |
| `steamer add-txt <domain> <sub\|""> <text>` | Create a `TXT` record | |
| `steamer rm <domain> <record-id>` | Delete a DNS record by ID | |
| `steamer email list <domain>` | List configured email forwards and quota usage | `--json` |
| `steamer email add <domain> <addr> <to>` | Create an email forwarding rule | |
| `steamer email rm <domain> <addr> <to>` | Remove an email forwarding rule | |

## 📚 Documentation
Check out `docs/api_reference.md` for a quick look at the Porkbun V3 API endpoints supported by this tool.

## 🐷 Why 'steamer'?
Because we're cooking with Pork(bun) and we're moving fast. Oink oink! 💨

---
*Built with ❤️ and Go and GeminiCLI.*

This is not an official Google product.