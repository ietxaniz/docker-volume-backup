package script

import (
	"fmt"
	"gos3/internal/config"
	"log"
	"os/exec"
	"path/filepath"
)

func KeyDecrypt2(inputFile, outputFile, encryptedPrivateKeyFile string, privateKeyPassword string, configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "key-decrypt2.sh")
	cmd := exec.Command(scriptPath, inputFile, outputFile, encryptedPrivateKeyFile)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating stdin pipe: %w", err)
	}

	go func() {
		defer stdin.Close()
		stdin.Write([]byte(privateKeyPassword + "\n"))
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing key-decrypt2.sh: %w\nOutput: %s", err, string(output))
	}

	log.Printf("KeyDecrypt2 output: %s", string(output))
	return nil
}
