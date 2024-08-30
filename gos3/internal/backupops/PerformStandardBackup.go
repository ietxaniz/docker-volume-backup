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
	log.Printf("Stopping containers: %v", def.Containers)

	err := stopContainers(def.Containers)
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

	for i, volumeName := range def.Volumes {
		backupFileName := generateBackupFileName(def.Name, volumeName, i)
		backupFilePath := filepath.Join(cfg.App.LocalBackupFolder, backupFileName)

		encryptedFilePath := backupFilePath + ".cpt"
		passFilePath := encryptedFilePath + ".pass"

		err = encryptBackup(backupFilePath, encryptedFilePath, cfg)
		if err != nil {
			return fmt.Errorf("failed to encrypt backup for volume %s: %w", volumeName, err)
		}

		s3Subfolder := s3.GenerateSubfolderName(cfg.App.BackupFrequency)

		// Upload .cpt file
		s3PathCpt := filepath.Join(cfg.S3.BackupFolder, s3Subfolder, backupFileName+".cpt")
		err = s3.UploadToS3(encryptedFilePath, s3PathCpt, cfg)
		if err != nil {
			return fmt.Errorf("failed to upload encrypted backup for volume %s to S3: %w", volumeName, err)
		}

		// Upload .pass file
		s3PathPass := filepath.Join(cfg.S3.BackupFolder, s3Subfolder, backupFileName+".cpt.pass")
		err = s3.UploadToS3(passFilePath, s3PathPass, cfg)
		if err != nil {
			return fmt.Errorf("failed to upload pass file for volume %s to S3: %w", volumeName, err)
		}

		// Clean up local files
		for _, filePath := range []string{backupFilePath, encryptedFilePath, passFilePath} {
			err = os.Remove(filePath)
			if err != nil {
				log.Printf("Warning: failed to remove local file %s: %v", filePath, err)
			}
		}
	}

	return nil
}
