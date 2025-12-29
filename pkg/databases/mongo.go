package databases

import (
	"db-backup-cli/pkg/core"
	"fmt"
	"os/exec"
)

type MongoDatabase struct{}

func (db *MongoDatabase) Backup(config core.Config, outputPath string) (string, error) {
	user := config["user"].(string)
	password := config["password"].(string)
	host := config["host"].(string)
	port := config["port"].(int)
	database := config["database"].(string)

	mongoUrl := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", user, password, host, port, database)

	cmd := exec.Command("mongodump", "--uri", mongoUrl, "--out", outputPath)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return outputPath, nil
}

func (db *MongoDatabase) Restore(config core.Config, backupPath string) error {
	user := config["user"].(string)
	password := config["password"].(string)
	host := config["host"].(string)
	port := config["port"].(int)
	database := config["database"].(string)

	mongoUrl := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", user, password, host, port, database)

	cmd := exec.Command("mongorestore", "--uri", mongoUrl, "--dir", backupPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (db *MongoDatabase) TestConnection(config core.Config) error {
	user := config["user"].(string)
	password := config["password"].(string)
	host := config["host"].(string)
	port := config["port"].(int)
	database := config["database"].(string)

	mongoUrl := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", user, password, host, port, database)

	cmd := exec.Command(
		"mongosh",
		mongoUrl,
		"--eval",
		"db.runCommand({ ping: 1 })",
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mongo test connection failed: %w", err)
	}

	return nil
}
