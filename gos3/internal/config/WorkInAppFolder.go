package config

import (
	"fmt"
	"os"
)

func WorkInAppFolder(cfg Config) (err error) {

	err = os.Chdir(cfg.AppFolders.AppStartFolder)
	if err != nil {
		return fmt.Errorf("failed to change to app folder: %w", err)
	}

	return nil
}
