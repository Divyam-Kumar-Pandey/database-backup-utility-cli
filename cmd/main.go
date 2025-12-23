package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "db-backup",
	Short: "Database Backup Utility Tool CLI",
}

func execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Configuration initialized.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func main() {
	execute()
}