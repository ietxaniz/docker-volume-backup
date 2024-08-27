#!/bin/bash

usage() {
    echo "Usage: $0 <volume_name> <backup_file_name> [options]"
    echo "Options:"
    echo "  -n, --no-compression     Create backup without compression"
    echo "  -h, --help               Display this help message"
}

if [ "$#" -lt 2 ]; then
    usage
    exit 1
fi

VOLUME_NAME=$1
BACKUP_FILE=$2
shift 2

COMPRESS=true

while [ "$#" -gt 0 ]; do
    case "$1" in
        -n|--no-compression) COMPRESS=false ;;
        -h|--help) usage; exit 0 ;;
        *) echo "Unknown option: $1"; usage; exit 1 ;;
    esac
    shift
done

start_time=$(date +%s.%N)

if [ "$COMPRESS" = true ]; then
    # Compressed backup with size reporting
    orig_size=$(docker run --rm -v $VOLUME_NAME:/volume alpine du -sb /volume | cut -f1)
    docker run --rm -v $VOLUME_NAME:/volume -v $(dirname $BACKUP_FILE):/backup alpine tar -czpf /backup/$(basename $BACKUP_FILE) -C /volume .
    final_size=$(stat -c%s "$BACKUP_FILE")
    compression_ratio=$(echo "scale=2; $final_size / $orig_size" | bc)
else
    # Uncompressed backup
    docker run --rm -v $VOLUME_NAME:/volume -v $(dirname $BACKUP_FILE):/backup alpine tar -cpf /backup/$(basename $BACKUP_FILE) -C /volume .
    final_size=$(stat -c%s "$BACKUP_FILE")
    orig_size=$final_size
    compression_ratio=1
fi

end_time=$(date +%s.%N)
elapsed=$(echo "$end_time - $start_time" | bc)
elapsed=$(printf "%.6f" $elapsed)

echo "Backup of volume $VOLUME_NAME created as $BACKUP_FILE"
echo "Original size: $orig_size bytes"
echo "Final size: $final_size bytes"
echo "Compression ratio: $compression_ratio"
echo "Time elapsed: $elapsed seconds"