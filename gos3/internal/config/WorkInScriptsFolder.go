package config

import (
	"fmt"
	"os"
)

func WorkInScriptsFolder(cfg Config) (err error) {

	err = os.Chdir(cfg.AppFolders.ScriptsFolder)
	if err != nil {
		return fmt.Errorf("failed to change to scripts folder: %w", err)
	}

	return nil
}
