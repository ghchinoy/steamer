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

// Package version holds build-time version metadata for steamer.
package version

import "runtime/debug"

// Build metadata injected at release time via -ldflags.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// Info represents resolved version metadata for the CLI binary.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// Current returns build metadata, falling back to runtime/debug.ReadBuildInfo
// when binaries are built with "go install".
func Current() Info {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return Info{Version: Version, Commit: Commit, Date: Date}
	}
	return resolve(Version, Commit, Date, bi)
}

func resolve(ver, commit, date string, bi *debug.BuildInfo) Info {
	info := Info{
		Version: ver,
		Commit:  commit,
		Date:    date,
	}
	if bi == nil {
		return info
	}

	if info.Version == "dev" && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		info.Version = bi.Main.Version
	}

	for _, setting := range bi.Settings {
		switch setting.Key {
		case "vcs.revision":
			if info.Commit == "none" && setting.Value != "" {
				rev := setting.Value
				if len(rev) > 12 {
					rev = rev[:12]
				}
				info.Commit = rev
			}
		case "vcs.time":
			if info.Date == "unknown" && setting.Value != "" {
				info.Date = setting.Value
			}
		}
	}

	return info
}
