# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and [Common Changelog](https://common-changelog.org/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-11-21

### Added

- DNS record management commands: `add-a`, `add-aaaa`, `add-txt`, `add-ns` (subdomain delegation), and `rm`.
- Domain availability search (`search`) with single-domain and multi-TLD phrase lookup plus Porkbun rate-limit pacing.
- TLD pricing browser (`list-tlds`) rendered as an interactive Bubble Tea table with 7-day local caching and `--force` refresh `[steamer-zrk]`.
- Email forwarding management commands: `email list`, `email add`, and `email rm`.
- Binary version reporting via `steamer --version` and `steamer version [--json]` using release `-ldflags` with `runtime/debug.ReadBuildInfo` fallback `[steamer-098]`.
- Deterministic `--json` output across query and mutation commands (`list-domains`, `list-records`, `list-tlds`, `search`, `add-*`, `rm`, `email *`, `version`).
- Cross-platform release automation via GoReleaser (GitHub Releases, checksums, and `ghchinoy/homebrew-tap` cask publishing) and GitHub Actions CI (`go test -race`, `golangci-lint`, `govulncheck`) `[steamer-rg4]`, `[steamer-098]`.

### Changed

- **Breaking:** `rm` and `email rm` now prompt for interactive confirmation before deleting resources and refuse non-interactive execution unless `-y`/`--yes` is passed `[steamer-d1s]`.
- Improved CLI discoverability with grouped `--help` categories, semantic adaptive Lipgloss styling via `internal/theme`, and `RunE`-based stderr error reporting `[steamer-x4j]`.
- Upgraded the Go module floor to `go 1.26.10` and updated transitive dependencies (`golang.org/x/text` to `v0.44.0`, `golang.org/x/sys` to `v0.50.0`) `[steamer-f42]`.

### Fixed

- Handled inconsistent Porkbun API JSON types (numeric vs. string values) for DNS record IDs, `securityLock`, `whoisPrivacy`, `autoRenew`, `notLocal`, and `Label.ID` using `FlexibleString` and `json.Decoder.UseNumber()`, preventing scientific notation formatting on large IDs `[steamer-63q]`, `[steamer-9p8]`.
- Fixed `search --json` serializing Go errors as empty `{}` objects instead of error message strings `[steamer-x4j]`.
- Fixed `domain/checkDomain` response unmarshaling, `email list` ANSI column alignment, root-domain leading-dot formatting in `add-*` output, double config loading, and missing `status == "SUCCESS"` response checks `[steamer-x4j]`.

### Security

- Validated and URL-escaped user-supplied domain and record ID path segments before issuing authenticated Porkbun API POST requests `[steamer-d1s]`.
- Hardened the Porkbun HTTP client with a default 30-second timeout, `context.Context` cancellation, a 10 MiB response body cap, and truncated non-JSON error output `[steamer-d1s]`.
- Scoped TLD cache reads and writes with `os.OpenRoot`, restricted current-directory `.env` loading to fallback-only when credentials are absent, added permission warnings for group/world-readable config files, and replaced personal email addresses and domains with RFC 2606 placeholders `[steamer-d1s]`.
- Resolved `govulncheck` package and module advisories (`GO-2026-5970`, `GO-2026-6629`, `GO-2026-5024`) `[steamer-f42]`.

## [0.1.0] - 2026-02-08

_Initial release._

### Added

- Interactive Bubble Tea TUI (`tui`) and core CLI commands (`list-domains`, `list-records`, `add-cname`) backed by the Porkbun API v3 client `[steamer-gqh]`, `[steamer-k9c]`, `[steamer-3ar]`, `[steamer-2i6]`, `[steamer-xq7]`.
- Multi-source credential configuration supporting XDG (`~/.config/steamer/config.yaml`), macOS Application Support, legacy `~/.steamer.yaml`, `.env`, and `PORKBUN_*` environment variables `[steamer-0w6]`, `[steamer-c9x]`.
- Extracted Markdown reference for the Porkbun API v3 (`docs/api_reference.md`) `[steamer-asg]`, `[steamer-4fj]`.

### Fixed

- Aligned `go.mod` module path to `github.com/ghchinoy/steamer` for `go install` compatibility and resolved initial `golangci-lint` issues `[steamer-m6f]`, `[steamer-rcu]`, `[steamer-7ft]`, `[steamer-eew]`, `[steamer-zod]`.

[Unreleased]: https://github.com/ghchinoy/steamer/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/ghchinoy/steamer/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/ghchinoy/steamer/releases/tag/v0.1.0
