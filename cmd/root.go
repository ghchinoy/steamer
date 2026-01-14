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
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "steamer",
	Short: "A CLI for managing Porkbun domains",
	Long:  `steamer is a CLI tool built with Go, Cobra, and Bubble Tea to manage your Porkbun domains and DNS records.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/steamer/config.yaml)")

	viper.SetDefault("apikey", "")
	viper.SetDefault("secretapikey", "")
}

func initConfig() {
	// 1. Try loading .env file
	_ = godotenv.Load()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// 2. Setup XDG-style config path
		home, _ := os.UserHomeDir()
		
		// Priority 1: ~/.config/steamer (XDG standard, common on macOS for CLI)
		viper.AddConfigPath(filepath.Join(home, ".config", "steamer"))

		// Priority 2: ~/Library/Application Support/steamer (macOS standard)
		configHome, err := os.UserConfigDir()
		if err == nil {
			viper.AddConfigPath(filepath.Join(configHome, "steamer"))
		}
		
		// Priority 3: $HOME (Legacy/simple)
		viper.AddConfigPath(home)

		viper.SetConfigType("yaml")
		
		// Try 'config' first (idiomatic), then 'steamer' (fallback)
		viper.SetConfigName("config")
		if err := viper.ReadInConfig(); err != nil {
			viper.SetConfigName("steamer")
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("PORKBUN")

	// Map .env style names to our viper keys if they aren't already set
	if v := os.Getenv("API_KEY"); v != "" && viper.GetString("apikey") == "" {
		viper.Set("apikey", v)
	}
	if v := os.Getenv("API_SECRET"); v != "" && viper.GetString("secretapikey") == "" {
		viper.Set("secretapikey", v)
	}

	// Final attempt to read if it wasn't 'config' and we switched to 'steamer'
	_ = viper.ReadInConfig()
}

func getClientConfig() (string, string, error) {
	apiKey := viper.GetString("apikey")
	if apiKey == "" {
		apiKey = viper.GetString("api_key")
	}

	secretKey := viper.GetString("secretapikey")
	if secretKey == "" {
		secretKey = viper.GetString("api_secret")
	}
	if secretKey == "" {
		secretKey = viper.GetString("apisecret")
	}
	if secretKey == "" {
		secretKey = viper.GetString("secretkey")
	}

	if apiKey == "" || secretKey == "" {
		return "", "", fmt.Errorf("Porkbun API Key and Secret must be provided via config file (~/.config/steamer/config.yaml), .env, or environment variables")
	}

	return apiKey, secretKey, nil
}
