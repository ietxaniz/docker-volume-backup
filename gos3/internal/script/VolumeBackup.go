package script

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"gos3/internal/config"
)

type VolumeBackupResult struct {
	OriginalSize     int64
	FinalSize        int64
	CompressionRatio float64
	TimeElapsed      float64
}

func VolumeBackup(volumeName, backupFileName string, compress bool, configuration config.Config) (*VolumeBackupResult, error) {
	scriptPath := filepath.Join(configuration.App.ScriptsFolder, "volume-backup.sh")

	args := []string{volumeName, backupFileName}
	if !compress {
		args = append(args, "--no-compression")
	}

	cmd := exec.Command(scriptPath, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting command: %w", err)
	}

	result := &VolumeBackupResult{}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			switch parts[0] {
			case "Original size":
				result.OriginalSize, _ = strconv.ParseInt(strings.Fields(parts[1])[0], 10, 64)
			case "Final size":
				result.FinalSize, _ = strconv.ParseInt(strings.Fields(parts[1])[0], 10, 64)
			case "Compression ratio":
				result.CompressionRatio, _ = strconv.ParseFloat(parts[1], 64)
			case "Time elapsed":
				result.TimeElapsed, _ = strconv.ParseFloat(strings.Fields(parts[1])[0], 64)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command finished with error: %w", err)
	}

	return result, nil
}
