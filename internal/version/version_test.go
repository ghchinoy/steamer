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

package version

import (
	"runtime/debug"
	"testing"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		name    string
		version string
		commit  string
		date    string
		bi      *debug.BuildInfo
		want    Info
	}{
		{
			name:    "ldflags_override_buildinfo",
			version: "0.2.0",
			commit:  "abc1234",
			date:    "2026-11-21T00:00:00Z",
			bi: &debug.BuildInfo{
				Main: debug.Module{Version: "v0.1.0"},
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "ffffffffffffffff"},
					{Key: "vcs.time", Value: "2026-01-01T00:00:00Z"},
				},
			},
			want: Info{Version: "0.2.0", Commit: "abc1234", Date: "2026-11-21T00:00:00Z"},
		},
		{
			name:    "fallback_to_buildinfo",
			version: "dev",
			commit:  "none",
			date:    "unknown",
			bi: &debug.BuildInfo{
				Main: debug.Module{Version: "v0.2.0"},
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "0123456789abcdef0123456789"},
					{Key: "vcs.time", Value: "2026-11-21T16:00:00Z"},
				},
			},
			want: Info{Version: "v0.2.0", Commit: "0123456789ab", Date: "2026-11-21T16:00:00Z"},
		},
		{
			name:    "devel_keeps_dev",
			version: "dev",
			commit:  "none",
			date:    "unknown",
			bi:      &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}},
			want:    Info{Version: "dev", Commit: "none", Date: "unknown"},
		},
		{
			name:    "nil_buildinfo",
			version: "dev",
			commit:  "none",
			date:    "unknown",
			bi:      nil,
			want:    Info{Version: "dev", Commit: "none", Date: "unknown"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := resolve(tc.version, tc.commit, tc.date, tc.bi)
			if got != tc.want {
				t.Errorf("resolve(%q, %q, %q, bi) = %+v, want %+v", tc.version, tc.commit, tc.date, got, tc.want)
			}
		})
	}
}
