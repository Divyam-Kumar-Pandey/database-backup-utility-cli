package databases

import (
	"db-backup-cli/pkg/core"
	"fmt"
	"os"
	"os/exec"
)

type MySQLDatabase struct{}

func (db *MySQLDatabase) Backup(config core.Config, outputPath string) (string, error) {
	user := config["user"].(string)
	password := config["password"].(string)
	host := config["host"].(string)
	port := config["port"].(int)
	database := config["database"].(string)

	cmd := exec.Command(
		"mysqldump",
		fmt.Sprintf("-u%s", user),
		fmt.Sprintf("-h%s", host),
		fmt.Sprintf("-P%d", port),
		database,
	)
	
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("MYSQL_PWD=%s", password),
	)

	file, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	cmd.Stdout = file
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("mysqldump failed: %w", err)
	}

	return outputPath, nil
}

func (db *MySQLDatabase) Restore(config core.Config, backupPath string) error {
	user := config["user"].(string)
	password := config["password"].(string)
	host := config["host"].(string)
	port := config["port"].(int)
	database := config["database"].(string)

	file, err := os.Open(backupPath)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd := exec.Command(
		"mysql",
		fmt.Sprintf("-u%s", user),
		fmt.Sprintf("-h%s", host),
		fmt.Sprintf("-P%d", port),
		database,
	)
	
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("MYSQL_PWD=%s", password),
	)

	cmd.Stdin = file
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysql restore failed: %w", err)
	}

	return nil
}

func (db *MySQLDatabase) TestConnection(config core.Config) error {
	return nil
}

