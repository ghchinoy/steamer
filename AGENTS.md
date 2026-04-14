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

### 🌐 Porkbun API DNS Limitations

- **Record Limits:** A domain can have a maximum of 200 DNS records across all types.
- **Root Record Conflicts:** The API will fail with a generic "We were unable to create the DNS record" if you try to add an `A` or `AAAA` record to the root (subdomain `""`) when a root `ALIAS` or `CNAME` record already exists. Ensure conflicting root alias records are deleted before adding IP-based root records.

### 🤖 CLI Design Best Practices (A2A Standard)

To ensure the CLI remains human-friendly and highly accessible to AI agents, `steamer` follows these best practices (adapted from the A2A CLI standard):

- **Structured Discoverability:** Use `GroupID` in Cobra commands to categorize them in the `--help` output (e.g., `GroupInfo`, `GroupManagement`, `GroupTUI`).
- **Three Pillars of Documentation:** Every command MUST have a `Short` description (action verb), a `Long` detailed explanation, and copy-pasteable `Example` blocks.
- **Agent-First Interoperability:** Data-returning commands (like `list-domains`) MUST support a `--json` flag for deterministic, non-interactive parsing.
- **Proactive Error Guidance:** Don't just fail on missing config; print a clear "Hint:" explaining exactly how to fix it (e.g., setting environment variables or creating a config file).
- **Semantic Color Usage (Lipgloss):** Use the centralized `internal/theme` package rather than raw colors.
  - `Accent`: Headers, landmarks (`#399ee6` / `#59c2ff`).
  - `Pass`: Success states (`#86b300` / `#c2d94c`).
  - `Warn`: Warnings, pending (`#f2ae49` / `#ffb454`).
  - `Fail`: Errors (`#f07171` / `#f07178`).
  - `Muted`: De-emphasis, types (`#828c99` / `#6c7680`).
  - `ID`: Identifiers (`#46ba94` / `#95e6cb`).





<!-- BEGIN BEADS INTEGRATION v:1 profile:minimal hash:ca08a54f -->
## Beads Issue Tracker

This project uses **bd (beads)** for issue tracking. Run `bd prime` to see full workflow context and commands.

### Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --claim  # Claim work
bd close <id>         # Complete work
```

### Rules

- Use `bd` for ALL task tracking — do NOT use TodoWrite, TaskCreate, or markdown TODO lists
- Run `bd prime` for detailed command reference and session close protocol
- Use `bd remember` for persistent knowledge — do NOT use MEMORY.md files

## Session Completion

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd dolt push
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
<!-- END BEADS INTEGRATION -->
