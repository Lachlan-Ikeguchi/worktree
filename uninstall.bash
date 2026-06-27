#!/usr/bin/env bash

# Uninstall worktree tools
# Removes Go binaries from ~/bin/

set -e

# Remove Go binaries
if [ -f ~/bin/worktree ]; then
    rm ~/bin/worktree
    echo "Removed ~/bin/worktree"
fi

echo "Uninstallation complete!"
