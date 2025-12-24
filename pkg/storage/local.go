package storage

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct{}

func (l *LocalStorage) Upload(localPath, remotePath string) (string, error) {
	if err := os.MkdirAll(filepath.Dir(remotePath), 0755); err != nil {
		return "", nil
	}

	src, err := os.Open(localPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(remotePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return remotePath, nil
}

func (l *LocalStorage) Download(remotePath, localPath string) (string, error) {
	src, err := os.Open(remotePath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", nil
	}

	return localPath, nil
}

func (l *LocalStorage) ListFiles(path string) ([]string, error) {
	return []string{}, nil
}