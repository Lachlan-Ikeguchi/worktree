#!/bin/bash

# Install completion scripts for worktree command (includes clone subcommand)
# This script installs completion scripts to ~/.bash_completion.d/

set -e

echo "Installing completion scripts..."

# Create the completion directory if it doesn't exist
COMPLETION_DIR="$HOME/.bash_completion.d"
mkdir -p "$COMPLETION_DIR"

# Build worktree (includes clone subcommand)
if [ -f cmd/worktree/main.go ]; then
    echo "Building worktree (includes clone subcommand)..."
    go build -o worktree ./cmd/worktree/
    ./worktree completion bash > "$COMPLETION_DIR/worktree"
    rm worktree
    echo "Installed worktree completion to $COMPLETION_DIR/worktree"
    echo "Note: worktree now includes the clone subcommand"
else
    echo "Warning: cmd/worktree/main.go not found, skipping worktree"
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
    echo "    for file in ~/.bash_completion.d/*; do"
    echo "        [ -f \"\$file\" ] && source \"\$file\""
    echo "    done"
fi

echo ""
echo "Installation complete!"
echo "Completion scripts installed to: $COMPLETION_DIR/"
echo "Restart your shell or run: source ~/.bashrc"
