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
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/ghchinoy/steamer/internal/tui"
)

func newTUICmd(app *cliApp) *cobra.Command {
	var domainFlag string

	cmd := &cobra.Command{
		Use:     "tui",
		Short:   "Start the interactive TUI",
		GroupID: GroupTUI,
		Long:    `Starts a rich, interactive terminal UI built with Bubble Tea. Use the TUI to visually browse your domains and inspect their DNS records without needing to run separate commands.`,
		Example: `  # Start the default TUI
  steamer tui

  # Start the TUI directly focused on a specific domain
  steamer tui -d example.com`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := app.newClient()
			if err != nil {
				return err
			}

			parentCtx := cmd.Context()
			newReqCtx := func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(parentCtx, defaultRequestTimeout)
			}

			p := tea.NewProgram(tui.NewModel(client, domainFlag, newReqCtx), tea.WithContext(parentCtx))
			if _, err := p.Run(); err != nil {
				return fmt.Errorf("running TUI: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&domainFlag, "domain", "d", "", "Start with records for a specific domain")
	return cmd
}
