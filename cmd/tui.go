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

	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var domainFlag string

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start the interactive TUI",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client := porkbun.NewClient(apiKey, secretKey)

		p := tea.NewProgram(tui.NewModel(client, domainFlag))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	tuiCmd.Flags().StringVarP(&domainFlag, "domain", "d", "", "Start with records for a specific domain")
	rootCmd.AddCommand(tuiCmd)
}
