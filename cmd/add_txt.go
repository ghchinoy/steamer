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

func newAddTXTCmd(app *cliApp) *cobra.Command {
	return newAddRecordCmd(app, recordCmdSpec{
		use:         "add-txt <domain> <subdomain> <text>",
		short:       "Add a new TXT record to a domain",
		long:        `Creates a new TXT record for the specified domain and subdomain. TXT records are commonly used for domain verification, SPF, and DMARC.`,
		recordType:  "TXT",
		successVerb: "with value",
		example: `  # Add a DMARC record
  steamer add-txt example.com _dmarc "v=DMARC1; p=none;"

  # Add a site verification record to the root domain
  steamer add-txt example.com "" "google-site-verification=abc123xyz"`,
	})
}
