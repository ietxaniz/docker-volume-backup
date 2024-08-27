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

# Check if private key file exists
if [ ! -f "$PRIVATE_KEY" ]; then
    echo "Error: Private key file not found."
    exit 1
fi

PASS_FILE="${INPUT_FILE}.pass"

# Check if pass file exists
if [ ! -f "$PASS_FILE" ]; then
    echo "Error: Encrypted password file not found: $PASS_FILE"
    exit 1
fi

# Decrypt the password using the private key
RANDOM_PASS=$(openssl pkeyutl -decrypt -inkey "$PRIVATE_KEY" -in "$PASS_FILE")

if [ -z "$RANDOM_PASS" ]; then
    echo "Error: Failed to decrypt the password."
    exit 1
fi

# Decrypt the input file using the decrypted password with PBKDF2
openssl enc -d -aes-256-cbc -pbkdf2 -iter 100000 -in "$INPUT_FILE" -out "$OUTPUT_FILE" -k "$RANDOM_PASS"

if [ $? -eq 0 ]; then
    echo "File decrypted successfully: $OUTPUT_FILE"
else
    echo "Error: Failed to decrypt the file."
    exit 1
fi
