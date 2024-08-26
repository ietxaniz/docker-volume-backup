#!/bin/bash
docker volume create mmx-dbgen-data-2
./volume-restore.sh mmx-dbgen-data-2 /mnt/data/volume-backups/mmx-dbgen-data.tar.gz