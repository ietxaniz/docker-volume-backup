package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type S3Config struct {
	Endpoint        string `yaml:"endpoint"`
	Bucket          string `yaml:"bucket"`
	AccessKeyID     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Region          string `yaml:"region"`
	MaxFileSize     string `yaml:"maxFileSize"`
	BackupFolder    string `yaml:"backupFolder"`
}

type AppConfig struct {
	ScriptsFolder      string `yaml:"scriptsFolder"`
	LocalBackupFolder  string `yaml:"localBackupFolder"`
	BackupFrequency    string `yaml:"backupFrequency"`
	PublicKeyFile      string `yaml:"publicKeyFile"`
	PrivateKeyFile     string `yaml:"privateKeyFile"`
	PrivateKeyMetadata string `yaml:"privateKeyMetadata"`
}

type BackupDefinition struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`
	Containers []string `yaml:"containers"`
	Volumes    []string `yaml:"volumes"`
}

type VolumeConfig struct {
	Name       string `yaml:"name"`
	BackupName string `yaml:"backupName"`
	Compress   bool   `yaml:"compress"`
}

type AppFolders struct {
	AppStartFolder string
	ScriptsFolder  string
}

type Config struct {
	S3                S3Config           `yaml:"s3"`
	App               AppConfig          `yaml:"app"`
	Volumes           []VolumeConfig     `yaml:"volumes"`
	BackupDefinitions []BackupDefinition `yaml:"backupDefinitions"`
	AppFolders        AppFolders
}

func LoadConfiguration(configFileName string) (Config, error) {
	var config Config

	if configFileName == "" {
		configFileName = os.Getenv("S3CONFIGFILE")
		if configFileName == "" {
			return config, fmt.Errorf("no configuration file provided and S3CONFIGFILE environment variable is not set")
		}
	}

	file, err := os.Open(configFileName)
	if err != nil {
		return config, fmt.Errorf("failed to open configuration file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("failed to decode configuration file: %w", err)
	}

	appStartFolder, err := os.Getwd()
	if err != nil {
		return config, fmt.Errorf("failed to get current working directory: %w", err)
	}

	config.AppFolders = AppFolders{
		AppStartFolder: appStartFolder,
	}

	config.App.ScriptsFolder, err = getAbsPath(config.App.ScriptsFolder, appStartFolder)
	if err != nil {
		return config, fmt.Errorf("failed to get absolute path for scripts folder: %w", err)
	}
	config.AppFolders.ScriptsFolder = config.App.ScriptsFolder

	config.App.LocalBackupFolder, err = getAbsPath(config.App.LocalBackupFolder, appStartFolder)
	if err != nil {
		return config, fmt.Errorf("failed to get absolute path for local backup folder: %w", err)
	}

	config.App.PublicKeyFile, err = getAbsPath(config.App.PublicKeyFile, appStartFolder)
	if err != nil {
		return config, fmt.Errorf("failed to get absolute path for public key file: %w", err)
	}

	config.App.PrivateKeyFile, err = getAbsPath(config.App.PrivateKeyFile, appStartFolder)
	if err != nil {
		return config, fmt.Errorf("failed to get absolute path for private key file: %w", err)
	}

	config.App.PrivateKeyMetadata, err = getAbsPath(config.App.PrivateKeyMetadata, appStartFolder)
	if err != nil {
		return config, fmt.Errorf("failed to get absolute path for private key metadata: %w", err)
	}

	for i, bd := range config.BackupDefinitions {
		for j, volume := range bd.Volumes {
			if isLikelyPath(volume) {
				absPath, err := getAbsPath(volume, appStartFolder)
				if err != nil {
					return config, fmt.Errorf("failed to get absolute path for volume in backup definition: %w", err)
				}
				config.BackupDefinitions[i].Volumes[j] = absPath
			}
			// If not a likely path, assume it's a volume name and leave it as is
		}
	}

	return config, nil
}

func isLikelyPath(s string) bool {
	return strings.Contains(s, string(os.PathSeparator)) ||
		strings.Contains(s, "/") ||
		strings.Contains(s, "\\") ||
		strings.HasPrefix(s, "~") ||
		strings.HasPrefix(s, ".") ||
		filepath.IsAbs(s)
}

func getAbsPath(path, basePath string) (string, error) {
	if path == "" {
		return "", nil
	}
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Abs(filepath.Join(basePath, path))
}
