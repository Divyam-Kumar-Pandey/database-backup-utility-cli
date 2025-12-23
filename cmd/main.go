package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func initConfig() {
	viper.SetConfigName("db_backup_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found, run `backup-tool init`")
	} else {
		fmt.Println("Config file initialized.")
	}
}


func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(initCmd)
}

func main() {
	execute()
}