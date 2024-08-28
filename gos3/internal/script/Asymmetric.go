package script

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"gos3/internal/config"
)

func KeyGenerate(configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "key-generate.sh")
	cmd := exec.Command(scriptPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing key-generate.sh: %w\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}

func KeyEncrypt(inputFile, outputFile, publicKeyFile string, configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "key-encrypt.sh")
	cmd := exec.Command(scriptPath, inputFile, outputFile, publicKeyFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing key-encrypt.sh: %w\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}

func KeyDecrypt(inputFile, outputFile, privateKeyFile string, configuration config.Config) error {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "key-decrypt.sh")
	cmd := exec.Command(scriptPath, inputFile, outputFile, privateKeyFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing key-decrypt.sh: %w\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}

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
