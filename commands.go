package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rene",
	Short: "René, twitter figur bot manager and automater",
	Run: func(cmd *cobra.Command, args []string) {
		ListenAndServe()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Setup a René property/account",
	Run: func(cmd *cobra.Command, args []string) {

		var count int
		err := Db.QueryRow("SELECT count(*) FROM accounts WHERE username = ?",
			args[0]).Scan(&count)

		if err != nil {
			logger.Error(
				`Something went wrong while reading from database`, err)
			return
		}
		if count != 0 {
			fmt.Println(`
The account you are trying to set is persistent in the database, please specifiy
flags if you wish to configure it`)
			return
		}

		for _, value := range args {
			err = AddAccount(value)
			if err != nil {
				logger.Info(fmt.Sprintf(
					`Couldn't add account %s to accounts database`, args[0]), err)
				return
			}
			logger.Info(fmt.Sprintf(
				`Successfully added %s to accounts database`, value))
		}
	},
}
