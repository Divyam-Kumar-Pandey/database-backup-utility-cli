package main

import (
	"db-backup-cli/pkg/core"
	"db-backup-cli/pkg/databases"
	"db-backup-cli/pkg/storage"
	"db-backup-cli/pkg/utils"
	"fmt"
	"os"
	"path/filepath"

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
		storageConfig := viper.GetStringMap("storage")
		dbType := dbConfig["type"].(string)

		storageAdaptor, err := getStorageAdaptor(storageConfig["type"].(string))
		if err != nil {
			fmt.Println(err)
			return
		}

		

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

		finalPath := filepath.Join(storageConfig["path"].(string), filepath.Base(compressedPath))

		uploadedPath, err := storageAdaptor.Upload(compressedPath, finalPath)
		if err != nil {
			fmt.Println("Upload failed:", err)
			return
		}
		os.Remove(compressedPath)
		fmt.Println("backup created at: ", uploadedPath)
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

func getStorageAdaptor(storageType string) (core.Storage, error) {
	switch storageType {
	case "local":
		return &storage.LocalStorage{}, nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
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