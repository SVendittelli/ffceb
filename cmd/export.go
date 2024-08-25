/*
Copyright © 2024 Sam Vendittelli <sam.vendittelli@hotmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("export called")

		// Connect to the FireFox permissions database
		db, err := sql.Open("sqlite3", `filepath`)
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		rows, err := db.Query("SELECT origin FROM moz_perms WHERE type = 'cookie' AND permission = 1 AND expireTime = 0")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var domains []string
		for rows.Next() {
			var origin string

			err = rows.Scan(&origin)
			if err != nil {
				log.Fatal(err)
			}

			// Strip protocol from origin
			split := strings.Split(origin, "://")
			domains = append(domains, split[1])
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		// Remove duplicates and sort
		domains = removeDuplicates(domains)
		sort.Strings(domains)

		// Write domains to file
		f, err := os.Create("export.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		for _, domain := range domains {
			_, err := f.WriteString(domain + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// removeDuplicates removes duplicates from a slice of strings
func removeDuplicates(s []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range s {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
