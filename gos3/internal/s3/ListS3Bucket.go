package s3

import (
	"fmt"
	"strings"

	"gos3/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Item struct {
	Name            string `json:"name" yaml:"name"`
	IsFolder        bool   `json:"isFolder" yaml:"isFolder"`
	LastModified    int64  `json:"lastModified" yaml:"lastModified"`
	LastModifiedStr string `json:"lastModifiedStr" yaml:"lastModifiedStr"`
	Size            int64  `json:"size" yaml:"size"`
}

func ListS3Bucket(s3config config.Config, prefix string, delimiter string) ([]S3Item, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3config.S3.Region),
		Endpoint:    aws.String(s3config.S3.Endpoint),
		Credentials: credentials.NewStaticCredentials(s3config.S3.AccessKeyID, s3config.S3.AccessKeySecret, ""),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	svc := s3.New(sess)

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s3config.S3.Bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String(delimiter),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("unable to list items in bucket %q, %v", s3config.S3.Bucket, err)
	}

	var items []S3Item

	for _, prefix := range result.CommonPrefixes {
		items = append(items, S3Item{
			Name:     *prefix.Prefix,
			IsFolder: true,
		})
	}

	for _, item := range result.Contents {
		lastModifiedUnix := item.LastModified.UnixMilli()
		lastModifiedStr := item.LastModified.UTC().Format("2006-01-02 15:04:05.000")

		items = append(items, S3Item{
			Name:            *item.Key,
			LastModified:    lastModifiedUnix,
			LastModifiedStr: lastModifiedStr,
			Size:            *item.Size,
		})
	}

	return items, nil
}

func PrintS3ItemList(items []S3Item) {
	maxNameWidth := len("Name")
	maxDateWidth := len("Last-Modified")
	maxSizeWidth := len("Size")
	maxIsfolderWidth := len("Is-Folder")

	for _, item := range items {
		if len(item.Name) > maxNameWidth {
			maxNameWidth = len(item.Name)
		}
		if len(item.LastModifiedStr) > maxDateWidth {
			maxDateWidth = len(item.LastModifiedStr)
		}
		sizeWidth := len(fmt.Sprintf("%d", item.Size))
		if sizeWidth > maxSizeWidth {
			maxSizeWidth = sizeWidth
		}
	}

	fmt.Printf("%-*s  %-*s  %-*s  %-*s\n", maxNameWidth, "Name", maxDateWidth, "Last-Modified", maxSizeWidth, "Size", maxIsfolderWidth, "Is-Folder")
	fmt.Println(strings.Repeat("-", maxNameWidth+maxDateWidth+maxSizeWidth+maxIsfolderWidth+6))

	for _, item := range items {
		isFolder := 0
		if item.IsFolder {
			isFolder = 1
		}
		fmt.Printf("%-*s  %-*s  %-*d  %-*d\n",
			maxNameWidth, item.Name,
			maxDateWidth, item.LastModifiedStr,
			maxSizeWidth, item.Size,
			maxIsfolderWidth, isFolder)
	}
}
