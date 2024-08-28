package backupops

import (
	"fmt"
	"gos3/internal/config"
	"gos3/internal/script"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func stopContainers(containers []string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(containers))

	for _, container := range containers {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			cmd := exec.Command("docker", "stop", c)
			if err := cmd.Run(); err != nil {
				errChan <- fmt.Errorf("failed to stop container %s: %w", c, err)
			}
		}(container)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func startContainers(containers []string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(containers))

	for _, container := range containers {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			cmd := exec.Command("docker", "start", c)
			if err := cmd.Run(); err != nil {
				errChan <- fmt.Errorf("failed to start container %s: %w", c, err)
			}
		}(container)
	}

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to start one or more containers: %v", errors)
	}

	return nil
}

func generateBackupFileName(backupName string, volumeName string, index int) string {
	if isVolumePath(volumeName) {
		return fmt.Sprintf("%s-%d.tar.gz", backupName, index)
	}
	return fmt.Sprintf("%s-%s.tar.gz", backupName, volumeName)
}

func encryptBackup(inputFile, outputFile string, cfg config.Config) error {
	err := script.KeyEncrypt(inputFile, outputFile, cfg.App.PublicKeyFile, cfg)
	if err != nil {
		return fmt.Errorf("failed to encrypt backup: %w", err)
	}

	log.Printf("Encrypted file: %s to %s", inputFile, outputFile)
	return nil
}

func isVolumePath(volumeName string) bool {
	// Check if the volume name is a file path
	return filepath.IsAbs(volumeName)
}

func changeBackupPermissions(backupFilePath string, cfg config.Config) error {
    cmd := exec.Command("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/backup", filepath.Dir(backupFilePath)),
        "alpine", "chown", "-R", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()), "/backup")
    return cmd.Run()
}
