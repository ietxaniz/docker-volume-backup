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

	err := stopContainers(def.Containers)
	if err != nil {
		return err
	}
	volumeCreationErrors := ""

	for i, volumeName := range def.Volumes {
		backupFileName := generateBackupFileName(def.Name, volumeName, i)
		backupFilePath := filepath.Join(cfg.App.LocalBackupFolder, backupFileName)

		_, err = script.VolumeBackup(volumeName, backupFilePath, true, cfg)
		if err != nil {
			volumeCreationErrors = err.Error()
			break
		}
	}

	err = startContainers(def.Containers)
	if err != nil {
		return err
	}

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
		err = s3.UploadToS3(encryptedFilePath, s3PathCpt, cfg.S3)
		if err != nil {
			return fmt.Errorf("failed to upload encrypted backup for volume %s to S3: %w", volumeName, err)
		}

		// Upload .pass file
		s3PathPass := filepath.Join(cfg.S3.BackupFolder, s3Subfolder, backupFileName+".cpt.pass")
		err = s3.UploadToS3(passFilePath, s3PathPass, cfg.S3)
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
