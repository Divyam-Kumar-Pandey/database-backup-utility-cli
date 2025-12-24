package utils

import (
	"compress/gzip"
	"io"
	"os"
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

