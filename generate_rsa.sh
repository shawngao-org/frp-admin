#!/bin/bash

function generate_key() {
    echo "Generating RSA key pair..."
    rm private_key.pem public_key.pem
    openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out private_key.pem
    openssl rsa -pubout -in private_key.pem -out public_key.pem
    echo "RSA key pair successfully generated."
}

function verify() {
  echo "Verify public key and private key..."
  echo "Q2ju-x8ru4s]tYCMt4KxAzxHTi@>V%407c~e.d}}A:mY:Bizggbr1g69zfzib3+c" > verify_data.txt
  openssl dgst -sha256 -sign private_key.pem -out signature.bin verify_data.txt
  openssl dgst -sha256 -verify public_key.pem -signature signature.bin verify_data.txt
  verify_result=$?
  rm verify_data.txt signature.bin
  if [ $verify_result -eq 0 ]; then
      echo "Verify successfully."
  else
      echo "Verify failed."
      generate_key
  fi
  echo "Done."
}

echo "Check public key and private key exists..."
if [ -f "private_key.pem" ] && [ -f "public_key.pem" ]; then
  echo "Public key and private key already exists."
  verify
  exit 0
fi
generate_key
verify
