#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <input_file> <output_file>"
    exit 1
fi

INPUT_FILE=$1
OUTPUT_FILE=$2

echo -n "Enter encryption password: "
read -s PASSWORD
echo

# Generate salt and derive key
KEY_INFO=$(./derive-key.sh "$PASSWORD")
SALT=$(echo "$KEY_INFO" | grep "Salt:" | cut -d' ' -f2)
ITERATIONS=$(echo "$KEY_INFO" | grep "Iterations:" | cut -d' ' -f2)
KEY=$(echo "$KEY_INFO" | grep "Key:" | cut -d' ' -f2)

# Generate a random IV for AES-256-CBC
IV=$(openssl rand -hex 16)

# Encrypt the file using the derived key and IV
openssl enc -aes-256-cbc -in "$INPUT_FILE" -out "$OUTPUT_FILE" -K "$KEY" -iv "$IV"

# Store metadata with IV
echo "Salt: $SALT" > "${OUTPUT_FILE}.metadata"
echo "Iterations: $ITERATIONS" >> "${OUTPUT_FILE}.metadata"
echo "IV: $IV" >> "${OUTPUT_FILE}.metadata"

echo "File encrypted successfully."
