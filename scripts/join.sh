#!/bin/bash

usage() {
    echo "Usage: $0 <folder_path>"
    echo "Example: $0 /path/to/backup/parts"
    echo "Note: This script will join all sets of *.part-* files in the specified folder and its subfolders."
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

# Find all split folders
SPLIT_FOLDERS=$(find . -type d -name "*-split_parts")

for SPLIT_FOLDER in $SPLIT_FOLDERS; do
    echo "Processing folder: $SPLIT_FOLDER"
    
    # Extract the original file name
    ORIGINAL_FILE=$(basename "$SPLIT_FOLDER" "-split_parts")
    
    # Find all part files
    PART_FILES=$(find "$SPLIT_FOLDER" -name "${ORIGINAL_FILE}.part-*" | sort -V)
    
    if [ -z "$PART_FILES" ]; then
        echo "Error: No part files found in $SPLIT_FOLDER."
        continue
    fi

    OUTPUT_FILE="$ORIGINAL_FILE"
    
    echo "Joining parts to create: $OUTPUT_FILE"

    # Join the parts
    cat $PART_FILES > "$OUTPUT_FILE"

    # Check if joining was successful
    if [ $? -eq 0 ]; then
        echo "Successfully joined parts into $OUTPUT_FILE"
        
        # Verify the joined file exists and has a non-zero size
        if [ -s "$OUTPUT_FILE" ]; then
            echo "Deleting split folder..."
            rm -r "$SPLIT_FOLDER"
            echo "Split folder deleted."
        else
            echo "Error: Joined file is empty. Split folder preserved."
            rm "$OUTPUT_FILE"  # Remove the empty file
            continue
        fi
    else
        echo "Error occurred during joining. Split folder preserved."
        continue
    fi

    # Display the size of the joined file
    JOINED_SIZE=$(stat -c%s "$OUTPUT_FILE")
    echo "Size of joined file: $JOINED_SIZE bytes"
    echo "Joined file created at: $FOLDER_PATH/$OUTPUT_FILE"
    echo "------------------------"
done

echo "All joining operations completed."