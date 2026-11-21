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

import "github.com/spf13/cobra"

func newAddACmd(app *cliApp) *cobra.Command {
	return newAddRecordCmd(app, recordCmdSpec{
		use:        "add-a <domain> <subdomain> <ip>",
		short:      "Add a new A record to a domain",
		long:       `Creates a new IPv4 'A' record for the specified domain and subdomain. Use an empty string "" for the root domain.`,
		recordType: "A",
		example: `  # Add an A record for www.example.com
  steamer add-a example.com www 192.0.2.1

  # Add an A record for the root domain (example.com)
  steamer add-a example.com "" 192.0.2.1`,
	})
}
