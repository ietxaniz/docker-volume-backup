package backupops

import (
	"fmt"
	"gos3/internal/config"
	"gos3/internal/s3"
	"gos3/internal/script"
	"log"
	"os"
	"path/filepath"
)

func PerformStandardBackup(def config.BackupDefinition, cfg config.Config) error {
	log.Printf("Starting backup process for: %s", def.Name)

	err := cleanLocalBackupFolder(cfg.App.LocalBackupFolder)
	if err != nil {
		return fmt.Errorf("failed to clean local backup folder: %w", err)
	}

	log.Printf("Stopping containers: %v", def.Containers)

	err = stopContainers(def.Containers)
	if err != nil {
		return err
	}
	volumeCreationErrors := ""
	log.Printf("Containers stopped successfully: %v", def.Containers)

	for i, volumeName := range def.Volumes {
		backupFileName := generateBackupFileName(def.Name, volumeName, i)
		backupFilePath := filepath.Join(cfg.App.LocalBackupFolder, backupFileName)
		log.Printf("Creating backup for volume: %s", volumeName)

		_, err = script.VolumeBackup(volumeName, backupFilePath, true, cfg)
		if err != nil {
			volumeCreationErrors = err.Error()
			log.Printf("Backup failed for volume: %s, %s", volumeName, err.Error())
			break
		}
		log.Printf("Backup created successfully for volume: %s", volumeName)

		err = changeBackupPermissions(backupFilePath)
		if err != nil {
			log.Printf("Warning: failed to change permissions for %s: %v", backupFilePath, err)
		}
	}

	log.Printf("Starting containers: %v", def.Containers)
	err = startContainers(def.Containers)
	if err != nil {
		return err
	}
	log.Printf("Containers started successfully: %v", def.Containers)

	if len(volumeCreationErrors) != 0 {
		return fmt.Errorf("error creating volumes: %s", volumeCreationErrors)
	}

	err = encryptBackupFiles(cfg)
	if err != nil {
		return fmt.Errorf("failed to encrypt backup files: %w", err)
	}

	err = script.Split(cfg.App.LocalBackupFolder, cfg.S3.MaxFileSize, cfg)
	if err != nil {
		return fmt.Errorf("failed to split backup files: %w", err)
	}

	s3Subfolder := s3.GenerateSubfolderName(cfg.App.BackupFrequency)
	s3Path := filepath.Join(cfg.S3.BackupFolder, s3Subfolder)

	err = s3.UploadFolderToS3(cfg.App.LocalBackupFolder, s3Path, cfg)
	if err != nil {
		return fmt.Errorf("failed to upload backup to S3: %w", err)
	}

	err = cleanLocalBackupFolder(cfg.App.LocalBackupFolder)
	if err != nil {
		log.Printf("Warning: failed to clean local backup folder after upload: %v", err)
	}

	return nil
}

func cleanLocalBackupFolder(folderPath string) error {
	err := os.RemoveAll(folderPath)
	if err != nil {
		return err
	}

	return os.MkdirAll(folderPath, 0755)
}

func encryptBackupFiles(cfg config.Config) error {
	files, err := os.ReadDir(cfg.App.LocalBackupFolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(cfg.App.LocalBackupFolder, file.Name())
		encryptedFilePath := filePath + ".cpt"

		err = script.KeyEncrypt(filePath, encryptedFilePath, cfg.App.PublicKeyFile, cfg)
		if err != nil {
			return fmt.Errorf("failed to encrypt file %s: %w", filePath, err)
		}

		err = os.Remove(filePath)
		if err != nil {
			log.Printf("Warning: failed to remove original file %s after encryption: %v", filePath, err)
		}
	}

	return nil
}
