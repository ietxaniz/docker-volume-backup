package cmd

import (
	"fmt"

	"gos3/internal/config"
	"gos3/internal/s3"

	"github.com/spf13/cobra"
)

var s3UploadCmd = &cobra.Command{
	Use:   "s3upload",
	Short: "Upload the configured local folder to S3",
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		localFolderOverride, _ := cmd.Flags().GetString("local")
		s3FolderOverride, _ := cmd.Flags().GetString("s3folder")

		localFolder := configuration.App.LocalBackupFolder
		if localFolderOverride != "" {
			localFolder = localFolderOverride
		}

		s3Folder := configuration.S3.BackupFolder
		if s3FolderOverride != "" {
			s3Folder = s3FolderOverride
		}

		fmt.Printf("Uploading from %s to S3 folder %s\n", localFolder, s3Folder)

		err = s3.UploadFolderToS3(localFolder, s3Folder, configuration)
		if err != nil {
			return fmt.Errorf("failed to upload folder to S3: %w", err)
		}

		fmt.Println("Folder uploaded successfully to S3")
		return nil
	},
}
