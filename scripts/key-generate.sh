#!/bin/bash

# Generate 16384-bit RSA key pair
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:16384
openssl rsa -pubout -in private_key.pem -out public_key.pem

echo "16384-bit RSA key pair generated:"
echo "Private key: private_key.pem"
echo "Public key: public_key.pem"

echo "Remember to encrypt and securely store the private key using the manual encryption method."
