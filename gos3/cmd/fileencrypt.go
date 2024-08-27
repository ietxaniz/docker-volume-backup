package cmd

import (
	"fmt"
	"gos3/internal/config"
	"gos3/internal/script"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var fileEncryptCmd = &cobra.Command{
	Use:   "fileencrypt <input_file> <output_encrypted_file>",
	Short: "Encrypt a file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		fmt.Print("Enter encryption password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println()

		err = script.FileEncrypt(args[0], args[1], string(password), configuration)
		if err != nil {
			return err
		}

		fmt.Println("File encrypted successfully.")
		return nil
	},
}
