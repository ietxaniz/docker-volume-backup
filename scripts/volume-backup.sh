#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <volume_name> <backup_file_name>"
    exit 1
fi

VOLUME_NAME=$1
BACKUP_FILE=$2

touch "$BACKUP_FILE"

docker run --rm -v $VOLUME_NAME:/volume -v $BACKUP_FILE:/backup.tar.gz alpine tar -czpf /backup.tar.gz -C /volume ./

echo "Backup of volume $VOLUME_NAME created as $BACKUP_FILE"
