#!/bin/bash

# Generate ECC key pair
openssl ecparam -name prime256v1 -genkey -noout -out private_key.pem
openssl ec -in private_key.pem -pubout -out public_key.pem

echo "ECC key pair generated:"
echo "Private key: private_key.pem"
echo "Public key: public_key.pem"

echo "Remember to encrypt and securely store the private key using the manual encryption method."
