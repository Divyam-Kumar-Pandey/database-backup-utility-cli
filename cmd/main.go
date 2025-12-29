package main

import (
	"db-backup-cli/pkg/core"
	"db-backup-cli/pkg/databases"
	"db-backup-cli/pkg/storage"
	"db-backup-cli/pkg/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "db-backup",
	Short: "Database Backup Utility Tool CLI",
}

func execute() {
	err := rootCmd.Execute()
	if err != nil {
		handleError("Error executing command", err)
		os.Exit(1)
	}
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LogInfo("Configuration initialized.")
	},
}

func initConfig() {
	viper.SetConfigName("db_backup_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		handleError("No config file found, run `db-backup-cli init`", err)
	} else {
		utils.LogInfo("Config file initialized.")
	}
}

var backupCmd = &cobra.Command{
	Use:   "backup [database]",
	Short: "Backup a Database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utils.LogInfo("Starting backup...")
		var databaseName = args[0]

		if !viper.IsSet("databases." + databaseName) {
			handleError(fmt.Sprintf("Database '%s' not found in config", databaseName), nil)
			return
		}

		dbConfig := viper.GetStringMap("databases." + databaseName)
		pwdKey := fmt.Sprintf("databases.%s.password", databaseName)
		if pwd := viper.GetString(pwdKey); pwd != "" {
			dbConfig["password"] = pwd
		}
		storageConfig := viper.GetStringMap("storage")
		dbType := dbConfig["type"].(string)

		storageAdaptor, err := getStorageAdapter()
		if err != nil {
			handleError("Failed to get storage adapter", err)
			return
		}

		dbAdapter, err := getDatabaseAdapter(dbType)
		if err != nil {
			handleError("Failed to get database adapter", err)
			return
		}

		outputFile := databaseName + "_backup.sql"
		path, err := dbAdapter.Backup(dbConfig, outputFile)
		if err != nil {
			handleError("Backup failed", err)
			return
		}
		utils.LogInfo("Database dump completed successfully")

		compressedPath, err := utils.CompressFile(path)
		if err != nil {
			handleError("Compression failed", err)
			return
		}
		utils.LogInfo("Compression completed successfully")
		os.Remove(path)

		finalPath := filepath.Join(storageConfig["path"].(string), filepath.Base(compressedPath))

		uploadedPath, err := storageAdaptor.Upload(compressedPath, finalPath)
		if err != nil {
			handleError("Upload failed", err)
			return
		}
		utils.LogInfo("Upload completed successfully")
		os.Remove(compressedPath)
		utils.LogInfo("Backup created at: " + uploadedPath)
	},
}

func getDatabaseAdapter(dbType string) (core.Database, error) {
	switch dbType {
	case "mysql":
		return &databases.MySQLDatabase{}, nil
	case "sqlite":
		return &databases.SQLiteDatabase{}, nil
	case "mongo":
		return &databases.MongoDatabase{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

func getStorageAdapter() (core.Storage, error) {
	storageType := viper.GetString("storage.type")

	switch storageType {
	case "local":
		return &storage.LocalStorage{}, nil
	case "s3":
		bucket := viper.GetString("storage.bucket")
		region := viper.GetString("storage.region")
		return storage.NewS3Storage(bucket, region)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

func handleError(userMsg string, err error) {
	if err != nil {
		utils.LogError(userMsg + ": " + err.Error())
	} else {
		utils.LogError(userMsg)
	}
	fmt.Println(userMsg)
}

var restoreCmd = &cobra.Command{
	Use:   "restore [backup_file] [db_name]",
	Short: "Restore a Database",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utils.LogInfo("Starting restore...")
		backupFile := args[0]
		dbName := args[1]

		if !viper.IsSet("databases." + dbName) {
			handleError(fmt.Sprintf("Database not found: %s", dbName), nil)
			return
		}

		dbConfig := viper.GetStringMap("databases." + dbName)
		pwdKey := fmt.Sprintf("databases.%s.password", dbName)
		if pwd := viper.GetString(pwdKey); pwd != "" {
			dbConfig["password"] = pwd
		}
		// storageConfig := viper.GetStringMap("storage")
		dbType := dbConfig["type"].(string)

		dbAdapter, err := getDatabaseAdapter(dbType)
		if err != nil {
			handleError("Failed to get database adapter", err)
			return
		}

		storageAdapter, err := getStorageAdapter()
		if err != nil {
			handleError("Failed to get storage adapter", err)
			return
		}

		tempPath := filepath.Join(os.TempDir(), filepath.Base(backupFile))

		utils.LogInfo("Downloading backup...")
		localPath, err := storageAdapter.Download(backupFile, tempPath)
		if err != nil {
			handleError("Download failed", err)
			return
		}

		restorePath := localPath
		if strings.HasSuffix(localPath, ".gz") {
			utils.LogInfo("Decompressing backup...")
			restorePath, err = utils.DecompressFile(localPath)
			if err != nil {
				handleError("Decompression failed", err)
				return
			}
		}

		utils.LogInfo("Restoring database...")
		if err := dbAdapter.Restore(dbConfig, restorePath); err != nil {
			handleError("Restore failed", err)
			return
		}

		os.Remove(localPath)
		if restorePath != localPath {
			os.Remove(restorePath)
		}

		utils.LogInfo("Restore completed successfully")
	},
}

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage backup schedules",
}

var scheduleAddCmd = &cobra.Command{
	Use: "add [db_name] [cron]",
	Short: "Schedule a database backup",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dbName := args[0]
		cronExpr := args[1]

		binary, _ := exec.LookPath(os.Args[0])

		command := fmt.Sprintf(
			"%s backup %s >> %s 2>&1",
			binary,
			dbName,
			"backup_tool.log",
		)

		err := utils.AddCronJob(dbName, cronExpr, command)
		if err != nil {
			handleError("Failed to add schedule", err)
			return
		}

		utils.LogInfo("Backup scheduled successfully")
	},
}

var scheduleRemoveCmd = &cobra.Command{
	Use:   "remove [db_name]",
	Short: "Remove scheduled backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := utils.RemoveCronJob(args[0]); err != nil {
			handleError("Failed to remove schedule", err)
			return
		}

		utils.LogInfo("Schedule removed successfully")
	},
}



func init() {
	utils.InitLogger()
	cobra.OnInitialize(initConfig)

	scheduleCmd.AddCommand(scheduleAddCmd)
	scheduleCmd.AddCommand(scheduleRemoveCmd)

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(scheduleCmd)
}

func main() {
	defer utils.CloseLogger()
	execute()
}
