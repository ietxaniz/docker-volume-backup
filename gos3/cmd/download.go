package cmd

import (
	"fmt"
	"gos3/internal/config"
	"gos3/internal/s3"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download and decrypt backups",
	Long:  `Download backups from S3, join split files if necessary, and decrypt the files`,
	RunE:  runDownload,
}

func runDownload(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfiguration("")
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	return s3.DownloadBackup(cfg)
}
