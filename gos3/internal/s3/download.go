package s3

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"gos3/internal/config"
	"gos3/internal/script"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/term"
)

func ListBackupVersions(cfg config.Config) ([]string, error) {
	sess, err := createS3Session(cfg.S3)
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(cfg.S3.Bucket),
		Prefix: aws.String(cfg.S3.BackupFolder),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	versions := make(map[string]struct{})
	for _, item := range resp.Contents {
		parts := strings.Split(*item.Key, "/")
		if len(parts) > 1 {
			versions[parts[1]] = struct{}{}
		}
	}

	result := make([]string, 0, len(versions))
	for v := range versions {
		result = append(result, v)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(result)))

	return result, nil
}

func DownloadBackup(cfg config.Config) error {
	versions, err := ListBackupVersions(cfg)
	if err != nil {
		return err
	}

	fmt.Println("Available backup versions:")
	for i, v := range versions {
		fmt.Printf("%d. %s\n", i+1, v)
	}

	var choice int
	fmt.Print("Enter the number of the version you want to download: ")
	_, err = fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > len(versions) {
		return fmt.Errorf("invalid choice")
	}

	selectedVersion := versions[choice-1]
	log.Printf("Downloading backup version: %s", selectedVersion)

	sess, err := createS3Session(cfg.S3)
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)

	prefix := filepath.Join(cfg.S3.BackupFolder, selectedVersion)
	err = svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(cfg.S3.Bucket),
		Prefix: aws.String(prefix),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, item := range page.Contents {
			localPath := filepath.Join(cfg.App.LocalBackupFolder, strings.TrimPrefix(*item.Key, prefix))
			if err := downloadFile(downloader, cfg.S3.Bucket, *item.Key, localPath); err != nil {
				log.Printf("Error downloading %s: %v", *item.Key, err)
				continue
			}
		}
		return true
	})

	if err != nil {
		return fmt.Errorf("failed to list objects: %w", err)
	}

	log.Println("All files downloaded. Joining split files...")
	if err := joinSplitFiles(cfg.App.LocalBackupFolder); err != nil {
		return fmt.Errorf("failed to join split files: %w", err)
	}

	log.Println("Split files joined. Decrypting files...")
	fmt.Print("Enter the private key password: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()

	if err := decryptFiles(cfg, string(password)); err != nil {
		return fmt.Errorf("failed to decrypt files: %w", err)
	}

	log.Println("Download process completed successfully.")
	return nil
}

func downloadFile(downloader *s3manager.Downloader, bucket, key, localPath string) error {
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	defer file.Close()

	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to download file %s: %w", key, err)
	}

	log.Printf("Downloaded file: %s", localPath)
	return nil
}

func joinSplitFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	splitFiles := make(map[string][]string)
	for _, file := range files {
		if strings.Contains(file.Name(), ".part") {
			baseName := strings.Split(file.Name(), ".part")[0]
			splitFiles[baseName] = append(splitFiles[baseName], file.Name())
		}
	}

	for baseName, parts := range splitFiles {
		sort.Strings(parts)
		outputFile, err := os.Create(filepath.Join(dir, baseName))
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outputFile.Close()

		for _, part := range parts {
			partPath := filepath.Join(dir, part)
			partFile, err := os.Open(partPath)
			if err != nil {
				return fmt.Errorf("failed to open part file: %w", err)
			}
			_, err = io.Copy(outputFile, partFile)
			partFile.Close()
			if err != nil {
				return fmt.Errorf("failed to copy part file: %w", err)
			}
			os.Remove(partPath)
		}
		log.Printf("Joined file: %s", baseName)
	}

	return nil
}

func decryptFiles(cfg config.Config, password string) error {
	files, err := os.ReadDir(cfg.App.LocalBackupFolder)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".cpt") {
			inputFile := filepath.Join(cfg.App.LocalBackupFolder, file.Name())
			outputFile := strings.TrimSuffix(inputFile, ".cpt")

			log.Printf("Decrypting file: %s", inputFile)
			err := script.KeyDecrypt2(inputFile, outputFile, cfg.App.PrivateKeyFile, password, cfg)
			if err != nil {
				return fmt.Errorf("failed to decrypt file %s: %w", file.Name(), err)
			}

			os.Remove(inputFile)
			os.Remove(inputFile + ".pass")
			log.Printf("Decrypted file: %s", outputFile)
		}
	}

	return nil
}

func createS3Session(cfg config.S3Config) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Endpoint:    aws.String(cfg.Endpoint),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.AccessKeySecret, ""),
	})
}
