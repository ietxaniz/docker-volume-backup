package cmd

import (
	"gos3/internal/config"
	"gos3/internal/script"

	"github.com/spf13/cobra"
)

var splitCmd = &cobra.Command{
	Use:   "split <folder_path> <split_size>",
	Short: "Split files in a folder",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return err
		}
		return script.Split(args[0], args[1], configuration)
	},
}

var joinCmd = &cobra.Command{
	Use:   "join <folder_path>",
	Short: "Join split files in a folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return err
		}
		return script.Join(args[0], configuration)
	},
}
