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

	"github.com/spf13/cobra"

	"github.com/ghchinoy/steamer/internal/theme"
	"github.com/ghchinoy/steamer/internal/version"
)

func newVersionCmd() *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print steamer version and build metadata",
		GroupID: GroupInfo,
		Long:    `Displays the steamer release version, git commit SHA, and build date. Supports --json for scripting and automation.`,
		Example: `  # Print version details
  steamer version

  # Print version details as JSON
  steamer version --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			info := version.Current()
			if outputJSON {
				return printJSON(cmd.OutOrStdout(), info)
			}

			_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s (%s %s, %s %s)\n",
				theme.Accent.Render("steamer"),
				info.Version,
				theme.Muted.Render("commit:"),
				info.Commit,
				theme.Muted.Render("built:"),
				info.Date,
			)
			return err
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output version metadata in JSON format")
	return cmd
}
