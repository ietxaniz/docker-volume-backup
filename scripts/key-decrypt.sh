#!/bin/bash

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <input_encrypted_file> <output_decrypted_file> <private_key_file>"
    exit 1
fi

INPUT_FILE=$1
OUTPUT_FILE=$2
PRIVATE_KEY=$3

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file not found."
    exit 1
fi

# Check if encrypted password file exists
if [ ! -f "${INPUT_FILE}.pass" ]; then
    echo "Error: Encrypted password file not found."
    exit 1
fi

# Check if private key file exists
if [ ! -f "$PRIVATE_KEY" ]; then
    echo "Error: Private key file not found."
    exit 1
fi

# Decrypt the random password
RANDOM_PASS=$(openssl pkeyutl -decrypt -inkey "$PRIVATE_KEY" -in "${INPUT_FILE}.pass")

# Decrypt the file
openssl enc -d -aes-256-cbc -salt -in "$INPUT_FILE" -out "$OUTPUT_FILE" -k "$RANDOM_PASS"

echo "File decrypted successfully: $OUTPUT_FILE"