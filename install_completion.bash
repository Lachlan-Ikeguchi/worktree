#!/bin/bash

# Install completion scripts for worktree and clone commands
# This script installs completion scripts to ~/.bash_completion.d/

set -e

echo "Installing completion scripts..."

# Create the completion directory if it doesn't exist
COMPLETION_DIR="$HOME/.bash_completion.d"
mkdir -p "$COMPLETION_DIR"

# Build worktree
if [ -f cmd/worktree/main.go ]; then
    echo "Building worktree..."
    go build -o worktree ./cmd/worktree/
    ./worktree completion bash > "$COMPLETION_DIR/worktree"
    rm worktree
    echo "Installed worktree completion to $COMPLETION_DIR/worktree"
else
    echo "Warning: cmd/worktree/main.go not found, skipping worktree"
fi

# Build clone
if [ -f cmd/clone/main.go ]; then
    echo "Building clone..."
    go build -o clone ./cmd/clone/
    ./clone completion bash > "$COMPLETION_DIR/clone"
    rm clone
    echo "Installed clone completion to $COMPLETION_DIR/clone"
else
    echo "Warning: cmd/clone/main.go not found, skipping clone"
fi

# Check if .bashrc exists and add sourcing if not already present
if [ -f "$HOME/.bashrc" ]; then
    if ! grep -q "bash_completion.d" "$HOME/.bashrc"; then
        echo "" >> "$HOME/.bashrc"
        echo "# Source bash completion scripts" >> "$HOME/.bashrc"
        echo "for file in ~/.bash_completion.d/*; do" >> "$HOME/.bashrc"
        echo "    [ -f \"\$file\" ] && source \"\$file\"" >> "$HOME/.bashrc"
        echo "done" >> "$HOME/.bashrc"
        echo "Added completion sourcing to ~/.bashrc"
    else
        echo "Completion sourcing already exists in ~/.bashrc"
    fi
else
    echo "Warning: ~/.bashrc not found, completion scripts installed but not auto-sourced"
    echo "To enable, add the following to your .bashrc:"
    echo "    for file in ~/.bash_autocomplete.d/*; do"
    echo "        [ -f \"\$file\" ] && source \"\$file\""
    echo "    done"
fi

echo ""
echo "Installation complete!"
echo "Completion scripts installed to: $COMPLETION_DIR/"
echo "Restart your shell or run: source ~/.bashrc"
