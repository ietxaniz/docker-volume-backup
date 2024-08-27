package cmd

import (
	"fmt"
	"syscall"

	"gos3/internal/config"
	"gos3/internal/script"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var keyGenerateCmd = &cobra.Command{
	Use:   "keygenerate",
	Short: "Generate a public-private key pair",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		return script.KeyGenerate(configuration)
	},
}

var keyEncryptCmd = &cobra.Command{
	Use:   "keyencrypt <input_file> <output_encrypted_file> <public_key_file>",
	Short: "Encrypt a file using a public key",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		return script.KeyEncrypt(args[0], args[1], args[2], configuration)
	},
}

var keyDecryptCmd = &cobra.Command{
	Use:   "keydecrypt <input_encrypted_file> <output_decrypted_file> <private_key_file>",
	Short: "Decrypt a file using a private key",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		return script.KeyDecrypt(args[0], args[1], args[2], configuration)
	},
}

var keyDecrypt2Cmd = &cobra.Command{
	Use:   "keydecrypt2 <input_encrypted_file> <output_decrypted_file> <encrypted_private_key_file>",
	Short: "Decrypt a file using an encrypted private key",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.LoadConfiguration("")
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		fmt.Print("Enter private key decryption password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println()

		return script.KeyDecrypt2(args[0], args[1], args[2], string(password), configuration)
	},
}
