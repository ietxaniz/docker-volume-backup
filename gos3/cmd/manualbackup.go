package cmd

import (
	"fmt"
	"gos3/internal/config"
	"gos3/internal/backupops"
	"log"

	"github.com/spf13/cobra"
)

var manualBackupCmd = &cobra.Command{
	Use:   "manualbackup",
	Short: "Perform a manual backup of all defined backups",
	Long:  `Perform a manual backup of all backup definitions from the configuration.`,
	RunE:  runManualBackup,
}


func runManualBackup(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfiguration("")
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	log.Println("Starting manual backup process for all defined backups")
	err = backupops.PerformBackups(cfg)
	if err != nil {
		return fmt.Errorf("failed to perform backups: %w", err)
	}
	log.Println("Manual backup process completed successfully")

	return nil
}