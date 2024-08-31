package s3

import (
	"fmt"
	"gos3/internal/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type BackupItem struct {
	Name         string
	DataItem     string
	PassItem     string
	IsDataFolder bool
	S3BaseFolder string
}

func GetBackupItems(cfg config.Config, date BackupDate) ([]BackupItem, error) {
	sess, err := createS3Session(cfg.S3)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 session: %w", err)
	}

	svc := s3.New(sess)

	parentFolder := cfg.S3.BackupFolder + "/" + date.FolderName + "/"

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(cfg.S3.Bucket),
		Prefix:    aws.String(parentFolder),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	folderItems := make([]string, 0)

	for _, prefix := range resp.CommonPrefixes {
		folderName := strings.TrimPrefix(*prefix.Prefix, cfg.S3.BackupFolder)
		folderName = strings.TrimPrefix(folderName, "/")
		folderName = strings.TrimSuffix(folderName, "/")

		if folderName != "" {
			folderItems = append(folderItems, *prefix.Prefix)
		}
	}

	fileItems := make([]string, 0)
	for _, obj := range resp.Contents {
		fileItems = append(fileItems, *obj.Key)
	}

	return getBackupItemsFromFolderAndFiles(parentFolder, folderItems, fileItems), nil

}

func DownloadBackupItem(item BackupItem, cfg config.Config) error {
	sess, err := createS3Session(cfg.S3)
	if err != nil {
		return fmt.Errorf("failed to create S3 session: %w", err)
	}

	downloader := s3manager.NewDownloader(sess)

	if item.IsDataFolder {
		err = downloadFolder(sess, downloader, item.S3BaseFolder, item.DataItem, cfg)
	} else {
		err = downloadFile(downloader, cfg.S3.Bucket, item.S3BaseFolder, item.DataItem, cfg)
	}
	if err != nil {
		return fmt.Errorf("failed to download data item: %w", err)
	}

	err = downloadFile(downloader, cfg.S3.Bucket, item.S3BaseFolder, item.PassItem, cfg)
	if err != nil {
		return fmt.Errorf("failed to download pass item: %w", err)
	}

	return nil
}

func getBackupItemsFromFolderAndFiles(parentFolder string, folderItems []string, fileItems []string) []BackupItem {
	items := make([]BackupItem, 0)
	for _, folder := range folderItems {
		name := strings.TrimPrefix(folder, parentFolder)
		name = strings.TrimSuffix(name, ".cpt-split_parts/")
		containsData := false
		containsPass := false
		dataName := parentFolder + name + ".cpt-split_parts/"
		passName := parentFolder + name + ".cpt.pass"
		for _, data := range folderItems {
			if data == dataName {
				containsData = true
				break
			}
		}
		for _, pass := range fileItems {
			if pass == passName {
				containsPass = true
				break
			}
		}
		if containsData && containsPass {
			items = append(items, BackupItem{
				Name:         name,
				DataItem:     dataName,
				PassItem:     passName,
				IsDataFolder: true,
				S3BaseFolder: parentFolder,
			})
		}
	}
	for _, file := range fileItems {
		if strings.HasSuffix(file, ".cpt") {
			name := strings.TrimPrefix(file, parentFolder)
			name = strings.TrimSuffix(name, ".cpt")
			containsData := false
			containsPass := false
			dataName := parentFolder + name + ".cpt"
			passName := parentFolder + name + ".cpt.pass"
			for _, data := range fileItems {
				if data == dataName {
					containsData = true
					break
				}
			}
			for _, pass := range fileItems {
				if pass == passName {
					containsPass = true
					break
				}
			}
			if containsData && containsPass {
				items = append(items, BackupItem{
					Name:         name,
					DataItem:     dataName,
					PassItem:     passName,
					IsDataFolder: true,
					S3BaseFolder: parentFolder,
				})
			}
		}
	}
	return items
}

func downloadFolder(sess *session.Session, downloader *s3manager.Downloader, baseFolder string, folderPath string, cfg config.Config) error {
	svc := s3.New(sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(cfg.S3.Bucket),
		Prefix: aws.String(folderPath),
	})
	if err != nil {
		return fmt.Errorf("failed to list objects in folder: %w", err)
	}

	for _, item := range resp.Contents {
		err := downloadFile(downloader, cfg.S3.Bucket, baseFolder, *item.Key, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadFile(downloader *s3manager.Downloader, bucket, baseFolder, filePath string, cfg config.Config) error {
	localPath := filepath.Join(cfg.App.LocalBackupFolder, strings.TrimPrefix(filePath, baseFolder))

	err := os.MkdirAll(filepath.Dir(localPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	fmt.Printf("Downloaded: %s\n", localPath)
	return nil
}
