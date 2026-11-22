// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/theme"
)

type sslRetrieveResult struct {
	Domain      string               `json:"domain"`
	CertPath    string               `json:"cert_path"`
	KeyPath     string               `json:"key_path"`
	CertUpdated bool                 `json:"cert_updated"`
	KeyUpdated  bool                 `json:"key_updated"`
	Certificate *porkbun.SSLCertInfo `json:"certificate"`
}

type sslStatusResult struct {
	Domain      string               `json:"domain"`
	Certificate *porkbun.SSLCertInfo `json:"certificate"`
}

func newSSLCmd(app *cliApp) *cobra.Command {
	sslCmd := &cobra.Command{
		Use:     "ssl",
		Short:   "Inspect and retrieve managed Porkbun SSL certificates",
		GroupID: GroupSSL,
		Long:    `Inspect Porkbun's auto-renewed wildcard SSL bundle metadata or write fullchain/private key PEM files atomically for services with TLS hot-reload.`,
		Example: `  # Inspect certificate issuer, SANs, and expiry without touching disk
  steamer ssl status example.com

  # Retrieve and atomically write fullchain.pem (0644) and privkey.pem (0600)
  steamer ssl retrieve example.com --out-dir ./certs`,
	}

	sslCmd.AddCommand(
		newSSLStatusCmd(app),
		newSSLRetrieveCmd(app),
	)
	return sslCmd
}

