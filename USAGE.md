# USAGE.md

Comprehensive usage guide for the worktree CLI tool. This document provides detailed usage patterns, examples, and best practices for all worktree commands.

## Table of Contents

- [Quick Start](#quick-start)
- [Command Overview](#command-overview)
- [worktree clone](#worktree-clone)
- [worktree (create)](#worktree-create)
- [worktree list](#worktree-list)
- [worktree -d (delete worktree)](#worktree--d-delete-worktree)
- [worktree --merge](#worktree---merge)
- [worktree --delete](#worktree---delete)
- [worktree completion](#worktree-completion)
- [Advanced Usage Patterns](#advanced-usage-patterns)
- [Error Handling and Troubleshooting](#error-handling-and-troubleshooting)
- [Best Practices](#best-practices)

---

## Quick Start

### Prerequisites

- Git installed and available in PATH
- Go 1.21+ (for building from source)
- Bash or compatible shell

### Installation

```bash
# Clone the repository
git clone https://github.com/lachlan/worktree.git
cd worktree

# Install the tool
./install.bash

# Verify installation
worktree --help
```

### Basic Workflow

```bash
# Clone a repository
worktree clone https://github.com/user/project.git

# Navigate to the cloned project
cd project/main

# Create a new feature branch with worktree
worktree feat/new-feature

# Work in the new worktree (automatically opens bash)
# ... make changes, commit, etc.

# Back in main repository, list branches
worktree list

# Merge the feature branch
worktree --merge --confirm feat/new-feature
```

---

## Command Overview

| Command | Description | Destination |
|---------|-------------|-------------|
| `worktree clone <url>` | Clone a repository into organized structure | Creates directory structure |
| `worktree <branch>` | Create new branch and worktree | `../<branch>/` |
| `worktree -r <branch>` | Create worktree from remote branch | `../<branch>/` |
| `worktree -e <branch>` | Create worktree from existing local branch | `../<branch>/` |
| `worktree list` | List all local branches | stdout |
| `worktree list -r` | List all remote branches | stdout |
| `worktree -d <branch>` | Delete worktree directory only | Removes directory |
| `worktree --merge <branch>` | Dry-run merge validation | stdout |
| `worktree --merge --confirm <branch>` | Execute merge and cleanup | Modifies git repo |
| `worktree --delete <branch>` | Dry-run delete validation | stdout |
| `worktree --delete --confirm <branch>` | Delete branch and worktree | Modifies git repo |
| `worktree completion <shell>` | Generate shell completion script | stdout |

---

## worktree clone

Clone a Git repository into a structured directory format.

### Syntax

```bash
worktree clone <repository-url>
```

### Description

The `clone` command creates a directory structure like:
```
current-directory/
└── project_name/
    └── [default-branch]/
        └── (repository contents)
```

Where `[default-branch]` is the repository's default branch (main, master, or trunk).

### Supported URL Formats

- **HTTPS**: `https://github.com/user/repo.git` or `https://github.com/user/repo`
- **SSH**: `git@github.com:user/repo.git` or `git@github.com:user/repo`
- **Git protocol**: `git://github.com/user/repo.git`

### Examples

```bash
# Clone using HTTPS
worktree clone https://github.com/user/project.git

# Clone using SSH
worktree clone git@github.com:user/project.git

# Clone without .git suffix
worktree clone https://github.com/user/project
```

### How It Works

1. **Extract project name**: Parses the URL to extract the repository name
2. **Create project directory**: Creates a directory with the project name
3. **Clone repository**: Uses `git clone` to clone into the project directory
4. **Determine default branch**: Checks origin/HEAD, then falls back to common branch names
5. **Rename directory**: Renames the cloned repository to the branch name

### Expected Output

```
Successfully cloned https://github.com/user/repo.git into project/main
```

### Common Issues

| Issue | Solution |
|-------|----------|
| URL not supported | Use standard GitHub/GitLab/Bitbucket URL formats |
| Authentication required | Ensure SSH keys or HTTPS credentials are set up |
| Repository not found | Verify the URL is correct and accessible |
| Permission denied | Check write permissions for current directory |

### Pro Tip

The project name extraction handles complex URLs:
- `git@github.com:org/repo.git` → `repo`
- `https://github.com/org/repo` → `repo`
- `git@gitlab.com:group/subgroup/project.git` → `project`

---

## worktree (create)

Create a new branch and worktree, or create a worktree from existing branches.

### Syntax

```bash
# Create new branch and worktree
worktree <branch-name>

# Create worktree from remote branch
worktree -r <branch-name>

# Create worktree from existing local branch
worktree -e <branch-name>
```

### Description

Creates a Git worktree at `../<branch-name>` relative to the current repository and starts a bash shell in that worktree.

### Flags

| Flag | Short | Description | Required |
|------|-------|-------------|----------|
| `--remote` | `-r` | Create local tracking branch from origin/<branch> | No |
| `--existing` | `-e` | Create worktree from existing local branch | No |

### Flag Combinations

| Combination | Behavior |
|-------------|----------|
| No flags | Creates new branch from current HEAD |
| `-r` only | Creates local tracking branch from remote |
| `-e` only | Creates worktree from existing local branch |
| `-r` + `-e` | ❌ Error: mutually exclusive |

### Examples

```bash
# Create new feature branch and worktree
worktree feat/new-feature

# Create worktree from existing remote branch
worktree -r origin/feat/existing-feature

# Create worktree from existing local branch
worktree -e feat/local-branch

# Create bug fix branch
worktree fix/bug-123

# Create documentation branch
worktree docs/readme-update
```

### How It Works

**Without flags (new branch)**:
1. Creates new branch from current HEAD: `git branch <branch-name>`
2. Creates worktree: `git worktree add ../<branch-name> <branch-name>`
3. Starts bash shell in the worktree

**With `-r` flag (remote branch)**:
1. Fetches remote branch: `git fetch origin <branch-name>:<branch-name>`
2. Creates worktree: `git worktree add ../<branch-name> <branch-name>`
3. Starts bash shell in the worktree

**With `-e` flag (existing local branch)**:
1. Validates branch exists locally
2. Creates worktree: `git worktree add ../<branch-name> <branch-name>`
3. Starts bash shell in the worktree

### Directory Creation

The worktree is created at `../<branch-name>` relative to the main repository:

```
If main repo is at: /home/user/projects/myapp/main
Creating worktree for: feat/new-feature
Worktree created at: /home/user/projects/myapp/feat/new-feature
```

**Note**: The path uses `/` as the separator. On Windows, it will use `\` automatically.

### Important Message

When creating a new branch (without `-r` or `-e`), you'll see:
```
Remember to push to create tracking branches in remote
```

This reminds you to push the new branch to the remote repository.

### Common Issues

| Issue | Solution |
|-------|----------|
| "WARNING: must be called from the main repository" | Run from directory where `.git` is a directory, not a file |
| "WARNING: not a git repository" | Ensure you're in a Git repository |
| "WARNING: failed to create branch" | Check branch name is valid and doesn't already exist |
| "WARNING: failed to create worktree" | Check worktree path doesn't already exist |
| "WARNING: cannot use switches -r and -e at the same time" | Use only one flag |

---

## worktree list

List all branches in the repository.

### Syntax

```bash
worktree list [flags]
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--remote` | `-r` | List remote branches instead of local |

### Examples

```bash
# List all local branches
worktree list

# List all remote branches
worktree list -r

# List with remote flag
worktree list --remote
```

### Output

**Local branches**:
```
  main
  feat/new-feature
  fix/bug-123
  docs/readme
```

**Remote branches** (with `-r`):
```
  origin/main
  origin/feat/new-feature
  origin/feat/existing
  origin/fix/bug-123
```

### How It Works

- **Local branches**: Runs `git branch --format=%(refname:short)`
- **Remote branches**: Runs `git branch -r --format=%(refname:short)`

### Common Issues

| Issue | Solution |
|-------|----------|
| No output | Repository has no branches |
| "WARNING: must be called from the main repository" | Run from main repository, not worktree |
| "WARNING: not a git repository" | Ensure you're in a Git repository |

---

## worktree -d (delete worktree)

Delete a worktree directory without affecting the branch.

### Syntax

```bash
worktree -d <branch-name>
worktree --delete-worktree <branch-name>
```

### Description

Removes the worktree directory and cleans up empty parent directories. The branch itself remains intact.

### Flags

| Flag | Long Form | Description |
|------|-----------|-------------|
| `-d` | `--delete-worktree` | Delete the worktree directory |

### Examples

```bash
# Delete worktree for feature branch
worktree -d feat/new-feature

# Using long form
worktree --delete-worktree feat/new-feature

# Delete multiple worktrees
worktree -d feat/old-feature
worktree -d fix/bug-456
```

### How It Works

1. Removes the worktree: `git worktree remove ../<branch-name>`
2. Cleans up empty parent directories recursively
3. Only removes directories that are empty (safe deletion)

### Cleanup Process

If worktree is at `../feat/new-feature`:
1. Remove worktree directory
2. Check if `../feat/` is now empty
3. If empty, remove `../feat/`
4. Continue up the directory tree until a non-empty directory is found

### Expected Output

```
Deleted worktree at ../feat/new-feature
```

### Common Issues

| Issue | Solution |
|-------|----------|
| "WARNING: could not remove worktree" | Worktree may not be registered or path is incorrect |
| "WARNING: must be called from the main repository" | Run from main repository |
| Worktree directory still exists | There may be files in it - clean up manually |

### Important Notes

- This only removes the worktree directory, not the branch
- The branch can still be checked out normally
- Empty parent directories are automatically cleaned up
- Non-empty directories are preserved

---

## worktree --merge

Merge a branch into the main branch and clean up the worktree.

### Syntax

```bash
# Dry-run (default) - validate before merging
worktree --merge <branch-name>

# Execute the merge
worktree --merge --confirm <branch-name>
```

### Description

Validates that a merge operation is possible, then optionally executes it with confirmation.

### Flags

| Flag | Description | Required |
|------|-------------|----------|
| `--merge` | Merge the branch into main and clean up | Yes |
| `--confirm` | Execute the operation (otherwise dry-run) | Yes for execution |

### Examples

```bash
# Check if merge is possible (dry-run)
worktree --merge feat/new-feature

# Execute the merge
worktree --merge --confirm feat/new-feature

# Multiple step process
worktree --merge feat/test  # Check first
worktree --merge --confirm feat/test  # Then execute
```

### Dry-Run Mode

When run without `--confirm`, the tool performs comprehensive validation:

1. **Branch status check**: Compares branch against main
2. **Ahead/behind analysis**: Shows commit differences
3. **Merge test**: Attempts a no-commit merge to validate
4. **Worktree validation**: Checks worktree exists and can be removed
5. **Branch deletion check**: Validates branch can be deleted
6. **Remote branch check**: Validates remote branch can be deleted

### Dry-Run Output Example

```
=== DRY RUN - Testing if operations are possible ===

Checking: Branch feat/new-feature status relative to main
  INFO: Branch feat/new-feature has 3 commit(s) ahead of main
  INFO: Branch feat/new-feature is up to date with main
  PASS: Merge is possible

Testing: Delete worktree at ../feat/new-feature
  PASS: Worktree can be removed

Testing: Delete local branch feat/new-feature
  PASS: Local branch can be deleted (fully merged)

Testing: Delete remote branch origin/feat/new-feature
  WARNING: Remote branch origin/feat/new-feature does not exist (will be skipped)

All operations are possible!

To execute, run: worktree --merge --confirm feat/new-feature
Warning: ensures no one else is working on this branch - it will be deleted
```

### Execution Mode

With `--confirm`, the tool:
1. Performs the actual merge: `git merge <branch>`
2. Removes the worktree directory
3. Cleans up empty parent directories
4. Deletes the local branch: `git branch -d <branch>`
5. Deletes the remote branch: `git push -d origin <branch>` (if exists)

### Important Notes

- **Main branch detection**: Uses origin/HEAD, then tries main, master, trunk
- **Merge validation**: Tests merge with `--no-commit --no-ff` first
- **Conflict detection**: If merge would fail, dry-run shows the error
- **Safety**: Requires explicit `--confirm` to execute

### Common Issues

| Issue | Solution |
|-------|----------|
| Branch is behind main | Rebase or merge main into the branch first |
| Merge conflicts | Resolve conflicts, then run with `--confirm` |
| Worktree not found | Ensure worktree exists at expected path |
| Branch not fully merged | Use `-D` flag for force deletion if needed |

---

## worktree --delete

Delete a branch, its remote counterpart, and the worktree directory.

### Syntax

```bash
# Dry-run (default) - validate before deleting
worktree --delete <branch-name>

# Execute the deletion
worktree --delete --confirm <branch-name>
```

### Description

Completely removes a branch including local, remote, and worktree. More comprehensive than `--merge` which only merges and cleans up.

### Flags

| Flag | Description | Required |
|------|-------------|----------|
| `--delete` | Delete the branch, remote branch, and worktree | Yes |
| `--confirm` | Execute the operation (otherwise dry-run) | Yes for execution |

### Examples

```bash
# Check if deletion is possible (dry-run)
worktree --delete feat/old-feature

# Execute the deletion
worktree --delete --confirm feat/old-feature

# Cannot combine with other flags
worktree --delete --confirm -r feat/test  # ❌ Error
```

### Differences from --merge

| Aspect | `--merge` | `--delete` |
|--------|-----------|------------|
| Merges branch | ✅ Yes | ❌ No |
| Deletes local branch | ✅ Yes | ✅ Yes |
| Deletes remote branch | ✅ Yes | ✅ Yes |
| Removes worktree | ✅ Yes | ✅ Yes |
| Cleanup empty dirs | ✅ Yes | ✅ Yes |

**Use `--merge`** when you want to incorporate changes into main.  
**Use `--delete`** when you want to discard the branch entirely.

### Dry-Run Output Example

```
=== DRY RUN - Testing if operations are possible ===

Checking: Branch feat/abandoned status relative to main
  INFO: Branch feat/abandoned has 1 commit(s) ahead of main
  FAIL: Branch feat/abandoned is 2 commit(s) behind main

Testing: Delete worktree at ../feat/abandoned
  PASS: Worktree can be removed

Testing: Delete local branch feat/abandoned
  PASS: Local branch can be deleted (with force -D)

Testing: Delete remote branch origin/feat/abandoned
  PASS: Remote branch can be deleted

Some operations would fail - see above for details

To execute, run: worktree --delete --confirm feat/abandoned
Warning: ensures no one else is working on this branch - it will be deleted
```

### Execution Mode

With `--confirm`, the tool:
1. Removes the worktree directory
2. Cleans up empty parent directories
3. Deletes the local branch (with force if not merged)
4. Deletes the remote branch

### Important Notes

- **Force deletion**: If branch is not fully merged, uses `-D` flag for local deletion
- **Remote branch**: Attempts deletion but continues if remote doesn't exist
- **Irreversible**: Once executed, branch and worktree are permanently removed
- **Safety**: Requires explicit `--confirm` to execute

### Common Issues

| Issue | Solution |
|-------|----------|
| Branch is behind main | Use `--delete` to force delete (changes will be lost) |
| Remote branch doesn't exist | Normal - will be skipped with warning |
| Branch in use | Ensure no one else is working on the branch |

---

## worktree completion

Generate shell completion scripts for worktree commands.

### Syntax

```bash
worktree completion <shell-type>
```

### Supported Shells

| Shell | Command |
|-------|---------|
| Bash | `worktree completion bash` |
| Zsh | `worktree completion zsh` |
| Fish | `worktree completion fish` |
| PowerShell | `worktree completion powershell` |

### Usage

#### Bash

```bash
# Load immediately (current session)
source <(worktree completion bash)

# Install permanently (all future sessions)
worktree completion bash > /etc/bash_completion.d/worktree

# Or to home directory
worktree completion bash > ~/.local/share/bash-completion/completions/worktree
```

#### Zsh

```zsh
# Load immediately
source <(worktree completion zsh)

# Install permanently (add to ~/.zshrc)
worktree completion zsh > "${fpath[1]}/_worktree"

# Then add to ~/.zshrc
# autoload -U compinit && compinit
```

#### Fish

```fish
# Load immediately
worktree completion fish | source

# Install permanently
worktree completion fish > ~/.config/fish/completions/worktree.fish
```

#### PowerShell

```powershell
# Load immediately
worktree completion powershell | Out-String | Invoke-Expression

# Install permanently (add to profile)
worktree completion powershell > $PROFILE.CurrentUserCurrentHost
```

### Completion Features

- **Branch name completion**: Tab completes branch names
- **Local and remote**: Shows both local and remote branches
- **Prefix filtering**: Only shows branches matching current input
- **Deduplication**: Removes duplicate branch names
- **Context-aware**: Only completes branch names for positional arguments

### Examples

```bash
# Type 'worktree ' then press Tab twice
worktree <TAB><TAB>
# Shows: feat/new-feature  fix/bug-123  main  docs/readme

# Type 'worktree feat/' then press Tab
worktree feat/<TAB>
# Shows: new-feature  existing-feature

# Works with flags too
worktree -r <TAB><TAB>
# Shows: feat/new-feature  fix/bug-123  main  docs/readme
```

### Common Issues

| Issue | Solution |
|-------|----------|
| Completion not working | Source the completion script or restart shell |
| No branches shown | Ensure you're in a Git repository |
| Slow completion | Large repositories may have slower completion |

---

## Advanced Usage Patterns

### Multi-Repository Workflow

```bash
# Clone multiple repositories
worktree clone https://github.com/user/frontend.git
worktree clone https://github.com/user/backend.git

# Organized structure
projects/
├── frontend/
│   └── main/
└── backend/
    └── main/
```

### Feature Branch Development

```bash
# Start new feature
worktree feat/user-authentication

# In the worktree: make changes, commit
# ... development work ...

# Push to remote
git push -u origin feat/user-authentication

# Back in main: merge when ready
worktree --merge --confirm feat/user-authentication
```

### Bug Fix with Remote Branch

```bash
# Create worktree from existing remote branch
worktree -r fix/urgent-bug

# Fix the issue, commit, push
# ... fix work ...

# Merge back to main
worktree --merge --confirm fix/urgent-bug
```

### Parallel Development

```bash
# Create multiple feature branches
worktree feat/login
worktree feat/dashboard
worktree feat/api-integration

# Work on each in separate terminals
# Each worktree has isolated working directory
```

### Cleanup Old Branches

```bash
# Check what can be deleted
worktree --delete feat/old-feature

# If all checks pass, execute deletion
worktree --delete --confirm feat/old-feature
```

### Complex Directory Structures

The tool creates organized structures:

```
projects/
├── myapp/
│   ├── main/
│   │   └── .git/  # Main repository
│   ├── feat/
│   │   ├── user-auth/
│   │   └── payment/
│   ├── fix/
│   │   └── security/
│   └── docs/
│       └── readme/
└── otherapp/
    ├── master/  # Some repos use master
    └── develop/
```

---

## Error Handling and Troubleshooting

### Common Error Messages

| Error | Meaning | Solution |
|-------|---------|----------|
| `WARNING: must be called from the main repository, not a worktree` | Command run from worktree directory | `cd` to main repository directory |
| `WARNING: not a git repository` | Not in a Git repo | Navigate to a Git repository |
| `WARNING: cannot use switches -r and -e at the same time` | Conflicting flags | Use only one flag |
| `WARNING: cannot use -r, -e, or -d with --merge or --delete` | Invalid flag combination | Use only merge/delete flags |
| `WARNING: cannot use --confirm is only used with --merge or --delete` | Misplaced confirm flag | Add with merge or delete flag |
| `WARNING: branch '<name>' does not exist` | Branch doesn't exist | Check branch name spelling |
| `WARNING: failed to create worktree` | Worktree path issue | Check directory permissions |

### Debugging Tips

1. **Check current directory**:
   ```bash
   pwd
   ls -la
   ```

2. **Verify Git status**:
   ```bash
   git status
   git branch -a
   ```

3. **Check worktree status**:
   ```bash
   git worktree list
   ```

4. **Test with --help**:
   ```bash
   worktree --help
   worktree clone --help
   ```

5. **Use dry-run mode**:
   ```bash
   worktree --merge feat/test
   worktree --delete feat/test
   ```

### Worktree Specific Issues

**Issue**: Worktree exists but not showing in `git worktree list`

```bash
# Check actual directory
ls -la ../branch-name

# Re-register worktree
git worktree add ../branch-name branch-name
```

**Issue**: Branch exists but worktree command can't find it

```bash
# Verify branch exists
git branch -a | grep branch-name

# Check if in main repo
ls -la .git  # Should be directory, not file
```

---

## Best Practices

### Directory Organization

1. **Use meaningful branch names**:
   - `feat/user-authentication` (not `feat/auth`)
   - `fix/login-bug-123` (not `fix/bug`)
   - `docs/api-documentation` (not `docs/fix`)

2. **Group related branches**:
   - Use `feat/`, `fix/`, `docs/`, `chore/` prefixes
   - This creates organized directory structures

3. **Keep main repository clean**:
   - Don't commit work-in-progress to main
   - Use feature branches for all development

### Workflow Efficiency

1. **Use tab completion**:
   - Branch names auto-complete
   - Reduces typing errors

2. **Always dry-run first**:
   - `worktree --merge branch` before `--confirm`
   - `worktree --delete branch` before `--confirm`

3. **Regular cleanup**:
   - Delete old worktrees with `worktree -d`
   - Merge and delete completed feature branches
   - Keep only active branches

4. **Team coordination**:
   - Communicate before deleting branches
   - Ensure no one is working on a branch before deletion
   - Use `--merge --confirm` for completed features

### Performance Tips

1. **Large repositories**:
   - Dry-run operations may be slow
   - Branch completion may be slow
   - Consider using `-r` flag to list only remote branches

2. **Many branches**:
   - Use specific branch name patterns
   - Tab completion filters as you type

3. **Network issues**:
   - Remote operations require network access
   - Use local operations when offline

### Security Best Practices

1. **Repository URLs**:
   - Use SSH for private repositories
   - Validate URLs before cloning

2. **Branch names**:
   - Avoid special characters in branch names
   - Use alphanumeric, hyphen, underscore, slash

3. **File permissions**:
   - Worktrees inherit repository permissions
   - Be cautious with file system operations

---

## Summary

The worktree tool provides a powerful and efficient way to manage Git branches and worktrees. By following the patterns and best practices outlined in this guide, you can maximize productivity while maintaining a clean and organized development environment.

**Remember**:
- Always run from the main repository
- Use dry-run mode for destructive operations
- Leverage autocompletion for branch names
- Keep your directory structure organized
- Communicate with your team before deleting branches

For more information, see:
- [README.md](./README.md) - Main documentation
- [WORKFLOW.md](./WORKFLOW.md) - Development workflow guide
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Technical architecture
- [AGENTS.md](./AGENTS.md) - AI agent guidance