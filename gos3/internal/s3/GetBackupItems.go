package s3

import (
	"fmt"
	"gos3/internal/config"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type BackupItem struct {
	Name         string
	DataItem     string
	PassItem     string
	IsDataFolder bool
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

func DonwloadBackupItem(item BackupItem, cfg config.Config) {
	// TODO: If it is folder download all the files in the folder and if it is file download the file. Both using DataItem string as full path
	// TODO: Download also PassItem as full path
	// TODO: If it is a folder call to script.Join
	// All items should be downloaded to localBackupFolder
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
				})
			}
		}
	}
	return items
}
