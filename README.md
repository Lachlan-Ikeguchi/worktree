# Worktree and Clone Scripts

These are wrappers around git to work with worktrees more efficiently.

## Overview

- **clone**: Clones a Git repository into `project_name/[branch]/` where branch is the default branch (main/master/trunk)
- **worktree**: Creates and manages Git worktrees for branch-based development workflows

## Installation

Run the install script to build and install the binaries:

```bash
./install.bash
```

This builds the Go binaries and installs them to `~/bin/`. Ensure `~/bin` is in your PATH.

## Uninstallation

```bash
./uninstall.bash
```

## Usage

### clone

Clone a Git repository:

```bash
clone <repository-url>
```

The repository will be cloned into a directory structure like:
```
project_name/[branch]/
```

Where `branch` is the repository's default branch (main, master, or trunk).

**Examples:**
```bash
clone https://github.com/user/repo.git
clone git@github.com:user/repo.git
```

### worktree

The worktree command creates and manages worktrees at `../[branch]` relative to the main repository.

#### Create a new worktree

Create a new branch and worktree:
```bash
worktree <branch-name>
```

Create a worktree from an existing remote branch:
```bash
worktree -r <branch-name>
```

Create a worktree from an existing local branch:
```bash
worktree -e <branch-name>
```

#### List branches

List all local branches:
```bash
worktree list
```

List all remote branches:
```bash
worktree list -r
```

#### Delete a worktree

Delete the worktree directory:
```bash
worktree -d <branch-name>
```

#### Merge and clean up

Merge a branch into main/master and clean up (dry-run by default):
```bash
worktree --merge <branch-name>
```

To actually perform the merge and cleanup:
```bash
worktree --merge --confirm <branch-name>
```


#### Delete branch and worktree

Delete the local branch, remote branch, and worktree (dry-run by default):
```bash
worktree --delete <branch-name>
```

To actually perform the deletion:
```bash
worktree --delete --confirm <branch-name>
```

## Options

| Flag | Description |
|------|-------------|
| `-r` | Create local tracking branch from origin/<branch> |
| `-e` | Create worktree from existing local branch |
| `-d` | Delete the worktree directory and clean up empty parent directories |
| `--merge` | Merge branch into main/master and clean up |
| `--delete` | Delete branch, remote branch, and worktree |
| `--confirm` | Confirm merge or delete operation (required for --merge and --delete) |
| `-h, --help` | Show help message |

## Autocompletion

Both `worktree` and `clone` commands support shell autocompletion through Cobra. The `worktree` command now includes **branch name autocompletion** - when you type `worktree <TAB><TAB>`, it will automatically suggest available local and remote branch names.

### Setup
**Bash:**
```bash
# Load immediately
source <(worktree completion bash)
# Or install permanently
worktree completion bash > /etc/bash_completion.d/worktree
clone completion bash > /etc/bash_completion.d/clone
```

**Zsh:**
```zsh
# Load immediately
source <(worktree completion zsh)
# Or install permanently (add to ~/.zshrc)
worktree completion zsh > "${fpath[1]}/_worktree"
clone completion zsh > "${fpath[1]}/_clone"
```

**Fish:**
```fish
# Load immediately
worktree completion fish | source
# Or install permanently
worktree completion fish > ~/.config/fish/completions/worktree.fish
clone completion fish > ~/.config/fish/completions/clone.fish
```

**PowerShell:**
```powershell
# Load immediately
worktree completion powershell | Out-String | Invoke-Expression
# Or install permanently (add to your profile)
worktree completion powershell > $PROFILE.CurrentUserCurrentHost
clone completion powershell > $PROFILE.CurrentUserCurrentHost
```

View completion help:
```bash
worktree completion --help
clone completion --help
```

### Branch Name Completion

The `worktree` command provides intelligent branch name completion:
- Tab completion after `worktree ` will show all available local and remote branches
- Branches are filtered as you type (e.g., `worktree feat<TAB>` will only show branches starting with "feat")
- Both local and remote branches are included, with duplicates removed

## Directory Structure

After cloning and creating worktrees, your directory structure will look like:

```
project_name/
├── main/           # Main branch (or master/trunk)
│   └── ...
├── feat/
│   ├── feature-a/  # Worktree for feature-a branch
│   │   └── ...
│   └── feature-b/  # Worktree for feature-b branch
│       └── ...
├── fix/
│   └── bug-fix/    # Worktree for bug-fix branch
│       └── ...
└── docs/
    └── readme/     # Worktree for docs/readme branch
        └── ...
```

## Notes

- All commands must be run from within the main repository (not from a worktree)
- The `-r` and `-e` flags are mutually exclusive
- `--merge` and `--delete` require `--confirm` to execute
  - Worktree and branch existence validation
