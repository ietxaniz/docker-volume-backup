package config

import "path/filepath"

func GetAbsPathRelativeToAppFolder(originalPath string, cfg Config) (string, error) {
	if filepath.IsAbs(originalPath) {
		return originalPath, nil
	}

	absPath := filepath.Join(cfg.AppFolders.AppStartFolder, originalPath)
	return filepath.Abs(absPath)
}
