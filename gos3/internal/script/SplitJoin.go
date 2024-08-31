package script

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"gos3/internal/config"
)

func Split(folderPath, splitSize string, configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "split.sh")
	cmd := exec.Command(scriptPath,
		config.MustGetAbsPathRelativeToAppFolder(folderPath, configuration),
		splitSize)
	cmd.Dir = configuration.AppFolders.ScriptsFolder

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing split.sh: %w\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}

func Join(folderPath string, configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "join.sh")
	cmd := exec.Command(scriptPath, config.MustGetAbsPathRelativeToAppFolder(folderPath, configuration))
	cmd.Dir = configuration.AppFolders.ScriptsFolder

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing join.sh: %w\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}
