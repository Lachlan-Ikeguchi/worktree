#!/usr/bin/env bash

# Install worktree tools
# Builds Go binaries and copies them to ~/bin/

set -e

# Create bin directory if it doesn't exist
mkdir -p ~/bin

# Build Go binaries
echo "Building clone..."
go build -o ~/bin/clone ./cmd/clone

echo "Building worktree..."
go build -o ~/bin/worktree ./cmd/worktree

echo "Installation complete!"
echo "Binaries installed to ~/bin/"
echo "Make sure ~/bin is in your PATH"
