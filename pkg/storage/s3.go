package storage

import (
	"context"
	"db-backup-cli/pkg/core"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(bucket, region string) (core.Storage, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Storage{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *S3Storage) Upload(localPath, remotePath string) (string, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	key := filepath.ToSlash(remotePath)

	// TODO: we have to check if the bucket exists, if not create it

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key: aws.String(key),
		Body: file,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("s3://%s/%s", s.bucket, key), nil
}

func (s *S3Storage) Download(remotePath, localPath string) (string, error) {
	key := filepath.ToSlash(remotePath)

	resp, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key: aws.String(key),
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := out.ReadFrom(resp.Body); err != nil {
		return "", err
	}

	return localPath, nil
}

func (s *S3Storage) ListFiles(prefix string) ([]string, error) {
	resp, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	var files []string
	for _, obj := range resp.Contents {
		files = append(files, *obj.Key)
	}

	return files, nil
}
