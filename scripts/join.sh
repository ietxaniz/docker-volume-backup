#!/bin/bash

usage() {
    echo "Usage: $0 <folder_path>"
    echo "Example: $0 /path/to/backup/parts"
    echo "Note: This script will join all sets of *.part-* files in the specified folder."
}

if [ "$#" -ne 1 ]; then
    usage
    exit 1
fi

FOLDER_PATH=$1

# Check if the folder exists
if [ ! -d "$FOLDER_PATH" ]; then
    echo "Error: The specified folder does not exist."
    exit 1
fi

# Change to the specified directory
cd "$FOLDER_PATH" || exit 1

# Find all unique base names
BASE_NAMES=$(ls *.part-* 2>/dev/null | sed 's/\.part-[0-9]*$//' | sort -u)

if [ -z "$BASE_NAMES" ]; then
    echo "Error: No part files found in the specified folder."
    exit 1
fi

for BASE_NAME in $BASE_NAMES; do
    echo "Processing $BASE_NAME..."
    
    PART_FILES=$(ls "${BASE_NAME}".part-* 2>/dev/null | sort)
    
    if [ -z "$PART_FILES" ]; then
        echo "Error: No part files found for $BASE_NAME."
        continue
    fi

    OUTPUT_FILE="$BASE_NAME"
    
    echo "Joining parts to create: $OUTPUT_FILE"

    # Join the parts
    cat $PART_FILES > "$OUTPUT_FILE"

    # Check if joining was successful
    if [ $? -eq 0 ]; then
        echo "Successfully joined parts into $OUTPUT_FILE"
        
        # Verify the joined file exists and has a non-zero size
        if [ -s "$OUTPUT_FILE" ]; then
            echo "Deleting part files..."
            rm "${BASE_NAME}".part-*
            echo "Part files deleted."
        else
            echo "Error: Joined file is empty. Part files preserved."
            rm "$OUTPUT_FILE"  # Remove the empty file
            continue
        fi
    else
        echo "Error occurred during joining. Part files preserved."
        continue
    fi

    # Display the size of the joined file
    JOINED_SIZE=$(stat -c%s "$OUTPUT_FILE")
    echo "Size of joined file: $JOINED_SIZE bytes"
    echo "Joined file created at: $FOLDER_PATH/$OUTPUT_FILE"
    echo "------------------------"
done

echo "All joining operations completed."