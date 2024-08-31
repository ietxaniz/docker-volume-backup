package script

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"

	"gos3/internal/config"
)

func VolumeRestore(volumeName, backupFileName string, configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "volume-restore.sh")

	cmd := exec.Command(scriptPath, volumeName, config.MustGetAbsPathRelativeToAppFolder(backupFileName, configuration))
	cmd.Dir = configuration.AppFolders.ScriptsFolder

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command finished with error: %w", err)
	}

	return nil
}
