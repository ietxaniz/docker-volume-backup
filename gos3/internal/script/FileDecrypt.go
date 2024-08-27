package script

import (
	"fmt"
	"gos3/internal/config"
	"io"
	"os/exec"
	"path/filepath"
)

func FileDecrypt(inputFile, outputFile string, password string, configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "file-decrypt.sh")
	cmd := exec.Command(scriptPath, inputFile, outputFile)

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
		return fmt.Errorf("error executing file-decrypt.sh: %w\nOutput: %s", err, string(output))
	}

	return nil
}
