#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <volume_name> <backup_file_name>"
    exit 1
fi

VOLUME_NAME=$1
BACKUP_FILE=$2

if [ ! -f "$BACKUP_FILE" ]; then
    echo "Error: Backup file $BACKUP_FILE does not exist."
    exit 1
fi

docker volume inspect $VOLUME_NAME > /dev/null 2>&1 || docker volume create $VOLUME_NAME

docker run --rm -v $VOLUME_NAME:/volume -v $BACKUP_FILE:/backup.tar.gz alpine sh -c "rm -rf /volume/* /volume/..?* /volume/.[!.]* ; tar -xzpf /backup.tar.gz -C /volume"

echo "Restore of $BACKUP_FILE to volume $VOLUME_NAME completed"
