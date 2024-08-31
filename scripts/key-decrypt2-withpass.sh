#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status.

if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <input_encrypted_file> <output_decrypted_file> <encrypted_private_key_file> <private_key_password>"
    exit 1
fi

INPUT_ENCRYPTED_FILE=$1
OUTPUT_DECRYPTED_FILE=$2
ENCRYPTED_PRIVATE_KEY_FILE=$3
PRIVATE_KEY_PASSWORD=$4

# Check if required files exist
for file in "$INPUT_ENCRYPTED_FILE" "$ENCRYPTED_PRIVATE_KEY_FILE"; do
    if [ ! -f "$file" ]; then
        echo "Error: File not found: $file"
        exit 1
    fi
done

PRIVATE_KEY_METADATA_FILE="${ENCRYPTED_PRIVATE_KEY_FILE}.metadata"

if [ ! -f "$PRIVATE_KEY_METADATA_FILE" ]; then
    echo "Metadata file not found: $PRIVATE_KEY_METADATA_FILE"
    exit 1
fi

# Read private key metadata
PRIVATE_KEY_SALT=$(grep "Salt:" "$PRIVATE_KEY_METADATA_FILE" | cut -d' ' -f2)
PRIVATE_KEY_ITERATIONS=$(grep "Iterations:" "$PRIVATE_KEY_METADATA_FILE" | cut -d' ' -f2)
PRIVATE_KEY_IV=$(grep "IV:" "$PRIVATE_KEY_METADATA_FILE" | cut -d' ' -f2)

echo "Debug: Private Key Salt: $PRIVATE_KEY_SALT"
echo "Debug: Private Key Iterations: $PRIVATE_KEY_ITERATIONS"
echo "Debug: Private Key IV: $PRIVATE_KEY_IV"

# Derive key for private key decryption
PRIVATE_KEY_INFO=$(./derive-key.sh "$PRIVATE_KEY_PASSWORD" "$PRIVATE_KEY_SALT" "$PRIVATE_KEY_ITERATIONS")
PRIVATE_KEY_ENCRYPTION_KEY=$(echo "$PRIVATE_KEY_INFO" | grep "Key:" | cut -d' ' -f2)

echo "Debug: Private Key Encryption Key (first 16 chars): ${PRIVATE_KEY_ENCRYPTION_KEY:0:16}..."

# Decrypt the private key (in memory)
DECRYPTED_PRIVATE_KEY=$(openssl enc -d -aes-256-cbc -in "$ENCRYPTED_PRIVATE_KEY_FILE" -K "$PRIVATE_KEY_ENCRYPTION_KEY" -iv "$PRIVATE_KEY_IV" 2>/dev/null)

if [ -z "$DECRYPTED_PRIVATE_KEY" ]; then
    echo "Error: Failed to decrypt the private key."
    exit 1
fi

echo "Private key decrypted successfully (in memory)."

# Now use the decrypted private key to decrypt the file password
PASS_FILE="${INPUT_ENCRYPTED_FILE}.pass"

if [ ! -f "$PASS_FILE" ]; then
    echo "Error: Encrypted password file not found: $PASS_FILE"
    exit 1
fi

RANDOM_PASS=$(echo "$DECRYPTED_PRIVATE_KEY" | openssl pkeyutl -decrypt -inkey /dev/stdin -in "$PASS_FILE")

if [ -z "$RANDOM_PASS" ]; then
    echo "Error: Failed to decrypt the file password."
    exit 1
fi

echo "Debug: Decrypted file password (first 8 chars): ${RANDOM_PASS:0:8}..."

# Finally, decrypt the input file using the decrypted password
openssl enc -d -aes-256-cbc -pbkdf2 -iter 100000 -in "$INPUT_ENCRYPTED_FILE" -out "$OUTPUT_DECRYPTED_FILE" -k "$RANDOM_PASS"

if [ $? -eq 0 ]; then
    echo "File decrypted successfully: $OUTPUT_DECRYPTED_FILE"
else
    echo "Error: Failed to decrypt the file."
    exit 1
fi

# Clear sensitive variables from memory
unset DECRYPTED_PRIVATE_KEY
unset RANDOM_PASS
unset PRIVATE_KEY_PASSWORD
unset PRIVATE_KEY_ENCRYPTION_KEY

echo "Decryption process completed."