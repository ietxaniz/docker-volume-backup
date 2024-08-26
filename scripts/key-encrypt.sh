#!/bin/bash

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <input_file> <output_encrypted_file> <public_key_file>"
    exit 1
fi

INPUT_FILE=$1
OUTPUT_FILE=$2
PUBLIC_KEY=$3

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file not found."
    exit 1
fi

# Check if public key file exists
if [ ! -f "$PUBLIC_KEY" ]; then
    echo "Error: Public key file not found."
    exit 1
fi

# Generate a random password
RANDOM_PASS=$(openssl rand -base64 32)

# Encrypt the file with the random password
openssl enc -aes-256-cbc -salt -in "$INPUT_FILE" -out "$OUTPUT_FILE" -k "$RANDOM_PASS"

# Encrypt the random password with the public key
echo "$RANDOM_PASS" | openssl pkeyutl -encrypt -pubin -inkey "$PUBLIC_KEY" -out "${OUTPUT_FILE}.pass"

echo "File encrypted successfully: $OUTPUT_FILE"
echo "Encrypted password: ${OUTPUT_FILE}.pass"