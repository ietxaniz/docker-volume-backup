# Docker Volume Backup, Restore, and File Encryption Scripts

This repository contains scripts for backing up and restoring Docker volumes, as well as encrypting and decrypting files, including the use of asymmetric encryption for secure automated backups.

## Scripts

Main scripts provided in the `scripts` folder:

1. `volume-backup.sh`: Creates backups of Docker volumes.
2. `volume-restore.sh`: Restores Docker volumes from backups.
3. `file-encrypt.sh`: Encrypts a file using symmetric encryption.
4. `file-decrypt.sh`: Decrypts an encrypted file.
5. `key-generate.sh`: Generates a public and private key pair for asymmetric encryption.
6. `key-encrypt.sh`: Encrypts data using the public key, typically for secure backup storage.
7. `key-decrypt.sh`: Decrypts data encrypted with the public key using the private key.

Example usage scripts:

- `sample-create.sh`: Demonstrates how to use the volume backup script.
- `sample-restore.sh`: Demonstrates how to use the volume restore script.
- `sample-fileencrypt.sh`: Demonstrates how to use the file encryption script.
- `sample-filedecrypt.sh`: Demonstrates how to use the file decryption script.
- `sample-filedecrypt-fail.sh`: Demonstrates decryption with an incorrect password.

## Use Cases

These scripts can be useful in various scenarios:

1. **Transfer Data Between Volumes**:
   - On the same machine or between different machines by transferring the backup file.

2. **Create Encrypted Backups**:
   - For point-in-time recovery, disaster recovery, or before significant system changes.

3. **Secure Data Migration**:
   - When upgrading applications or moving to different cloud providers or infrastructure.

4. **Cloning Environments with Data Protection**:
   - For creating secure development or staging environments from production data.

5. **Secure Archiving**:
   - For long-term storage of historical data in compliance with data retention policies.

6. **Testing and Debugging with Sensitive Data**:
   - To create reproducible environments for testing and capture system states for debugging.

7. **Scaling and Load Balancing with Data Security**:
   - Quickly provisioning new instances with existing data.

8. **Enhanced Encryption Scenarios**:
   - Encrypt backup files for secure storage or transfer and protect sensitive files with strong encryption.

## Usage

For detailed usage instructions, refer to the comments in the scripts and the example scripts provided.

### File Encryption and Decryption

To encrypt a file:
```bash
./file-encrypt.sh <input_file> <output_encrypted_file> <password>
```

To decrypt a file:
```bash
./file-decrypt.sh <input_encrypted_file> <output_decrypted_file> <password>
```

Example:
```bash
./file-encrypt.sh /mnt/data/volume-backups/mydata.tar.gz /mnt/data/volume-backups/mydata.cpt MyStrongPassword
./file-decrypt.sh /mnt/data/volume-backups/mydata.cpt /mnt/data/volume-backups/mydata-restored.tar.gz MyStrongPassword
```

### Asymmetric Encryption for Automated Backups

Using asymmetric encryption, backup files can be securely encrypted using a public key, allowing only the holder of the private key to decrypt them. This is particularly useful in automated systems where manual password entry is not feasible.

## Security Considerations

- Use strong, unique passwords for file encryption and robust key management for asymmetric encryption.
- Store encrypted files, passwords, and keys separately and securely.
- The scripts employ robust encryption methods, but their security also depends on the strength of the chosen passwords and keys, and the security of the environment where they are used.
- Always verify the integrity of decrypted files before use.
- Implement regular security audits and updates to ensure the encryption tools and methods remain secure against new vulnerabilities.
