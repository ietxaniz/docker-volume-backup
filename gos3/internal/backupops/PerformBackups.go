package backupops

import (
	"gos3/internal/config"
	"log"
)

func PerformBackups(cfg config.Config) error {
	for _, backupDef := range cfg.BackupDefinitions {
		log.Printf("Starting backup process for: %s", backupDef.Name)

		var err error
		switch backupDef.Type {
		case "standard":
			err = PerformStandardBackup(backupDef, cfg)
		default:
			log.Printf("Unknown backup type: %s for backup: %s", backupDef.Type, backupDef.Name)
			continue
		}

		if err != nil {
			log.Printf("Error performing backup %s: %v", backupDef.Name, err)
			return err
		}

		log.Printf("Completed backup process for: %s", backupDef.Name)
	}

	return nil
}
