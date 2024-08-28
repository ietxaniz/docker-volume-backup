package config

import (
	"fmt"
	"os"

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
	ScriptsFolder     string `yaml:"scriptsFolder"`
	LocalBackupFolder string `yaml:"localBackupFolder"`
	BackupFrequency   string `yaml:"backupFrequency"`
	PublicKeyFile     string `yaml:"publicKeyFile"`
	PrivateKeyFile    string `yaml:"privateKeyFile"`
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

type Config struct {
	S3                S3Config           `yaml:"s3"`
	App               AppConfig          `yaml:"app"`
	Volumes           []VolumeConfig     `yaml:"volumes"`
	BackupDefinitions []BackupDefinition `yaml:"backupDefinitions"`
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

	return config, nil
}
