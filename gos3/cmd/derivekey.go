package cmd

import (
	"gos3/internal/config"
	"gos3/internal/script"

	"github.com/spf13/cobra"
)

var derivekeyCmd = &cobra.Command{
	Use:   "derivekey",
	Short: "derivekey",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return err
		}
		password := args[0]
		salt, _ := cmd.Flags().GetString("salt")
		iterations, _ := cmd.Flags().GetInt("iterations")
		result, err := script.DeriveKey(password, salt, iterations, configuration)
		if err != nil {
			return err
		}
		script.PrintDeriveKeyResult(result)
		return nil
	},
}
