package script

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"gos3/internal/config"
)

func KeyGenerate(configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "key-generate.sh")
	cmd := exec.Command(scriptPath)
	cmd.Dir = configuration.AppFolders.ScriptsFolder

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing key-generate.sh: %w\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}

func KeyEncrypt(inputFile, outputFile, publicKeyFile string, configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "key-encrypt.sh")
	cmd := exec.Command(scriptPath,
		config.MustGetAbsPathRelativeToAppFolder(inputFile, configuration),
		config.MustGetAbsPathRelativeToAppFolder(outputFile, configuration),
		config.MustGetAbsPathRelativeToAppFolder(publicKeyFile, configuration))
	cmd.Dir = configuration.AppFolders.ScriptsFolder

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing key-encrypt.sh: %w\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}

func KeyDecrypt(inputFile, outputFile, privateKeyFile string, configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "key-decrypt.sh")
	cmd := exec.Command(scriptPath,
		config.MustGetAbsPathRelativeToAppFolder(inputFile, configuration),
		config.MustGetAbsPathRelativeToAppFolder(outputFile, configuration),
		config.MustGetAbsPathRelativeToAppFolder(privateKeyFile, configuration))
	cmd.Dir = configuration.AppFolders.ScriptsFolder

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing key-decrypt.sh: %w\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}
