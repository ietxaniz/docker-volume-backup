package script

import (
	"fmt"
	"gos3/internal/config"
	"io"
	"os/exec"
	"path/filepath"
)

func FileEncrypt(inputFile, outputFile string, password string, configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "file-encrypt.sh")
	cmd := exec.Command(scriptPath,
		config.MustGetAbsPathRelativeToAppFolder(inputFile, configuration),
		config.MustGetAbsPathRelativeToAppFolder(outputFile, configuration))
	cmd.Dir = configuration.AppFolders.ScriptsFolder

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating stdin pipe: %w", err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, password+"\n")
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing file-encrypt.sh: %w\nOutput: %s", err, string(output))
	}

	return nil
}
