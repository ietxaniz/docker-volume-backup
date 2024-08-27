package cmd

import (
	"fmt"
	"gos3/internal/config"
	"gos3/internal/script"

	"github.com/spf13/cobra"
)

var volumebackupCmd = &cobra.Command{
	Use:   "volumebackup <volume_name> <backup_file_name>",
	Short: "Backup a Docker volume",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		volumeName := args[0]
		backupFileName := args[1]

		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		noCompression, _ := cmd.Flags().GetBool("no-compression")
		compress := !noCompression

		result, err := script.VolumeBackup(volumeName, backupFileName, compress, configuration)
		if err != nil {
			return fmt.Errorf("volume backup failed: %w", err)
		}

		fmt.Printf("Backup of volume %s created as %s\n", volumeName, backupFileName)
		fmt.Printf("Original size: %d bytes\n", result.OriginalSize)
		fmt.Printf("Final size: %d bytes\n", result.FinalSize)
		fmt.Printf("Compression ratio: %.2f\n", result.CompressionRatio)
		fmt.Printf("Time elapsed: %.6f seconds\n", result.TimeElapsed)

		return nil
	},
}
