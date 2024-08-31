package script

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"gos3/internal/config"
)

func KeyDecrypt2WithPass(inputFile, outputFile, encryptedPrivateKeyFile, privateKeyPassword string, configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "key-decrypt2-withpass.sh")
	cmd := exec.Command(scriptPath,
		config.MustGetAbsPathRelativeToAppFolder(inputFile, configuration),
		config.MustGetAbsPathRelativeToAppFolder(outputFile, configuration),
		config.MustGetAbsPathRelativeToAppFolder(encryptedPrivateKeyFile, configuration),
		privateKeyPassword)
	cmd.Dir = configuration.AppFolders.ScriptsFolder

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing key-decrypt2-withpass.sh: %w\nOutput: %s", err, string(output))
	}

	log.Printf("KeyDecrypt2WithPass output: %s", string(output))
	return nil
}
