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

// Package theme provides standard semantic colors for the CLI and TUI.
package theme

import "github.com/charmbracelet/lipgloss"

var (
	// Accent is used for landmarks, headers, group titles, and section labels.
	Accent = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#399ee6", Dark: "#59c2ff"})

	// Command is used for scan targets like command names and flags.
	Command = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#5c6166", Dark: "#bfbdb6"})

	// Pass is used for success states and completed tasks.
	Pass = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#86b300", Dark: "#c2d94c"})

	// Warn is used for active tasks, warnings, and pending states.
	Warn = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#f2ae49", Dark: "#ffb454"})

	// Fail is used for errors, failed tasks, and rejected states.
	Fail = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#f07171", Dark: "#f07178"})

	// Muted is used to de-emphasize metadata, types, defaults, and previews.
	Muted = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#828c99", Dark: "#6c7680"})

	// ID is used for identifiers like TaskIDs or unique Record IDs.
	ID = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#46ba94", Dark: "#95e6cb"})
)