func newSSLStatusCmd(app *cliApp) *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "status <domain>",
		Short: "Inspect managed SSL certificate issuer, SANs, and expiry for a domain",
		Long:  `Retrieves the Porkbun SSL bundle for a domain and parses the X.509 leaf certificate metadata without writing or printing the private key.`,
		Example: `  # Inspect SSL certificate status
  steamer ssl status example.com

  # Output certificate metadata as JSON
  steamer ssl status example.com --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := app.newClient()
			if err != nil {
				return err
			}

			domain := args[0]
			ctx, cancel := context.WithTimeout(cmd.Context(), defaultRequestTimeout)
			defer cancel()

			bundle, err := client.RetrieveSSLBundle(ctx, domain)
			if err != nil {
				return fmt.Errorf("retrieving ssl bundle for %q: %w", domain, err)
			}

			info, err := porkbun.ParseCertificateChain(bundle.CertificateChain, time.Now().UTC())
			if err != nil {
				return fmt.Errorf("parsing ssl certificate for %q: %w", domain, err)
			}

			if outputJSON {
				return printJSON(cmd.OutOrStdout(), sslStatusResult{
					Domain:      domain,
					Certificate: info,
				})
			}

			out := cmd.OutOrStdout()
			expiryStyle := theme.Pass
			switch {
			case info.DaysRemaining <= 14:
				expiryStyle = theme.Fail
			case info.DaysRemaining <= 30:
				expiryStyle = theme.Warn
			}

			_, _ = fmt.Fprintf(out, "%s %s\n", theme.Accent.Render("Domain:"), domain)
			_, _ = fmt.Fprintf(out, "%s %s\n", theme.Accent.Render("Subject:"), info.Subject)
			_, _ = fmt.Fprintf(out, "%s %s\n", theme.Accent.Render("Issuer:"), info.Issuer)
			_, _ = fmt.Fprintf(out, "%s %s\n", theme.Accent.Render("SANs:"), strings.Join(info.DNSNames, ", "))
			_, _ = fmt.Fprintf(out, "%s %s (%s)\n",
				theme.Accent.Render("Expires:"),
				info.NotAfter.Format(time.RFC3339),
				expiryStyle.Render(fmt.Sprintf("%d days remaining", info.DaysRemaining)),
			)
			_, _ = fmt.Fprintf(out, "%s %d certs\n", theme.Muted.Render("Chain:"), info.ChainCount)
			_, err = fmt.Fprintf(out, "%s %s\n", theme.Muted.Render("SHA-256:"), info.SHA256Fingerprint)
			return err
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output certificate metadata in JSON format")
	return cmd
}

func newSSLRetrieveCmd(app *cliApp) *cobra.Command {
	var (
		outDir     string
		certFile   string
		keyFile    string
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "retrieve <domain>",
		Short: "Download and atomically save Porkbun SSL fullchain and private key PEM files",
		Long:  `Retrieves the Porkbun managed SSL bundle for a domain and atomically writes the certificate chain (0644) and private key (0600). Existing files with identical bytes are left untouched so file-watching servers only reload when the certificate actually changes.`,
		Example: `  # Save fullchain.pem and privkey.pem into a directory
  steamer ssl retrieve example.com --out-dir ./certs

  # Save with custom filenames and JSON output for automation
  steamer ssl retrieve example.com --out-dir /etc/tls/certs --cert-file cert.pem --key-file key.pem --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			certName, err := cleanFilename("cert-file", certFile)
			if err != nil {
				return err
			}
			keyName, err := cleanFilename("key-file", keyFile)
			if err != nil {
				return err
			}

			client, err := app.newClient()
			if err != nil {
				return err
			}

			domain := args[0]
			ctx, cancel := context.WithTimeout(cmd.Context(), defaultRequestTimeout)
			defer cancel()

			bundle, err := client.RetrieveSSLBundle(ctx, domain)
			if err != nil {
				return fmt.Errorf("retrieving ssl bundle for %q: %w", domain, err)
			}

			info, err := porkbun.ParseCertificateChain(bundle.CertificateChain, time.Now().UTC())
			if err != nil {
				return fmt.Errorf("parsing ssl certificate for %q: %w", domain, err)
			}

			targetDir := strings.TrimSpace(outDir)
			if targetDir == "" {
				targetDir = "."
			}
			if err := os.MkdirAll(targetDir, 0o700); err != nil {
				return fmt.Errorf("creating output directory %q: %w", targetDir, err)
			}

			root, err := os.OpenRoot(targetDir)
			if err != nil {
				return fmt.Errorf("opening output directory %q: %w", targetDir, err)
			}
			defer func() { _ = root.Close() }()

			certPath := filepath.Join(targetDir, certName)
			keyPath := filepath.Join(targetDir, keyName)

			certUpdated, err := writeFileIfChangedAtomic(root, certName, []byte(bundle.CertificateChain), 0o644)
			if err != nil {
				return fmt.Errorf("writing certificate %q: %w", certPath, err)
			}
			keyUpdated, err := writeFileIfChangedAtomic(root, keyName, []byte(bundle.PrivateKey), 0o600)
			if err != nil {
				return fmt.Errorf("writing private key %q: %w", keyPath, err)
			}

			res := sslRetrieveResult{
				Domain:      domain,
				CertPath:    certPath,
				KeyPath:     keyPath,
				CertUpdated: certUpdated,
				KeyUpdated:  keyUpdated,
				Certificate: info,
			}
			if outputJSON {
				return printJSON(cmd.OutOrStdout(), res)
			}

			statusText := "unchanged (already up to date)"
			if certUpdated || keyUpdated {
				statusText = "updated"
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\n  Cert: %s\n  Key:  %s\n  SANs: %s (expires %s, %d days remaining)\n",
				theme.Pass.Render(fmt.Sprintf("SSL bundle for %s %s", domain, statusText)),
				certPath,
				keyPath,
				strings.Join(info.DNSNames, ", "),
				info.NotAfter.Format(time.RFC3339),
				info.DaysRemaining,
			)
			return err
		},
	}

	cmd.Flags().StringVar(&outDir, "out-dir", ".", "Directory to write certificate and private key PEM files")
	cmd.Flags().StringVar(&certFile, "cert-file", "fullchain.pem", "Filename for the PEM certificate chain")
	cmd.Flags().StringVar(&keyFile, "key-file", "privkey.pem", "Filename for the PEM private key (written 0600)")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output result in JSON format (private key is never printed)")
	return cmd
}

func cleanFilename(flagName, raw string) (string, error) {
	base := filepath.Base(strings.TrimSpace(raw))
	if base == "" || base == "." || base == ".." || base == string(filepath.Separator) {
		return "", fmt.Errorf("invalid --%s %q: must be a filename", flagName, raw)
	}
	return base, nil
}

func writeFileIfChangedAtomic(root *os.Root, name string, data []byte, perm os.FileMode) (bool, error) {
	if existing, err := root.ReadFile(name); err == nil && bytes.Equal(existing, data) {
		_ = root.Chmod(name, perm)
		return false, nil
	}
	tmp := name + ".tmp"
	if err := root.WriteFile(tmp, data, perm); err != nil {
		return false, err
	}
	if err := root.Chmod(tmp, perm); err != nil {
		_ = root.Remove(tmp)
		return false, err
	}
	if err := root.Rename(tmp, name); err != nil {
		_ = root.Remove(tmp)
		return false, err
	}
	return true, nil
}
