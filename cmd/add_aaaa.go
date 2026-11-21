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

func newAddAAAACmd(app *cliApp) *cobra.Command {
	return newAddRecordCmd(app, recordCmdSpec{
		use:        "add-aaaa <domain> <subdomain> <ip>",
		short:      "Add a new AAAA record to a domain",
		long:       `Creates a new IPv6 'AAAA' record for the specified domain and subdomain. Use an empty string "" for the root domain.`,
		recordType: "AAAA",
		example: `  # Add an AAAA record for ipv6.example.com
  steamer add-aaaa example.com ipv6 2001:db8::1

  # Add an AAAA record for the root domain (example.com)
  steamer add-aaaa example.com "" 2001:db8::1`,
	})
}
