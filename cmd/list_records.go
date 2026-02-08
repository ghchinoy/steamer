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

var listRecordsCmd = &cobra.Command{
	Use:   "list-records [domain]",
	Short: "List DNS records for a specific domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		domain := args[0]
		client := porkbun.NewClient(apiKey, secretKey)
		records, err := client.RetrieveRecords(domain)
		if err != nil {
			fmt.Printf("Error retrieving records for %s: %v\n", domain, err)
			os.Exit(1)
		}

		fmt.Printf("%-10s %-25s %-10s %-30s\n", "ID", "NAME", "TYPE", "CONTENT")
		for _, r := range records {
			fmt.Printf("%-10v %-25s %-10s %-30s\n", r.ID, r.Name, r.Type, r.Content)
		}
	},
}

func init() {
	rootCmd.AddCommand(listRecordsCmd)
}
