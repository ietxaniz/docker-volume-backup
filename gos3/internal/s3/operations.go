package s3

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gos3/internal/config"
	"gos3/internal/script"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func GenerateSubfolderName(frequency string) string {
	now := time.Now()
	switch frequency {
	case "daily":
		return now.Format("2006-01-02")
	case "weekly":
		year, week := now.ISOWeek()
		return fmt.Sprintf("%d-W%02d", year, week)
	case "hourly":
		return fmt.Sprintf("%s-%02d", now.Format("2006-01-02"), now.Hour())
	case "4hourly":
		return fmt.Sprintf("%s-%02d", now.Format("2006-01-02"), (now.Hour()/4)*4)
	default:
		return now.Format("2006-01-02")
	}
}

func EnsureS3FolderExists(folderPath string, cfg config.S3Config) error {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(cfg.Region),
		Endpoint: aws.String(cfg.Endpoint),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	svc := s3.New(sess)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(folderPath + "/"),
		Body:   strings.NewReader(""),
	})
	if err != nil {
		return fmt.Errorf("failed to create folder in S3: %w", err)
	}

	return nil
}

func UploadFolderToS3(localFolder, s3Folder string, cfg config.Config) error {
	dateSubfolder := GenerateSubfolderName(cfg.App.BackupFrequency)
	s3FullPath := filepath.Join(s3Folder, dateSubfolder)

	err := filepath.Walk(localFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(localFolder, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		s3Path := filepath.Join(s3FullPath, relPath)
		s3Path = strings.ReplaceAll(s3Path, "\\", "/") // Ensure forward slashes for S3 paths

		if info.Size() > parseSize(cfg.S3.MaxFileSize) {
			fmt.Printf("File %s exceeds maximum size. Splitting...\n", path)
			splitFiles, err := splitFile(path, cfg)
			if err != nil {
				return fmt.Errorf("failed to split file %s: %w", path, err)
			}

			for _, splitFile := range splitFiles {
				splitRelPath, _ := filepath.Rel(filepath.Dir(path), splitFile)
				splitS3Path := filepath.Join(s3FullPath, splitRelPath)
				splitS3Path = strings.ReplaceAll(splitS3Path, "\\", "/")

				fmt.Printf("Uploading split file %s to s3://%s/%s\n", splitFile, cfg.S3.Bucket, splitS3Path)
				err = UploadToS3(splitFile, splitS3Path, cfg.S3)
				if err != nil {
					return fmt.Errorf("failed to upload split file %s: %w", splitFile, err)
				}
			}
		} else {
			fmt.Printf("Uploading %s to s3://%s/%s\n", path, cfg.S3.Bucket, s3Path)
			err = UploadToS3(path, s3Path, cfg.S3)
			if err != nil {
				return fmt.Errorf("failed to upload file %s: %w", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking through local folder: %w", err)
	}

	return nil
}

func splitFile(filePath string, cfg config.Config) ([]string, error) {
	tempDir, err := os.MkdirTemp("", "split_files")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	err = script.Split(filePath, cfg.S3.MaxFileSize, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to split file: %w", err)
	}

	splitFiles, err := filepath.Glob(filepath.Join(tempDir, "*"))
	if err != nil {
		return nil, fmt.Errorf("failed to get split files: %w", err)
	}

	return splitFiles, nil
}

func UploadToS3(localPath, remotePath string, cfg config.S3Config) error {
	log.Printf("Starting upload of %s to S3 path %s", localPath, remotePath)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Endpoint:    aws.String(cfg.Endpoint),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.AccessKeySecret, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", localPath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	maxSize := parseSize(cfg.MaxFileSize)
	if fileInfo.Size() > maxSize {
		return uploadLargeFile(sess, localPath, remotePath, cfg, maxSize)
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(remotePath),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	log.Printf("Successfully uploaded %s to S3", localPath)
	return nil
}

func uploadLargeFile(sess *session.Session, localPath, remotePath string, cfg config.S3Config, maxSize int64) error {
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", localPath, err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	totalSize := fileInfo.Size()
	totalParts := int(math.Ceil(float64(totalSize) / float64(maxSize)))

	log.Printf("File %s exceeds max size. Splitting into %d parts", localPath, totalParts)

	// Create a subfolder for split files
	fileName := filepath.Base(localPath)
	subfolderPath := filepath.Join(filepath.Dir(remotePath), strings.TrimSuffix(fileName, filepath.Ext(fileName)))
	
	uploader := s3manager.NewUploader(sess)
	
	splitFiles := []string{}

	for i := 0; i < totalParts; i++ {
		partSize := int64(math.Min(float64(maxSize), float64(totalSize-int64(i)*maxSize)))
		partBuffer := make([]byte, partSize)
		
		_, err := file.Read(partBuffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file part: %w", err)
		}

		partNumber := i + 1
		splitFileName := fmt.Sprintf("%s.part%d", fileName, partNumber)
		splitFilePath := filepath.Join(filepath.Dir(localPath), splitFileName)
		
		err = os.WriteFile(splitFilePath, partBuffer, 0644)
		if err != nil {
			return fmt.Errorf("failed to write split file: %w", err)
		}
		
		splitFiles = append(splitFiles, splitFilePath)

		partRemotePath := filepath.Join(subfolderPath, splitFileName)
		
		log.Printf("Uploading part %d/%d of %s", partNumber, totalParts, fileName)
		
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(cfg.Bucket),
			Key:    aws.String(partRemotePath),
			Body:   strings.NewReader(string(partBuffer)),
		})
		if err != nil {
			return fmt.Errorf("failed to upload file part to S3: %w", err)
		}
	}

	log.Printf("Successfully uploaded all %d parts of %s", totalParts, fileName)

	// Clean up split files
	for _, splitFile := range splitFiles {
		err := os.Remove(splitFile)
		if err != nil {
			log.Printf("Warning: failed to remove split file %s: %v", splitFile, err)
		}
	}

	return nil
}

func parseSize(size string) int64 {
	var multiplier int64 = 1
	if strings.HasSuffix(size, "K") {
		multiplier = 1024
	} else if strings.HasSuffix(size, "M") {
		multiplier = 1024 * 1024
	} else if strings.HasSuffix(size, "G") {
		multiplier = 1024 * 1024 * 1024
	}

	numericSize, _ := strconv.ParseInt(strings.TrimRight(size, "KMG"), 10, 64)
	return numericSize * multiplier
}
