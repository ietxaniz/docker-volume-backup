package script

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"gos3/internal/config"
)

type DeriveKeyResult struct {
	Salt       string
	Iterations string
	Key        string
	TimeTaken  string
}

func generateRandomSalt() (string, error) {
	bytes := make([]byte, 8) // 8 bytes for a 16-character hex string
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func DeriveKey(password string, salt string, iterations int, configuration config.Config) (*DeriveKeyResult, error) {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "derive-key.sh")

	args := []string{password}

	// If iterations are provided but salt is not, generate a random salt
	if iterations > 0 && salt == "" {
		var err error
		salt, err = generateRandomSalt()
		if err != nil {
			return nil, fmt.Errorf("failed to generate random salt: %w", err)
		}
	}

	// Add salt and iterations to args if they are provided
	if salt != "" {
		args = append(args, salt)
		if iterations > 0 {
			args = append(args, fmt.Sprintf("%d", iterations))
		}
	}

	cmd := exec.Command(scriptPath, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting command: %w", err)
	}

	result := &DeriveKeyResult{}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			switch parts[0] {
			case "Salt":
				result.Salt = parts[1]
			case "Iterations":
				result.Iterations = parts[1]
			case "Key":
				result.Key = parts[1]
			case "Time taken":
				result.TimeTaken = parts[1]
			}
		}
	}

	errOutput, _ := bufio.NewReader(stderr).ReadString('\n')
	if errOutput != "" {
		return nil, fmt.Errorf("script error: %s", strings.TrimSpace(errOutput))
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command finished with error: %w", err)
	}

	if result.Key == "" {
		return nil, fmt.Errorf("failed to derive key")
	}

	return result, nil
}

func PrintDeriveKeyResult(result *DeriveKeyResult) {
	fmt.Println("Derive Key Result:")
	fmt.Printf("  Salt: %s\n", result.Salt)
	fmt.Printf("  Iterations: %s\n", result.Iterations)
	fmt.Printf("  Key: %s\n", result.Key)
	fmt.Printf("  Time taken: %s\n", result.TimeTaken)
}
