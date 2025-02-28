#!/bin/bash

cmd_prefix="replay"

for folder in */; do
    # Remove the trailing slash
    folder="${folder%/}"

    # Change to the current folder
    cd "$folder"
    
    # Print the desired string. Replace "String for folder" with the actual string you want to print.
    # echo "String for folder: $cmd_prefix.$folder"

    go mod tidy
    go build -o "$cmd_prefix.$folder" main.go

    # Move back to the target directory
    cd ../
done