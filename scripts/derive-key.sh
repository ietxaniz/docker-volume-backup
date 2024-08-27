#!/bin/bash

usage() {
    echo "Usage: $0 <password> [salt] [iterations]"
    echo "  If salt is not provided, a random one will be generated."
    echo "  Default iterations: 100000"
}

if [ "$#" -lt 1 ] || [ "$#" -gt 3 ]; then
    usage
    exit 1
fi

PASSWORD=$1
SALT=${2:-$(openssl rand -hex 8)}
ITERATIONS=${3:-1000000}

start_time=$(date +%s.%N)

# Use the provided salt in the key derivation
KEY=$(printf "%s" "$PASSWORD" | openssl enc -aes-256-cbc -pbkdf2 -iter $ITERATIONS -S "$SALT" -pass stdin -P 2>/dev/null | grep "key=" | cut -d'=' -f2)

if [ -z "$KEY" ]; then
    echo "Error: Failed to derive key" >&2
    exit 1
fi

end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)

echo "Salt: $SALT"
echo "Iterations: $ITERATIONS"
echo "Key: $KEY"
echo "Time taken: $duration seconds"