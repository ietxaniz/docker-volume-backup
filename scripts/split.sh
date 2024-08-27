#!/bin/bash

usage() {
    echo "Usage: $0 <folder_path> <split_size>"
    echo "Example: $0 /path/to/backups 100M"
    echo "Note: split_size can be in bytes (default), K, M, or G"
}

if [ "$#" -ne 2 ]; then
    usage
    exit 1
fi

FOLDER_PATH=$1
SPLIT_SIZE=$2

# Check if the folder exists
if [ ! -d "$FOLDER_PATH" ]; then
    echo "Error: The specified folder does not exist."
    exit 1
fi

# Convert SPLIT_SIZE to bytes
case ${SPLIT_SIZE: -1} in
    K|k) SPLIT_BYTES=$((${SPLIT_SIZE%[Kk]} * 1024)) ;;
    M|m) SPLIT_BYTES=$((${SPLIT_SIZE%[Mm]} * 1024 * 1024)) ;;
    G|g) SPLIT_BYTES=$((${SPLIT_SIZE%[Gg]} * 1024 * 1024 * 1024)) ;;
    *) SPLIT_BYTES=$SPLIT_SIZE ;;
esac

# Change to the specified directory
cd "$FOLDER_PATH" || exit 1

# Process each file in the folder
for file in *; do
    # Skip if it's not a file
    [ -f "$file" ] || continue
    
    # Skip if it's already a part file
    [[ $file == *.part-* ]] && continue

    # Get file size
    FILE_SIZE=$(stat -c%s "$file")

    if [ $FILE_SIZE -gt $SPLIT_BYTES ]; then
        echo "Processing $file (size: $FILE_SIZE bytes)..."
        
        # Split the file
        split -b $SPLIT_SIZE -d "$file" "${file}.part-"
        
        # Check if split was successful
        if [ $? -eq 0 ]; then
            echo "$file split into parts of $SPLIT_SIZE bytes"
            
            # Verify that at least one part file was created
            if ls "${file}.part-"* > /dev/null 2>&1; then
                # Delete the original file
                rm "$file"
                echo "Original file deleted"
            else
                echo "Error: No part files created. Original file preserved."
            fi
        else
            echo "Error occurred during file splitting. Original file preserved."
        fi
        echo "------------------------"
    else
        echo "$file size ($FILE_SIZE bytes) does not exceed split size ($SPLIT_BYTES bytes). Skipping."
    fi
done

echo "All splitting operations completed."