package cmd

import (
	"gos3/internal/config"
	"gos3/internal/s3"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list",
	RunE: func(cmd *cobra.Command, args []string) error {
		s3config, err := config.LoadConfiguration("")
		if err != nil {
			return err
		}
		prefix, _ := cmd.Flags().GetString("prefix")
		delimiter, _ := cmd.Flags().GetString("delimiter")
		items, err := s3.ListS3Bucket(s3config, prefix, delimiter)
		if err != nil {
			return err
		}
		s3.PrintS3ItemList(items)
		return nil
	},
}
