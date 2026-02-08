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



The Porkbun API (and potentially others) may return fields like `autoRenew`, `notLocal`, or even record `id`s as either a `string` ("123") or a `number` (123) depending on the account state or endpoint.



- **Pattern:** Use `interface{}` in Go structs for these fields to prevent unmarshaling errors.



- **Example:**



  ```go



  ID interface{} `json:"id"`



  ```



- **Formatting:** Use `%v` in `fmt.Sprintf` or `fmt.Printf` to safely handle `interface{}` values that might contain strings or integers.







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







### 📦 Go Module Paths and `go install`

- **Pattern:** For tools intended to be installed via `go install github.com/user/repo@latest`, the `module` path in `go.mod` MUST match the repository path.
- **Problem:** If `go.mod` declares `module steamer` but is hosted at `github.com/ghchinoy/steamer`, `go install` will fail with a version constraints conflict.
- **Fix:** Ensure `module github.com/ghchinoy/steamer` is used in `go.mod` and all internal imports are updated accordingly.

### 🧹 Linting and Style

- **Tooling:** Use `golangci-lint` with a comprehensive set of linters (including `revive`, `govet`, `staticcheck`, `errcheck`) to maintain high code quality.
- **Documentation:** All exported symbols (types, functions, methods) and packages must have documentation comments to comply with Google Go style.
- **Error Handling (Defer):** To satisfy `errcheck` for `resp.Body.Close()` without complex error handling in a defer, use an anonymous function: `defer func() { _ = resp.Body.Close() }()`.
- **Error Strings:** In Go, error strings should be lowercase and not end with punctuation (e.g., `fmt.Errorf("api error: %s", msg)`).




