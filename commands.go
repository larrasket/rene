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
		for _, value := range args {
			err := AddAccount(value)
			if err != nil {
				logger.Info(fmt.Sprintf(
					`Couldn't add account %s to accounts database`, value))
			}
		}
	},
}
