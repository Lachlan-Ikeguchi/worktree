# Worktree

A wrapper around git to work with worktrees and repository cloning more efficiently.

## Overview

The `worktree` command includes functionality for:
- **worktree**: Creates and manages Git worktrees for branch-based development workflows
- **worktree clone**: Clones a Git repository into `project_name/[branch]/` where branch is the default branch (main/master/trunk)

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

### worktree init

Initialize a new project with worktree structure:

```bash
worktree init <project-name>
```

Creates a project directory with a subdirectory named after the main branch (master/main/trunk) and initializes a git repository in that subdirectory.

**Examples:**
```bash
worktree init myproject
worktree init test_project
```

The resulting structure will be:
```
.
└── myproject/
    └── master/  (or main/trunk based on git config)
        └── .git/ (newly initialized git repository)
```

The main branch name is determined from git config (`init.defaultBranch`) or defaults to 'master' if not configured.

### worktree clone

Clone a Git repository:

```bash
worktree clone <repository-url>
```

The repository will be cloned into a directory structure like:
```
project_name/[branch]/
```

Where `branch` is the repository's default branch (main, master, or trunk).

**Examples:**
```bash
worktree clone https://github.com/user/repo.git
worktree clone git@github.com:user/repo.git
```

### worktree

The worktree command creates and manages worktrees at `../[branch]` relative to the main repository.

#### Create a new worktree

Create a worktree with auto-detection (checks remote first, then local, then creates new from HEAD):
```bash
worktree <branch-name>
```

When both remote and local branches exist, the remote branch is prioritized.

#### List branches

List all branches (local and remote):
```bash
worktree list
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
| `-d, --delete-worktree` | Delete the worktree directory and clean up empty parent directories |
| `--merge` | Merge branch into main/master and clean up |
| `--delete` | Delete branch, remote branch, and worktree |
| `--confirm` | Confirm merge or delete operation (required for --merge and --delete) |
| `-h, --help` | Show help message |

## Subcommands

| Command | Description |
|---------|-------------|
| `worktree init <project>` | Initialize a new project with worktree structure |
| `worktree clone <repo-url>` | Clone a repository into project_name/[branch]/ structure |
| `worktree list` | List all branches (local and remote) |

## Autocompletion

The `worktree` command supports shell autocompletion through Cobra, including the `clone` subcommand. The `worktree` command provides **branch name autocompletion** - when you type `worktree <TAB><TAB>`, it will automatically suggest available local and remote branch names.

### Setup
**Bash:**
```bash
# Load immediately
source <(worktree completion bash)
# Or install permanently
worktree completion bash > /etc/bash_completion.d/worktree
```

**Zsh:**
```zsh
# Load immediately
source <(worktree completion zsh)
# Or install permanently (add to ~/.zshrc)
worktree completion zsh > "${fpath[1]}/_worktree"
```

**Fish:**
```fish
# Load immediately
worktree completion fish | source
# Or install permanently
worktree completion fish > ~/.config/fish/completions/worktree.fish
```

**PowerShell:**
```powershell
# Load immediately
worktree completion powershell | Out-String | Invoke-Expression
# Or install permanently (add to your profile)
worktree completion powershell > $PROFILE.CurrentUserCurrentHost
```

View completion help:
```bash
worktree completion --help
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
- `--merge` and `--delete` require `--confirm` to execute
- Dry-run mode for `--merge` and `--delete` includes worktree and branch existence validation
- Dry-run output uses color-coded messages: green (PASS), red (FAIL), yellow (WARNING), blue (INFO)
- Auto-detection prioritizes remote branches, then local branches, then creates new from HEAD
