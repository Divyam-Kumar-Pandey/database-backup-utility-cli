package utils

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)


func CompressFile(inputPath string) (string, error) {
	in, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer in.Close()

	outputPath := inputPath + ".gz"
	out, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	gzWriter := gzip.NewWriter(out)
	gzWriter.Close()

	if _, err := io.Copy(gzWriter, in); err != nil {
		return "", err
	}

	return outputPath, nil
}

func DecompressFile(gzipPath string) (string, error) {
	if filepath.Ext(gzipPath) != ".gz" {
		return "", fmt.Errorf("not a gzip file")
	}

	in, err := os.Open(gzipPath)
	if err != nil {
		return "", err
	}
	defer in.Close()

	gzReader, err := gzip.NewReader(in)
	if err != nil {
		return "", err
	}
	defer gzReader.Close()

	outputPath := strings.TrimSuffix(gzipPath, ".gz")
	out, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

    if _, err := io.Copy(out, gzReader); err != nil {
		return "", err
	}

	return outputPath, nil
}

