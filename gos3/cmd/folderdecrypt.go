package cmd

import (
	"fmt"
	"gos3/internal/backupops"
	"gos3/internal/config"

	"github.com/spf13/cobra"
)

var folderdecryptCmd = &cobra.Command{
	Use:   "folderdecrypt <folder_containing_files> <private_key_filename>",
	Short: "Decrypt a folder and delete encrypted files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}
		return backupops.KeyDecrypt2Folder(args[0], args[1], configuration)
	},
}
