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
	"bufio"
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/SVendittelli/ffceb/repository"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var overwrite bool

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		log.Info("importing domains")
		log.Debug("flags", "overwrite", overwrite)

		f, err := os.Open("export.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		var domains []string

		log.Debug("reading domains from file")
		for scanner.Scan() {
			log.Debug("reading", "domain", scanner.Text())
			domains = append(domains, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		db_file := viper.GetString("profile") + "/permissions.sqlite"
		log.Debug("accessing db", "db_file", db_file)

		// Connect to the Firefox permissions database
		db, err := sql.Open("sqlite3", db_file)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		queries := repository.New(db)

		log.Debug("starting transaction")
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback()

		qtx := queries.WithTx(tx)

		var impCount int64 = 0
		var delCount int64 = 0

		if overwrite {
			log.Debug("deleting existing domains")
			delCount, err = qtx.DeleteAll(ctx)
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Debugf("deleted %d domains", delCount)
		}

		log.Debug("inserting imported domains")
		now := time.Now().UnixMilli()
		for _, domain := range domains {
			if !overwrite {
				log.Debug("checking for existing domain", "domain", domain)

				count, err := qtx.CountExcludedDomain(ctx, domain)
				if err != nil {
					log.Fatal(err)
				}

				if count > 0 {
					log.Debug("domain exists, skipping import", "domain", domain)
					continue
				}
			}

			log.Debug("inserting", "domain", domain, "modificationTime", now)
			err := qtx.CreateExcludedDomain(ctx, repository.CreateExcludedDomainParams{Domain: domain, Now: now})
			if err != nil {
				log.Fatal(err)
				return
			}

			// We need to insert two rows for each domain
			impCount += 2
		}

		log.Debug("committing transaction")
		if err = tx.Commit(); err != nil {
			return
		}

		log.Infof("imported %d new origins, deleted %d existing origins", impCount, delCount)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	importCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing domains, any domains not in the export file will be deleted")
}
