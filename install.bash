#!/usr/bin/env bash

# Install worktree tools
# Builds Go binaries and copies them to ~/bin/

set -e

# Create bin directory if it doesn't exist
mkdir -p ~/bin

# Build Go binaries
echo "Building worktree (includes clone subcommand)..."
go build -o ~/bin/worktree ./cmd/worktree

echo "Installation complete!"
echo "Binary installed to ~/bin/worktree"
echo "Make sure ~/bin is in your PATH"
echo ""
echo "Usage: worktree clone <repo-url> to clone repositories"
