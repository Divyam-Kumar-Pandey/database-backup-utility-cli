package databases

import (
	"db-backup-cli/pkg/core"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type SQLiteDatabase struct{}

func (db *SQLiteDatabase) Backup(config core.Config, outputPath string) (string, error) {
	dbPath := config["path"].(string)
	src, err := os.Open(dbPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", err
	}

	dst, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return outputPath, nil
}

func (db *SQLiteDatabase) Restore(config core.Config, backupPath string) error {
	dbPath := config["path"].(string)

	src, err := os.Open(backupPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dbPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func (db *SQLiteDatabase) TestConnection(config core.Config) error {
	dbPath := config["path"].(string)
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("database file does not exist: %s", dbPath)
	}
	if err != nil {
		return err
	}
	return nil
}
