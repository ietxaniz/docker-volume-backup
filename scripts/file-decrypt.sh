#!/bin/bash

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <input_file> <output_file> <password>"
    exit 1
fi

INPUT_FILE=$1
OUTPUT_FILE=$2
PASSWORD=$3
METADATA_FILE="${INPUT_FILE}.metadata"

if [ ! -f "$METADATA_FILE" ]; then
    echo "Metadata file not found: $METADATA_FILE"
    exit 1
fi

# Read metadata
SALT=$(grep "Salt:" "$METADATA_FILE" | cut -d' ' -f2)
ITERATIONS=$(grep "Iterations:" "$METADATA_FILE" | cut -d' ' -f2)
IV=$(grep "IV:" "$METADATA_FILE" | cut -d' ' -f2)

# Derive key using metadata
KEY_INFO=$(./derive-key.sh "$PASSWORD" "$SALT" "$ITERATIONS")
KEY=$(echo "$KEY_INFO" | grep "Key:" | cut -d' ' -f2)

# Decrypt the file
openssl enc -d -aes-256-cbc -in "$INPUT_FILE" -out "$OUTPUT_FILE" -K "$KEY" -iv "$IV"

echo "File decrypted successfully."
