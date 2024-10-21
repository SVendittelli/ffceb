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
	"context"
	"database/sql"
	"os"
	"sort"
	"strings"

	"github.com/SVendittelli/ffceb/repository"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		ctx := context.Background()

		log.Info("exporting domains")

		db_file := viper.GetString("profile") + "/permissions.sqlite"
		log.Debug("accessing db", "db_file", db_file)

		// Connect to the Firefox permissions database
		db, err := sql.Open("sqlite3", db_file)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		queries := repository.New(db)

		origins, err := queries.ListExcludedOrigins(ctx)
		if err != nil {
			log.Fatal(err)
		}

		var domains []string
		for _, origin := range origins {
			log.Debug("exporting", "origin", origin.String)

			// Strip protocol from origin
			split := strings.Split(origin.String, "://")
			domains = append(domains, split[1])
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

		log.Infof("exported %d domains", len(domains))
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
	log.Debug("de-duplicating", "count", len(s))
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range s {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	log.Debug("done de-duplicating", "count", len(list))
	return list
}
