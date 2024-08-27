package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type S3Config struct {
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	Bucket          string `json:"bucket" yaml:"bucket"`
	AccessKeyID     string `json:"accessKeyId" yaml:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret"`
	Region          string `json:"region" yaml:"region"`
}

type AppConfig struct {
	ScriptsFolder string `json:"scriptsFolder" yaml:"scriptsFolder"`
}

type Config struct {
	S3  S3Config  `json:"s3" yaml:"s3"`
	App AppConfig `json:"app" yaml:"app"`
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
