#!/bin/bash
# This script is for demonstration purposes only. In production the right part of the script should be used and the password should be added manualy.

# Attempt to decrypt the file using an incorrect password to demonstrate a decryption failure.
echo "My678-strong-bass" | ./file-decrypt.sh /mnt/data/volume-backups/mmx-dbgen-data.cpt /mnt/data/volume-backups/mmx-dbgen-data-3.tar.gz
