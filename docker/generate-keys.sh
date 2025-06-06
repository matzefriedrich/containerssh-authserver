#!/bin/sh

scriptPath="$(cd "$(dirname "$0")" && pwd)"
cd $scriptPath

user=johndoe

# Generate the SSH host key
openssl genrsa > keys/ssh_host_rsa_key.pem

# Generate a new SSH key for the demo user
ssh-keygen -t ed25519 -f keys/$user.pem -N "" -C $user@localhost -q
