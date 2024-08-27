package cmd

import (
	"fmt"
	"gos3/internal/config"
	"gos3/internal/script"

	"github.com/spf13/cobra"
)

var volumerestoreCmd = &cobra.Command{
	Use:   "volumerestore <volume_name> <backup_file_name>",
	Short: "Restore a Docker volume from a backup",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		volumeName := args[0]
		backupFileName := args[1]

		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		err = script.VolumeRestore(volumeName, backupFileName, configuration)
		if err != nil {
			return fmt.Errorf("volume restore failed: %w", err)
		}

		fmt.Printf("Volume %s restored successfully from %s\n", volumeName, backupFileName)
		return nil
	},
}
