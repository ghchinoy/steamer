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

// Package cmd provides the command-line interface for steamer.
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/theme"
	"github.com/ghchinoy/steamer/internal/version"
)

const (
	// GroupInfo is for commands that retrieve and display data.
	GroupInfo = "info"
	// GroupManagement is for commands that mutate DNS records.
	GroupManagement = "management"
	// GroupTUI is for interactive terminal interface commands.
	GroupTUI = "tui"
	// GroupEmail is for email forwarding commands.
	GroupEmail = "email"
)

type cliApp struct {
	v         *viper.Viper
	cfgFile   string
	newClient func() (*porkbun.Client, error)
}

// ExecuteContext runs the root CLI command with the given context.
func ExecuteContext(ctx context.Context) error {
	return NewRootCmd().ExecuteContext(ctx)
}

// NewRootCmd constructs a fresh steamer root command with all subcommands attached.
func NewRootCmd() *cobra.Command {
	app := &cliApp{
		v: viper.New(),
	}
	app.newClient = app.defaultClient

	rootCmd := &cobra.Command{
		Use:           "steamer",
		Short:         "A CLI for managing Porkbun domains",
		Long:          `steamer is a CLI tool built with Go, Cobra, and Bubble Tea to manage your Porkbun domains, DNS records, and email forwarding.`,
		Version:       version.Current().Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return app.initConfig(cmd)
		},
	}

	rootCmd.AddGroup(
		&cobra.Group{ID: GroupInfo, Title: "Information Commands:"},
		&cobra.Group{ID: GroupManagement, Title: "Record Management Commands:"},
		&cobra.Group{ID: GroupTUI, Title: "Interactive Commands:"},
		&cobra.Group{ID: GroupEmail, Title: "Email Forwarding Commands:"},
	)

	rootCmd.PersistentFlags().StringVar(&app.cfgFile, "config", "", "config file (default is $HOME/.config/steamer/config.yaml)")

	app.v.SetDefault("apikey", "")
	app.v.SetDefault("secretapikey", "")

	rootCmd.AddCommand(
		newAddACmd(app),
		newAddAAAACmd(app),
		newAddCNAMECmd(app),
		newAddNSCmd(app),
		newAddTXTCmd(app),
		newEmailCmd(app),
		newListDomainsCmd(app),
		newListRecordsCmd(app),
		newListTLDsCmd(app),
		newRmCmd(app),
		newSearchCmd(app),
		newTUICmd(app),
		newVersionCmd(),
	)

	return rootCmd
}

func (a *cliApp) initConfig(cmd *cobra.Command) error {
	if err := a.loadConfigFile(); err != nil {
		return err
	}

	if used := a.v.ConfigFileUsed(); used != "" {
		warnIfConfigPermsTooOpen(cmd, used)
	}

	a.v.SetEnvPrefix("PORKBUN")
	a.v.AutomaticEnv()
	a.applyFallbackEnv()

	// Only fall back to CWD .env if credentials are not already provided via
	// config file or shell environment variables.
	if a.resolveAPIKey() == "" || a.resolveSecretKey() == "" {
		if err := godotenv.Load(); err == nil {
			a.applyFallbackEnv()
		}
	}

	return nil
}

func (a *cliApp) loadConfigFile() error {
	if a.cfgFile != "" {
		a.v.SetConfigFile(a.cfgFile)
		if err := a.v.ReadInConfig(); err != nil {
			return fmt.Errorf("reading config file %q: %w", a.cfgFile, err)
		}
		return nil
	}

	home, err := os.UserHomeDir()
	if err == nil {
		// Priority 1: ~/.config/steamer (XDG standard, common on macOS for CLI)
		a.v.AddConfigPath(filepath.Join(home, ".config", "steamer"))
	}

	// Priority 2: ~/Library/Application Support/steamer (macOS standard)
	if configHome, cfgErr := os.UserConfigDir(); cfgErr == nil {
		a.v.AddConfigPath(filepath.Join(configHome, "steamer"))
	}

	if err == nil {
		// Priority 3: $HOME (legacy/simple)
		a.v.AddConfigPath(home)
	}

	a.v.SetConfigType("yaml")
	a.v.SetConfigName("config")
	err = a.v.ReadInConfig()
	if err == nil {
		return nil
	}
	if !isConfigNotFound(err) {
		return fmt.Errorf("reading config file: %w", err)
	}

	a.v.SetConfigName("steamer")
	if err := a.v.ReadInConfig(); err != nil && !isConfigNotFound(err) {
		return fmt.Errorf("reading config file: %w", err)
	}
	return nil
}

func isConfigNotFound(err error) bool {
	var notFound viper.ConfigFileNotFoundError
	return errors.As(err, &notFound)
}

func warnIfConfigPermsTooOpen(cmd *cobra.Command, path string) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	if info.Mode().Perm()&0o077 != 0 {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s\n", theme.Warn.Render(fmt.Sprintf(
			"Warning: config file %q has permissions %04o; consider running: chmod 600 %s",
			path, info.Mode().Perm(), path,
		)))
	}
}

func (a *cliApp) applyFallbackEnv() {
	if v := os.Getenv("API_KEY"); v != "" && a.v.GetString("apikey") == "" {
		a.v.Set("apikey", v)
	}
	if v := os.Getenv("API_SECRET"); v != "" && a.v.GetString("secretapikey") == "" {
		a.v.Set("secretapikey", v)
	}
}

func (a *cliApp) resolveAPIKey() string {
	for _, key := range []string{"apikey", "api_key"} {
		if val := a.v.GetString(key); val != "" {
			return val
		}
	}
	return ""
}

func (a *cliApp) resolveSecretKey() string {
	for _, key := range []string{"secretapikey", "api_secret", "apisecret", "secretkey"} {
		if val := a.v.GetString(key); val != "" {
			return val
		}
	}
	return ""
}

func (a *cliApp) loadClientConfig() (string, string, error) {
	apiKey := a.resolveAPIKey()
	secretKey := a.resolveSecretKey()

	if apiKey == "" || secretKey == "" {
		return "", "", errors.New("porkbun API key and secret must be provided via config file (~/.config/steamer/config.yaml), .env, or environment variables - hint: create a config file at ~/.config/steamer/config.yaml or set PORKBUN_APIKEY and PORKBUN_SECRETAPIKEY environment variables")
	}

	return apiKey, secretKey, nil
}

func (a *cliApp) defaultClient() (*porkbun.Client, error) {
	apiKey, secretKey, err := a.loadClientConfig()
	if err != nil {
		return nil, err
	}
	return porkbun.NewClient(apiKey, secretKey), nil
}
