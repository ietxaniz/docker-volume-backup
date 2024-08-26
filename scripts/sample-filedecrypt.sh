#!/bin/bash
# This script is for demonstration purposes only. In production the right part of the script should be used and the password should be added manualy.

# Decrypt the file using a hard-coded password.
echo "My678-strong-pass" | ./file-decrypt.sh /mnt/data/volume-backups/mmx-dbgen-data.cpt /mnt/data/volume-backups/mmx-dbgen-data-2.tar.gz