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

	"github.com/spf13/cobra"
)

var addACmd = &cobra.Command{
	Use:   "add-a [domain] [subdomain] [ip]",
	Short: "Add a new A record to a domain",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		domain := args[0]
		subdomain := args[1]
		ip := args[2]

		client := porkbun.NewClient(apiKey, secretKey)
		id, err := client.CreateRecord(domain, porkbun.CreateRecordRequest{
			Name:    subdomain,
			Type:    "A",
			Content: ip,
		})
		if err != nil {
			fmt.Printf("Error creating A record: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created A record for %s.%s pointing to %s (ID: %s)\n", subdomain, domain, ip, id)
	},
}

func init() {
	rootCmd.AddCommand(addACmd)
}
