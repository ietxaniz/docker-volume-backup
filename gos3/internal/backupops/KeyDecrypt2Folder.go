package backupops

import (
	"fmt"
	"gos3/internal/config"
	"gos3/internal/script"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

func KeyDecrypt2Folder(workingFolder, encryptedPrivateKeyFile string, configuration config.Config) error {
	// Prompt for private key password
	fmt.Print("Enter private key decryption password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	privateKeyPassword := string(passwordBytes)
	fmt.Println() // Print a newline after password input

	// Get list of encrypted files
	encryptedFiles, err := listEncryptedFiles(workingFolder)
	if err != nil {
		return fmt.Errorf("failed to list encrypted files: %w", err)
	}

	// Decrypt each file
	for i, encryptedFile := range encryptedFiles {
		fmt.Printf("Decrypting file %d of %d: %s\n", i+1, len(encryptedFiles), encryptedFile)

		outputFile := filepath.Join(workingFolder, strings.TrimSuffix(filepath.Base(encryptedFile), ".cpt"))
		err := script.KeyDecrypt2WithPass(encryptedFile, outputFile, encryptedPrivateKeyFile, privateKeyPassword, configuration)
		if err != nil {
			return fmt.Errorf("failed to decrypt %s: %w", encryptedFile, err)
		}

		// Delete original encrypted file and its .pass file
		if err := deleteEncryptedFiles(encryptedFile); err != nil {
			fmt.Printf("Warning: Failed to delete %s: %v\n", encryptedFile, err)
		}
	}

	fmt.Println("Folder decryption completed successfully.")
	return nil
}

func listEncryptedFiles(folder string) ([]string, error) {
	var encryptedFiles []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".cpt") {
			passFile := path + ".pass"
			if _, err := os.Stat(passFile); err == nil {
				encryptedFiles = append(encryptedFiles, path)
			}
		}
		return nil
	})
	return encryptedFiles, err
}

func deleteEncryptedFiles(encryptedFile string) error {
	if err := os.Remove(encryptedFile); err != nil {
		return fmt.Errorf("failed to delete encrypted file: %w", err)
	}
	passFile := encryptedFile + ".pass"
	if err := os.Remove(passFile); err != nil {
		return fmt.Errorf("failed to delete pass file: %w", err)
	}
	return nil
}
