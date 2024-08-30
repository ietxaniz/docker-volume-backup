package cmd

import (
	"fmt"
	"gos3/internal/config"
	"gos3/internal/s3"
	"log"

	"github.com/spf13/cobra"
)

var testJoinDecryptCmd = &cobra.Command{
	Use:   "testjoindecrypt",
	Short: "Test joining split files",
	RunE:  runTestJoinDecrypt,
}

func runTestJoinDecrypt(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfiguration("")
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	log.Println("Starting test join process...")

	err = s3.JoinSplitFiles(cfg)
	if err != nil {
		return fmt.Errorf("failed to join split files: %w", err)
	}

	log.Println("Test join process completed successfully.")
	return nil
}
