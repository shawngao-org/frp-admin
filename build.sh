#!/usr/bin/env bash

read -r -p "Enter target system[windows, darwin, linux]: " targetSystem
read -r -p "Enter target platform[386, amd64, arm]: " targetPlatform
CGO_ENABLED=0 GOOS=$targetSystem GOARCH=$targetPlatform go build cloud-lock-go-gin
