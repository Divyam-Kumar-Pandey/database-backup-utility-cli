package main

import (
	"db-backup-cli/pkg/core"
	"db-backup-cli/pkg/databases"
	"db-backup-cli/pkg/utils"
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

var backupCmd = &cobra.Command{
	Use: "backup [database]",
	Short: "Backup a Database",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var databaseName = args[0]

		if !viper.IsSet("databases." + databaseName) {
			fmt.Printf("Database '%s' not found in config\n", databaseName); 
			return
		}

		dbConfig := viper.GetStringMap("databases." + databaseName)
		dbType := dbConfig["type"].(string)

		dbAdapter, err := getDatabaseAdapter(dbType)
		if err != nil {
			fmt.Println(err); return;
		}

		outputFile := databaseName + "_backup.sql";
		path, err := dbAdapter.Backup(dbConfig, outputFile)
		if err != nil {
			fmt.Println("Backup failed", err)
			return;
		}

		compressedPath, err := utils.CompressFile(path)
		if err != nil {
			fmt.Println("Compression failed:", err)
			return
		}
		os.Remove(path)
		fmt.Println("backup created at: ", compressedPath)
	},
}

func getDatabaseAdapter(dbType string) (core.Database, error) {
	switch dbType {
	case "mysql":
		return &databases.MySQLDatabase{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}



func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(backupCmd)
}

func main() {
	execute()
}