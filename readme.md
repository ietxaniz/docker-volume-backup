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
8. `key-decrypt2.sh`: Decrypts data when the private key itself is encrypted.

Example usage scripts:

- `sample-create.sh`: Demonstrates how to use the volume backup script.
- `sample-restore.sh`: Demonstrates how to use the volume restore script.
- `sample-fileencrypt.sh`: Demonstrates how to use the file encryption script.
- `sample-filedecrypt.sh`: Demonstrates how to use the file decryption script.
- `sample-filedecrypt-fail.sh`: Demonstrates decryption with an incorrect password.
- `sample-asym-fileencrypt.sh`: Demonstrates asymmetric encryption of a file.
- `sample-asym-filedecrypt.sh`: Demonstrates asymmetric decryption of a file with an unencrypted private key.
- `sample-asym-filedecrypt2.sh`: Demonstrates asymmetric decryption of a file with an encrypted private key.

## Recommended Workflow

For secure backup and restoration of Docker volumes:

1. Create a public-private key pair:
   ```
   ./key-generate.sh
   ```

2. Encrypt the private key using symmetric encryption:
   ```
   ./file-encrypt.sh private_key.pem private_key.pem.cpt
   ```
   Store the encrypted private key (`private_key.pem.cpt`) and its associated metadata file securely.

3. Back up the Docker volume:
   ```
   ./volume-backup.sh <volume_name> <backup_file.tar.gz>
   ```

4. Encrypt the backup data using the public key:
   ```
   ./key-encrypt.sh <backup_file.tar.gz> <encrypted_backup.cpt> public_key.pem
   ```

5. Store the encrypted backup data securely.

For restoration:

6. Retrieve the encrypted backup data.

7. If the private key was encrypted:
   ```
   ./key-decrypt2.sh <encrypted_backup.cpt> <decrypted_backup.tar.gz> private_key.pem.cpt
   ```
   This script will prompt for the password to decrypt the private key, then use it to decrypt the backup.

   If the private key was not encrypted:
   ```
   ./key-decrypt.sh <encrypted_backup.cpt> <decrypted_backup.tar.gz> private_key.pem
   ```

8. Restore the Docker volume:
   ```
   ./volume-restore.sh <volume_name> <decrypted_backup.tar.gz>
   ```

## Use Cases

These scripts can be useful in various scenarios:

1. **Transfer Data Between Volumes**: On the same machine or between different machines by transferring the encrypted backup file.

2. **Create Encrypted Backups**: For point-in-time recovery, disaster recovery, or before significant system changes.

3. **Secure Data Migration**: When upgrading applications or moving to different cloud providers or infrastructure.

4. **Cloning Environments with Data Protection**: For creating secure development or staging environments from production data.

5. **Secure Archiving**: For long-term storage of historical data in compliance with data retention policies.

6. **Testing and Debugging with Sensitive Data**: To create reproducible environments for testing and capture system states for debugging.

7. **Scaling and Load Balancing with Data Security**: Quickly provisioning new instances with existing data.

8. **Enhanced Encryption Scenarios**: Encrypt backup files for secure storage or transfer and protect sensitive files with strong encryption.

## Usage

### Docker Volume Backup and Restore

To backup a Docker volume:
```bash
./volume-backup.sh <volume_name> <backup_file.tar.gz>
```

To restore a Docker volume:
```bash
./volume-restore.sh <volume_name> <backup_file.tar.gz>
```

Example:
```bash
./volume-backup.sh my_volume /path/to/backup.tar.gz
./volume-restore.sh my_volume /path/to/backup.tar.gz
```

### Symmetric File Encryption and Decryption

To encrypt a file:
```bash
./file-encrypt.sh <input_file> <output_encrypted_file>
```

To decrypt a file:
```bash
./file-decrypt.sh <input_encrypted_file> <output_decrypted_file>
```

Example:
```bash
./file-encrypt.sh /mnt/data/myfile.txt /mnt/data/myfile.cpt
./file-decrypt.sh /mnt/data/myfile.cpt /mnt/data/myfile_decrypted.txt
```

### Asymmetric Key Generation

To generate a public-private key pair:
```bash
./key-generate.sh
```

This will create two files: `private_key.pem` and `public_key.pem`.

### Asymmetric Encryption and Decryption

To encrypt a file using the public key:
```bash
./key-encrypt.sh <input_file> <output_encrypted_file> <public_key_file>
```

To decrypt a file using the unencrypted private key:
```bash
./key-decrypt.sh <input_encrypted_file> <output_decrypted_file> <private_key_file>
```

To decrypt a file using an encrypted private key:
```bash
./key-decrypt2.sh <input_encrypted_file> <output_decrypted_file> <encrypted_private_key_file>
```

Example:
```bash
./key-encrypt.sh /mnt/data/backup.tar.gz /mnt/data/backup.cpt public_key.pem
./key-decrypt.sh /mnt/data/backup.cpt /mnt/data/restored_backup.tar.gz private_key.pem
./key-decrypt2.sh /mnt/data/backup.cpt /mnt/data/restored_backup.tar.gz private_key.pem.cpt
```

## Sample Scripts Explained

- `sample-asym-fileencrypt.sh`: This script demonstrates how to encrypt a Docker volume backup using asymmetric encryption. It uses the public key to encrypt the backup, allowing for secure storage or transfer.

- `sample-asym-filedecrypt.sh`: This script shows how to decrypt an asymmetrically encrypted backup using an unencrypted private key. This is useful when the private key is already secure and doesn't need an additional layer of encryption.

- `sample-asym-filedecrypt2.sh`: This script demonstrates the decryption of an asymmetrically encrypted backup using an encrypted private key. This provides an extra layer of security for the private key itself.

These sample scripts serve as practical examples of how to use the encryption and decryption scripts in various scenarios, showcasing the flexibility and security of the system.

