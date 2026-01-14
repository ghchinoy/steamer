## Issue Tracking

This project uses **bd (beads)** for issue tracking.
Run `bd prime` for workflow context, or install hooks (`bd hooks install`) for auto-injection.

**Quick reference:**
- `bd ready` - Find unblocked work
- `bd create "Title" --type task --priority 2` - Create issue
- `bd close <id>` - Complete work
- `bd sync` - Sync with git (run at session end)

For full workflow details: `bd prime`



## Lessons Learned & Patterns



### 🧩 Handling Inconsistent API Types

The Porkbun API (and potentially others) may return fields like `autoRenew` or `notLocal` as either a `string` ("0") or a `number` (0) depending on the account state or endpoint.

- **Pattern:** Use `interface{}` in Go structs for these fields to prevent unmarshaling errors.

- **Example:**

  ```go

  AutoRenew interface{} `json:"autoRenew"`

  ```



### ⚙️ Resilient Configuration (Viper + XDG)

- **Flexible Keys:** Support multiple naming conventions (e.g., `apikey`, `api_key`, `API_KEY`) in `getClientConfig` to accommodate variations between `.env` files and YAML configs.

- **XDG Paths:** Prioritize `~/.config/steamer/config.yaml` for XDG compliance, but maintain search paths for `~/.steamer.yaml` (legacy) and `~/Library/Application Support/steamer` (macOS native) to ensure a smooth user experience across platforms.

- **Format Awareness:** Be prepared for users to accidentally use `.env` syntax (`KEY=VALUE`) in a `.yaml` file. While not valid YAML, clear error messages or flexible loaders help.



### 🎨 TUI Alignment



When using `fmt.Sprintf` for TUI tables (e.g., in Bubble Tea), ensure format strings like `%-10s` are clean of accidental escape characters (like `%-\10s`) which cause the `fmt` package to output raw type information instead of the formatted string.







### 📝 Documentation Strategy: Binary to Markdown



- **Pattern:** Proactively extract full API documentation into a single, structured Markdown file (e.g., `docs/api_reference.md`).



- **Benefit:** This makes the documentation "agent-friendly" (searchable and parseable by AI) and allows for a "binary-free" repository.







### 🐹 Go Idiomatics: Installation



- **Pattern:** For Go-based tools, prioritize `go install github.com/user/repo@latest` instructions. It's the standard way to distribute Go binaries and ensures they land in the user's `$GOPATH/bin`.







### ⚠️ Configuration Syntax Pitfalls



- **Insight:** Users often accidentally use `.env` syntax (`KEY=VALUE`) in YAML files.



- **Pattern:** When config loading fails or returns empty values, include a "Syntax Check" hint in error messages to remind the user to use YAML colons (`key: value`).




