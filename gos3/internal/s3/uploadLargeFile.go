package s3

import (
	"fmt"
	"gos3/internal/config"
	"gos3/internal/script"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func uploadLargeFile(sess *session.Session, localPath, remotePath string, cfg config.Config, maxSize int64) error {
	log.Printf("File %s exceeds max size. Splitting into parts", localPath)

	// Create a temporary directory for split files
	tempDir, err := os.MkdirTemp("", "split_files")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	err = script.Split(localPath, fmt.Sprintf("%d", maxSize), cfg)
	if err != nil {
		return fmt.Errorf("failed to split file: %w", err)
	}

	// Create a subfolder for split files in S3
	fileName := filepath.Base(localPath)
	subfolderPath := filepath.Join(filepath.Dir(remotePath), strings.TrimSuffix(fileName, filepath.Ext(fileName)))

	uploader := s3manager.NewUploader(sess)

	// Upload each split file
	splitFiles, err := filepath.Glob(filepath.Join(tempDir, "*"))
	if err != nil {
		return fmt.Errorf("failed to list split files: %w", err)
	}

	for i, splitFile := range splitFiles {
		file, err := os.Open(splitFile)
		if err != nil {
			return fmt.Errorf("failed to open split file %s: %w", splitFile, err)
		}
		defer file.Close()

		splitFileName := filepath.Base(splitFile)
		partRemotePath := filepath.Join(subfolderPath, splitFileName)

		log.Printf("Uploading part %d/%d of %s", i+1, len(splitFiles), fileName)

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(cfg.S3.Bucket),
			Key:    aws.String(partRemotePath),
			Body:   file,
		})
		if err != nil {
			return fmt.Errorf("failed to upload file part to S3: %w", err)
		}
	}

	log.Printf("Successfully uploaded all %d parts of %s", len(splitFiles), fileName)

	return nil
}
