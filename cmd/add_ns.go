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

func newAddNSCmd(app *cliApp) *cobra.Command {
	return newAddRecordCmd(app, recordCmdSpec{
		use:             "add-ns <domain> <subdomain> <nameserver>",
		short:           "Add a new NS record to a domain",
		long:            `Creates a new NS (Name Server) record for the specified domain and subdomain, delegating DNS authority for the subdomain to the target nameserver.`,
		recordType:      "NS",
		valueFieldJSON:  "nameserver",
		errorHintFormat: "hint: verify the domain exists in your Porkbun account and no conflicting CNAME record exists for %q",
		example: `  # Delegate int.example.com to a Google Cloud DNS nameserver
  steamer add-ns example.com int ns-cloud-b1.googledomains.com

  # Delegate and return output as JSON
  steamer add-ns example.com int ns-cloud-b1.googledomains.com --json`,
	})
}
