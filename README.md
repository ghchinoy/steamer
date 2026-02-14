# 🐷 steamer

**steamer** is a high-speed, interactive CLI and TUI for managing your Porkbun domains. Built with Go, Cobra, and Charm's Bubble Tea, it's designed to make DNS management feel less like a chore and more like a breeze.

Whether you're adding a CNAME for a new project or just browsing your portfolio, **steamer** has you covered.

## ✨ Features

- **🚀 Quick Commands:** Add CNAMEs or list records in a single line.
- **🎮 Interactive TUI:** A beautiful terminal interface to browse your domains and records.
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
```bash
# List all your domains
steamer list-domains

# See records for a specific domain
steamer list-records aaie.cloud

# Add an A record
steamer add-a aaie.cloud home 1.2.3.4

# Add a CNAME in a flash
steamer add-cname aaie.cloud blog ghs.google.com
```

## 📚 Documentation
Check out `docs/api_reference.md` for a quick look at the Porkbun V3 API endpoints supported by this tool.

## 🐷 Why 'steamer'?
Because we're cooking with Pork(bun) and we're moving fast. Oink oink! 💨

---
*Built with ❤️ and Go and GeminiCLI.*

This is not an official Google product.