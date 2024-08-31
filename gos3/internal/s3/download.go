package s3

import (
	"bufio"
	"fmt"
	"gos3/internal/config"
	"gos3/internal/script"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type BackupDate struct {
	FolderName string
}

func GetBackupDates(cfg config.Config) ([]BackupDate, error) {
	sess, err := createS3Session(cfg.S3)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 session: %w", err)
	}

	svc := s3.New(sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(cfg.S3.Bucket),
		Prefix:    aws.String(cfg.S3.BackupFolder + "/"),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var dates []BackupDate
	for _, prefix := range resp.CommonPrefixes {
		folderName := strings.TrimPrefix(*prefix.Prefix, cfg.S3.BackupFolder)
		folderName = strings.TrimPrefix(folderName, "/")
		folderName = strings.TrimSuffix(folderName, "/")

		if folderName != "" {
			dates = append(dates, BackupDate{
				FolderName: folderName,
			})
		}
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].FolderName > dates[j].FolderName
	})

	return dates, nil
}

func DownloadBackup(cfg config.Config) error {
	dates, err := GetBackupDates(cfg)
	if err != nil {
		return fmt.Errorf("failed to get backup dates: %w", err)
	}

	if len(dates) == 0 {
		return fmt.Errorf("no backups found")
	}

	fmt.Println("Available backup dates:")
	for i, date := range dates {
		fmt.Printf("%d - %s\n", i+1, date.FolderName)
	}

	selectedDate, err := selectBackupDate(dates)
	if err != nil {
		return err
	}

	backupItems, err := GetBackupItems(cfg, selectedDate)
	if err != nil {
		return err
	}
	for _, backupItem := range backupItems {
		err = DownloadBackupItem(backupItem, cfg)
		if err != nil {
			return err
		}
	}

	err = script.Join(cfg.App.LocalBackupFolder, cfg)
	if err != nil {
		return fmt.Errorf("failed to join split files: %w", err)
	}

	return nil
}

func selectBackupDate(dates []BackupDate) (BackupDate, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter the number of the backup date you want to download: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return BackupDate{}, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(dates) {
			fmt.Println("Invalid input. Please enter a number between 1 and", len(dates))
			continue
		}

		return dates[index-1], nil
	}
}
