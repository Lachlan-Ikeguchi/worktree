# ARCHITECTURE.md

Technical architecture documentation for the worktree CLI tool. This document explains the internal structure, design decisions, and technical implementation of the tool.

## Table of Contents

- [System Overview](#system-overview)
- [Architecture Diagrams](#architecture-diagrams)
- [Source Code Structure](#source-code-structure)
- [Command Hierarchy](#command-hierarchy)
- [Core Components](#core-components)
- [Key Technical Decisions](#key-technical-decisions)
- [Integration with Git](#integration-with-git)
- [Shell Autocompletion System](#shell-autocompletion-system)
- [Color Coding System](#color-coding-system)
- [Error Handling Strategy](#error-handling-strategy)
- [Validation and Safety Mechanisms](#validation-and-safety-mechanisms)
- [Performance Considerations](#performance-considerations)
- [Security Considerations](#security-considerations)
- [Extensibility and Future Architecture](#extensibility-and-future-architecture)

---

## System Overview

### High-Level Architecture

The worktree tool is a **Go-based CLI application** that wraps Git commands to provide enhanced worktree management. It follows a **command-pattern architecture** with Cobra as the CLI framework.

The system is organized into several interacting layers:

**User Interface Layer** handles all user interaction:
- stdin/stdout for command input and output
- Exit codes for success/failure signaling
- Interactive shell sessions (automatically starts bash in new worktrees)

**Cobra CLI Framework** provides the command infrastructure:
- Command parsing and argument handling
- Flag definition and processing
- Built-in help system generation
- Subcommand organization and hierarchical command structure

**Git Integration Layer** executes all git operations:
- Uses `exec.Command()` to run git commands as subprocesses
- Handles command execution, output capture, and error processing
- Provides abstraction over git operations like branch management, worktree operations, and merges

**Shell Autocompletion System** provides intelligent command completion:
- Branch name completion combining local and remote branches
- Context-aware completion for different argument positions
- Support for Bash, Zsh, Fish, and PowerShell
- Deduplication and prefix filtering for efficient completion

**Repository State Validation** ensures commands run in valid contexts:
- `isGitRepo()` checks if current directory contains a .git directory (main repository)
- `isWorktree()` checks if current directory contains a .git file (worktree)
- Validates repository state before executing operations

**Color Coding System** provides visual feedback:
- ANSI color codes for PASS (green), FAIL (red), WARNING (yellow), INFO (blue)
- Consistent color semantics across all output

These layers work together with the following relationships:
- User Interface connects to Cobra CLI for command processing
- Cobra CLI calls Git Integration Layer for repository operations
- Git Integration Layer depends on Repository State Validation
- Shell Autocompletion interacts with both Cobra CLI (for command structure) and Git Integration (for branch data)
- Color Coding is used by all layers for user feedback

### Technology Stack

| Component | Technology | Version | Purpose |
|-----------|------------|---------|---------|
| **Language** | Go | 1.21+ | Core application logic |
| **CLI Framework** | Cobra | v1.7.0 | Command structure and parsing |
| **Dependency Management** | Go Modules | Built-in | Dependency resolution |
| **Shell Integration** | Standard OS | N/A | Shell execution and completion |
| **Git Integration** | Git CLI | 2.x+ | Repository operations |

### Module Structure

```
github.com/lachlan/worktree
├── go.mod          # Module definition
├── go.sum          # Dependency checksums
└── cmd/
    └── worktree/
        └── main.go  # Application entry point
```

---

## Command Flow Description

### Command Execution Flow

When a user runs `worktree feat/new-feature`, the command processing follows this sequence:

1. **Cobra Root Command Processing** (rootCmd in main.go):
   - Parse the command name and arguments
   - Validate flag combinations are valid
   - Execute PersistentPreRun function if applicable (validates git context)
   - Call the rootCmd.Run() function to handle the command

2. **Validation Phase**:
   - `isWorktree()` - Check if current directory is a worktree (has .git as a file)
   - `isGitRepo()` - Check if current directory is the main git repository (has .git as a directory)
   - Validate argument count matches command requirements
   - Validate flag exclusivity (e.g., -r vs -e, merge vs delete cannot be combined)

3. **Command Routing**: Based on flags and arguments:
   - If `mergeMode` OR `deleteMode` is true: call `mergeOrDelete(branch, flags, confirm)`
   - If `deleteFlag` is true: call `deleteWorktree(branch)`
   - If none of the above: call `createWorktree(branch, remote, existing)`

4. **Action Execution**: The routed function performs its specific operations:
   - `createWorktree()`: Creates the branch if needed, creates the worktree directory, and starts a bash shell in it
   - `deleteWorktree()`: Removes the worktree using git worktree remove, and cleans up empty parent directories
   - `mergeOrDelete()`: In dry-run mode validates all operations; with --confirm executes the merge/delete and cleanup

### Clone Command Flow

When a user runs `worktree clone https://github.com/user/repo.git`:

1. **Clone Command Validation**:
   - Validate that a repository URL argument was provided
   - Get the current working directory as the base for the new project
   - Extract the project name from the URL using `extractProjectName(url)` function

2. **Directory Setup**:
   - Create the project directory structure using `os.MkdirAll(projectPath, 0755)`
   - Execute git clone into a temporary location within the project directory using `exec.Command("git clone", repoURL, projectName)` with Dir set to projectPath

3. **Branch Detection**:
   - Get the current branch of the cloned repository using `getCurrentBranch(clonedRepoPath)` which runs `exec.Command("git rev-parse --abbrev-ref HEAD")`
   - First tries to read from origin/HEAD to get the explicit default branch
   - Falls back to checking common branch names in order: main, master, trunk using `getMainBranch()`

4. **Directory Renaming**:
   - Rename the cloned repository directory from its default name to the detected branch name
   - Uses `os.Rename(clonedRepoPath, branchPath)` where:
     - from: `projectPath/projectName` (the cloned repo directory)
     - to: `projectPath/branchName` (renamed to the branch name)

5. **Success Output**: Prints "Successfully cloned <url> into <path>" confirming the operation

---

## Source Code Structure

### File Organization

The entire application is contained in a single file: `cmd/worktree/main.go`

```
cmd/worktree/main.go
├── Package and Imports (Lines 1-11)
│   ├── Standard library imports
│   │   ├── fmt, os, exec, path/filepath, strings
│   │   └── ...
│   └── Third-party imports
│       └── github.com/spf13/cobra
│
├── Constants (Lines 13-20)
│   └── ANSI color codes (colorReset, colorRed, colorGreen, colorYellow, colorBlue)
│
├── Color Helper Functions (Lines 22-37)
│   ├── red(text string) string
│   ├── green(text string) string
│   ├── yellow(text string) string
│   └── blue(text string) string
│
├── Global Variables (Lines 39-97)
│   ├── rootCmd *cobra.Command
│   ├── Flags: remoteFlag, existingFlag, deleteFlag, mergeMode, deleteMode, confirmFlag
│   └── listRemoteFlag bool
│
├── init() Function (Lines 99-226)
│   ├── Root command configuration
│   ├── Global flags setup
│   ├── Subcommands registration
│   │   ├── listCmd
│   │   ├── cloneCmd
│   │   └── completionCmd
│   └── Main Run function assignment
│
├── main() Function (Lines 354-359)
│   └── Execute rootCmd
│
├── Utility Functions (Lines 361-952)
│   ├── isWorktree() bool
│   ├── isGitRepo() bool
│   ├── getMainBranch() (string, error)
│   ├── getDefaultBranchName() string
│   ├── branchExists(branch string) bool
│   ├── getLocalBranches() ([]string, error)
│   ├── getRemoteBranches() ([]string, error)
│   ├── branchNameCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)
│   ├── mergeOrDelete(branch string, mergeMode, deleteMode, confirm bool)
│   ├── deleteWorktree(branch string)
│   ├── createWorktree(branch string, remote, existing bool)
│   ├── extractProjectName(url string) string
│   └── getCurrentBranch(repoPath string) (string, error)
│
└── Command Functions (Embedded in init())
    ├── listCmd.Run
    ├── cloneCmd.Run
    └── completionCmd.RunE
```

### Line Count by Component

| Component | Lines | Percentage |
|-----------|-------|------------|
| Imports and Package | 11 | 1.3% |
| Color Constants and Helpers | 27 | 3.1% |
| Global Variables | 59 | 6.9% |
| init() Function | 128 | 14.9% |
| main() Function | 5 | 0.6% |
| Command Execution Logic | ~270 | 31.4% |
| Utility Functions | ~340 | 34.9% |
| Completion System | 54 | 6.3% |
| **Total** | **~910** | **100%** |

---

## Command Hierarchy

### Command Tree Structure

```
worktree (root)
├── init <project-name>
│   └── Initializes new project with worktree structure
│
├── clone <repository-url>
│   └── Validates and clones repository into structured directory
│
├── list [flags]
│   └── Lists branches (local or remote)
│       └── -r, --remote: List remote branches
│
├── <branch> [flags]
│   ├── (no flags): Create new branch and worktree
│   ├── -r, --remote: Create worktree from remote branch
│   ├── -e, --existing: Create worktree from existing local branch
│   └── -d, --delete-worktree: Delete worktree directory
│
├── --merge [flags] <branch>
│   ├── (no --confirm): Dry-run merge validation
│   └── --confirm: Execute merge and cleanup
│
├── --delete [flags] <branch>
│   ├── (no --confirm): Dry-run delete validation
│   └── --confirm: Execute deletion of branch and worktree
│
└── completion <shell-type>
    ├── bash: Generate Bash completion script
    ├── zsh: Generate Zsh completion script
    ├── fish: Generate Fish completion script
    └── powershell: Generate PowerShell completion script
```

### Flag Hierarchy and Precedence

```
Global Flags (applicable to root command):
├── -r, --remote: Create from remote branch
├── -e, --existing: Create from existing local branch
├── -d, --delete-worktree: Delete worktree directory
├── --merge: Merge branch into main
├── --delete: Delete branch and worktree
└── --confirm: Confirm destructive operation

List Command Flags (applicable to list subcommand):
└── -r, --remote: List remote branches

Mutual Exclusions:
├── -r and -e cannot be used together
├── -r, -e, -d cannot be used with --merge or --delete
├── --merge and --delete cannot be used together
└── --confirm can only be used with --merge or --delete
```

---

## Core Components

### 1. Cobra Command Framework

**Purpose**: Provides CLI command structure, argument parsing, and help system.

**Key Features Used**:
- Command tree with subcommands
- Flag parsing with short/long forms
- Built-in help system
- Argument validation
- Shell completion support

**Integration Points**:
- `cobra.Command` for command definitions
- `PersistentFlags()` for global flags
- `Flags()` for command-specific flags
- `Args` field for argument validation
- `Run` and `RunE` fields for command execution
- `ValidArgsFunction` for autocompletion

### 2. Git Integration Layer

**Purpose**: Executes Git commands and processes their output.

**Methods Used**:

| Git Command | Go Implementation | Purpose |
|-------------|-------------------|---------|
| `git init` | `exec.Command("git", "init")` | Initialize new repository |
| `git config` | `exec.Command("git", "config", "--get", "init.defaultBranch")` | Get default branch configuration |
| `git branch` | `exec.Command("git", "branch", ...)` | List branches |
| `git worktree add` | `exec.Command("git", "worktree", "add", ...)` | Create worktree |
| `git worktree remove` | `exec.Command("git", "worktree", "remove", ...)` | Remove worktree |
| `git merge` | `exec.Command("git", "merge", ...)` | Merge branches |
| `git fetch` | `exec.Command("git", "fetch", ...)` | Fetch remote branches |
| `git branch -d` | `exec.Command("git", "branch", "-d", ...)` | Delete local branch |
| `git push -d` | `exec.Command("git", "push", "-d", ...)` | Delete remote branch |
| `git symbolic-ref` | `exec.Command("git", "symbolic-ref", ...)` | Get symbolic ref |
| `git show-ref` | `exec.Command("git", "show-ref", ...)` | Check ref existence |
| `git rev-parse` | `exec.Command("git", "rev-parse", ...)` | Get current branch |
| `git log` | `exec.Command("git", "log", ...)` | Check commit history |

**Error Handling**:
- All Git commands are executed with `exec.Command()`
- Errors are captured and wrapped with descriptive messages
- Exit codes are preserved where appropriate
- User-friendly error messages with "WARNING:" prefix

### 3. Repository State Validation

**Purpose**: Ensure commands are run from valid contexts.

**Key Functions**:

```go
// Check if current directory is a worktree
func isWorktree() bool {
    info, err := os.Stat(".git")
    if err == nil && !info.IsDir() {
        return true  // .git is a file (worktree)
    }
    return false
}

// Check if current directory is the main git repository
func isGitRepo() bool {
    info, err := os.Stat(".git")
    if err != nil {
        return false
    }
    return info.IsDir()  // .git is a directory (main repo)
}
```

**Validation Logic**:
- Most commands require being in the main repository
- Clone command can run from any directory
- Completion command doesn't require a Git repository
- Validation happens in `PersistentPreRun` and individual command `Run` functions

### 4. Directory Structure Management

**Purpose**: Create and manage the organized directory structure.

**Key Functions**:

```go
// Create project directory structure
projectPath := filepath.Join(currentDir, projectName)
os.MkdirAll(projectPath, 0755)

// Clone repository
cloneCmd := exec.Command("git", "clone", repoURL, projectName)
cloneCmd.Dir = projectPath

// Rename to branch name
os.Rename(clonedRepoPath, branchPath)

// Create worktree at relative path
worktreePath := filepath.Join("..", branch)
cmd := exec.Command("git", "worktree", "add", worktreePath, branch)

// Cleanup empty directories (recursive)
cleanupPath := filepath.Dir(worktreePath)
for cleanupPath != "." && cleanupPath != ".." {
    // Check if directory is empty
    entries, err := os.ReadDir(cleanupPath)
    if err == nil && len(entries) == 0 {
        os.Remove(cleanupPath)
    }
    cleanupPath = filepath.Dir(cleanupPath)
}
```

---

## Key Technical Decisions

### 1. Single File Architecture

**Decision**: Entire application in one file (`main.go`)

**Rationale**:
- **Pros**: Simple deployment, no import paths, easier dependency management
- **Cons**: Harder to maintain as code grows, limited code organization
- **Justification**: Tool is small (~860 lines), manageable in single file

**Future Consideration**: Split into multiple files if complexity grows:
- `cmd/clone.go` - Clone command
- `cmd/worktree.go` - Worktree management
- `cmd/list.go` - List command
- `utils/git.go` - Git integration helpers
- `utils/validation.go` - Validation functions

### 2. Cobra Framework Selection

**Decision**: Use Cobra for CLI framework

**Alternatives Considered**:
- `flag` package (standard library)
- `urfave/cli`
- `kingpin`

**Rationale**:
- **Pros**: 
  - Built-in subcommand support
  - Automatic help generation
  - Shell completion support
  - Flag inheritance
  - Well-maintained and popular
- **Cons**: Slightly more complex than standard `flag`
- **Justification**: Need for subcommands and completion made Cobra the best choice

### 3. Worktree Relative Path Strategy

**Decision**: Create worktrees at `../<branch>` relative to main repository

**Rationale**:
- **Pros**:
  - Creates organized, nested directory structure
  - Groups branches by type (feat/, fix/, docs/)
  - Visually represents project state
  - Easy to navigate and understand
- **Cons**:
  - Requires being in main repository to run commands
  - Slightly more complex path handling
- **Justification**: The organized directory structure is the tool's main value proposition

### 4. Dry-Run by Default

**Decision**: Destructive operations default to dry-run mode

**Rationale**:
- **Pros**:
  - Safety first approach
  - Prevents accidental data loss
  - Allows validation before execution
  - Comprehensive feedback on what will happen
- **Cons**:
  - Requires extra step (`--confirm`) for execution
  - Slightly more verbose output
- **Justification**: Safety is paramount for destructive operations

### 5. Shell Integration via exec.Command

**Decision**: Use `exec.Command()` instead of Go Git libraries

**Rationale**:
- **Pros**:
  - Uses system Git (same version user has installed)
  - Supports all Git features and versions
  - No additional dependencies
  - Consistent with user's Git configuration
- **Cons**:
  - More error-prone (string-based command construction)
  - Requires careful argument escaping
  - Platform-dependent behavior
- **Justification**: Maximizes compatibility and reduces dependencies

### 6. Color-Coded Output

**Decision**: Implement ANSI color coding for user feedback

**Rationale**:
- **Pros**:
  - Improves user experience
  - Quick visual feedback on operation status
  - Consistent with modern CLI tools
- **Cons**:
  - May not work on all terminals
  - Can be disabled if needed
- **Justification**: Enhanced usability outweighs compatibility concerns

### 7. Bash Shell Assumption

**Decision**: Start bash shell in new worktrees

**Rationale**:
- **Pros**:
  - Consistent user experience
  - Bash is widely available
  - Supports all expected shell features
- **Cons**:
  - Not ideal for non-Bash users
  - Platform-dependent (Linux/macOS)
- **Justification**: Bash is the most common shell and the tool is primarily for Unix-like systems

**Future Consideration**: Make shell configurable via environment variable or flag:
```bash
worktree --shell zsh feat/my-branch
```

---

## Integration with Git

### Git Worktree Feature Utilization

The tool leverages Git's built-in worktree functionality:

**Git Worktree Commands Used**:
```bash
# Add a worktree
git worktree add <path> <branch>

# List worktrees
git worktree list

# Remove a worktree
git worktree remove <path>

# Prune stale worktrees
git worktree prune
```

### Worktree Detection

Git distinguishes between main repository and worktrees:

| Location | `.git` | Type |
|----------|--------|------|
| Main repository | Directory | Main repository |
| Worktree | File (symlink) | Worktree |

**Detection Code**:
```go
func isWorktree() bool {
    if _, err := os.Stat(".git"); err == nil {
        info, err := os.Stat(".git")
        if err == nil && !info.IsDir() {
            return true  // .git is a file, so this is a worktree
        }
    }
    return false
}
```

### Branch Management

The tool handles three types of branch operations:

1. **New Branch Creation**:
   ```go
   cmd := exec.Command("git", "branch", branch)
   ```

2. **Remote Branch Fetching**:
   ```go
   cmd := exec.Command("git", "fetch", "origin", fmt.Sprintf("%s:%s", branch, branch))
   ```

3. **Existing Branch Validation**:
   ```go
   cmd := exec.Command("git", "show-ref", "--quiet", "refs/heads/"+branch)
   return cmd.Run() == nil
   ```

### Merge Operations

**Merge Strategy**:
- Uses `git merge` with `--no-commit --no-ff` for dry-run testing
- Uses `git merge` without flags for actual execution
- Aborts merge with `git merge --abort` if testing fails

**Conflict Detection**:
- Tests merge in dry-run mode first
- Captures merge output and exit code
- Provides detailed feedback on merge failures

### Default Branch Detection

**Detection Order**:
1. Check `refs/remotes/origin/HEAD` for explicit default
2. Check for common branch names: `main`, `master`, `trunk`
3. Return error if none found

**Code**:
```go
func getMainBranch() (string, error) {
    // Try origin/HEAD first
    cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
    output, err := cmd.Output()
    if err == nil {
        ref := strings.TrimSpace(string(output))
        parts := strings.Split(ref, "/")
        if len(parts) >= 4 {
            return parts[3], nil
        }
    }
    
    // Fallback to common branches
    for _, candidate := range []string{"main", "master", "trunk"} {
        cmd := exec.Command("git", "show-ref", "--quiet", "refs/heads/"+candidate)
        if err := cmd.Run(); err == nil {
            return candidate, nil
        }
    }
    
    return "", fmt.Errorf("cannot determine main branch")
}
```

---

## Shell Autocompletion System

### Cobra Integration

Cobra provides built-in shell completion support through:
- `ValidArgsFunction` field on commands
- `ShellCompDirective` for completion behavior
- Automatic generation of completion scripts

### Branch Name Completion

**Implementation**:
```go
func branchNameCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    // Handle argument context
    if len(args) > 0 && !strings.HasPrefix(args[len(args)-1], "-") {
        return nil, cobra.ShellCompDirectiveNoFileComp
    }
    
    // Get all branches (local and remote)
    seen := make(map[string]bool)
    var allBranches []string
    
    // Add local branches
    localBranches, err := getLocalBranches()
    if err == nil {
        for _, b := range localBranches {
            if !seen[b] {
                seen[b] = true
                allBranches = append(allBranches, b)
            }
        }
    }
    
    // Add remote branches (without origin/ prefix)
    remoteBranches, err := getRemoteBranches()
    if err == nil {
        for _, b := range remoteBranches {
            if !seen[b] {
                seen[b] = true
                allBranches = append(allBranches, b)
            }
        }
    }
    
    // Filter by prefix
    var completions []string
    for _, branch := range allBranches {
        if strings.HasPrefix(branch, toComplete) {
            completions = append(completions, branch)
        }
    }
    
    return completions, cobra.ShellCompDirectiveNoFileComp
}
```

### Completion Features

1. **Context-Aware**: Only completes branch names for positional arguments
2. **Deduplication**: Removes duplicate branch names (local vs remote)
3. **Prefix Filtering**: Only shows branches matching current input
4. **Multi-Source**: Combines local and remote branches
5. **Shell-Specific**: Generates appropriate scripts for each shell type

### Completion Script Generation

The `completion` command uses Cobra's built-in methods:

```go
completionCmd.RunE = func(cmd *cobra.Command, args []string) error {
    switch args[0] {
    case "bash":
        return rootCmd.GenBashCompletion(os.Stdout)
    case "zsh":
        return rootCmd.GenZshCompletion(os.Stdout)
    case "fish":
        return rootCmd.GenFishCompletion(os.Stdout, true)
    case "powershell":
        return rootCmd.GenPowerShellCompletion(os.Stdout)
    default:
        return fmt.Errorf("unsupported shell: %s", args[0])
    }
}
```

### Completion Integration

**How it works**:
1. User types partial command and presses Tab
2. Shell sends completion request to worktree
3. Cobra calls `ValidArgsFunction` for the current command
4. Function returns list of possible completions
5. Shell displays matching completions

**Example Flow**:
```
User types: worktree feat/<TAB>
Shell sends: worktree __completeNoDesc feat/
Cobra calls: branchNameCompletion(cmd, ["feat/"], "feat/")
Function returns: ["feat/login", "feat/payment", "feat/auth"]
Shell displays: feat/login  feat/payment  feat/auth
```

---

## Color Coding System

### ANSI Color Constants

```go
const (
    colorReset  = "\033[0m"
    colorRed    = "\033[31m"
    colorGreen  = "\033[32m"
    colorYellow = "\033[33m"
    colorBlue   = "\033[34m"
)
```

### Color Helper Functions

```go
func red(text string) string {
    return colorRed + text + colorReset
}

func green(text string) string {
    return colorGreen + text + colorReset
}

func yellow(text string) string {
    return colorYellow + text + colorReset
}

func blue(text string) string {
    return colorBlue + text + colorReset
}
```

### Color Usage Semantics

| Color | Usage | Example |
|-------|-------|---------|
| **Green** | Success, PASS | `green("PASS")` - operation succeeded |
| **Red** | Failure, ERROR | `red("FAIL")` - operation failed |
| **Yellow** | Warning, CAUTION | `yellow("WARNING")` - user attention needed |
| **Blue** | Information, STATUS | `blue("INFO")` - informational message |

### Color Output Examples

```
=== DRY RUN - Testing if operations are possible ===

Checking: Branch feat/test status relative to main
  INFO: Branch feat/test has 2 commit(s) ahead of main
  PASS: Branch feat/test is up to date with main
  PASS: Merge is possible

Testing: Delete worktree at ../feat/test
  PASS: Worktree can be removed

Testing: Delete local branch feat/test
  WARNING: Local branch can be deleted (with force -D)

Testing: Delete remote branch origin/feat/test
  PASS: Remote branch can be deleted

All operations are possible!
```

### Terminal Compatibility

- **ANSI Support**: Works on most Unix-like terminals
- **Windows**: May not display colors without ANSI support
- **Color Disabling**: Colors can be disabled by setting `NO_COLOR` environment variable (standard convention)

---

## Error Handling Strategy

### Error Classification

The tool categorizes errors into different types with appropriate handling:

| Error Type | Handling | User Message |
|------------|----------|--------------|
| **Validation Error** | Immediate exit with message | "WARNING: [description]" |
| **Git Command Error** | Wrapped with context | "WARNING: [git operation] failed: [error]" |
| **File System Error** | Wrapped with context | "WARNING: failed to [operation]: [error]" |
| **User Input Error** | Help message | Shows usage and error |

### Error Handling Patterns

**1. Validation Errors (Early Exit)**:
```go
if isWorktree() {
    fmt.Fprintln(os.Stderr, "WARNING: must be called from the main repository, not a worktree")
    os.Exit(1)
}
```

**2. Git Command Errors**:
```go
cmd := exec.Command("git", "branch", branch)
if err := cmd.Run(); err != nil {
    fmt.Fprintf(os.Stderr, "WARNING: failed to create branch: %v\n", err)
    os.Exit(1)
}
```

**3. File System Errors**:
```go
if err := os.MkdirAll(projectPath, 0755); err != nil {
    fmt.Fprintf(os.Stderr, "Error: failed to create directory: %v\n", err)
    os.Exit(1)
}
```

### Error Recovery Strategy

**Dry-Run Mode**: Tests operations without making changes:
```go
// Test merge without committing
cmd := exec.Command("git", "merge", "--no-commit", "--no-ff", branch)
if err := cmd.Run(); err != nil {
    // Merge would fail - report error
    fmt.Printf("  %s: Merge would fail: %v\n", red("FAIL"), err)
    allPossible = false
} else {
    // Reset the merge attempt
    cmd := exec.Command("git", "merge", "--abort")
    cmd.Run() // Ignore error
    fmt.Printf("  %s: Merge is possible\n", green("PASS"))
}
```

**Cleanup on Failure**:
```go
// In clone command
cloneCmd := exec.Command("git", "clone", repoURL, projectName)
cloneCmd.Dir = projectPath
if err := cloneCmd.Run(); err != nil {
    fmt.Fprintf(os.Stderr, "Error: failed to clone repository: %v\n", err)
    // Clean up the directory we created
    os.RemoveAll(projectPath)
    os.Exit(1)
}
```

### Exit Codes

| Exit Code | Meaning | Usage |
|-----------|---------|-------|
| 0 | Success | Normal operation completion |
| 1 | Error | Most failures |

**Rationale**: Simple exit code system sufficient for CLI tool usage.

---

## Validation and Safety Mechanisms

### Input Validation

**1. Argument Count Validation**:
```go
cloneCmd.Args = func(cmd *cobra.Command, args []string) error {
    if len(args) < 1 {
        return fmt.Errorf("requires a repository URL argument")
    }
    return nil
}
```

**2. Flag Combination Validation**:
```go
if remoteFlag && existingFlag {
    fmt.Fprintln(os.Stderr, "WARNING: cannot use switches -r and -e at the same time")
    os.Exit(1)
}

if (mergeMode || deleteMode) && (remoteFlag || existingFlag || deleteFlag) {
    fmt.Fprintln(os.Stderr, "WARNING: cannot use -r, -e, or -d with --merge or --delete")
    os.Exit(1)
}
```

**3. Context Validation**:
```go
if isWorktree() {
    fmt.Fprintln(os.Stderr, "WARNING: must be called from the main repository, not a worktree")
    os.Exit(1)
}

if !isGitRepo() {
    fmt.Fprintln(os.Stderr, "WARNING: not a git repository")
    os.Exit(1)
}
```

### Dry-Run Validation

The `mergeOrDelete()` function performs comprehensive validation:

**Validation Checks**:
1. **Branch Existence**: Verify branch exists locally
2. **Main Branch Detection**: Determine main branch name
3. **Commit Status**: Check if branch is ahead/behind main
4. **Merge Test**: Attempt no-commit merge to validate
5. **Worktree Existence**: Verify worktree can be removed
6. **Branch Merged Status**: Check if branch is fully merged
7. **Remote Branch Status**: Check remote branch deletion possibility

**Dry-Run Output**:
- Green (PASS): Operation will succeed
- Red (FAIL): Operation will fail
- Yellow (WARNING): Operation may have issues but can proceed
- Blue (INFO): Informational messages

### Confirmation Requirement

**Safety Principle**: Destructive operations require explicit confirmation.

```go
if confirm {
    // Execute the operation
    performMergeOrDelete(branch, mergeMode, deleteMode)
} else {
    // Dry-run: validate and show what would happen
    mergeOrDelete(branch, mergeMode, deleteMode, false)
    fmt.Printf("To execute, run: worktree %s --confirm %s\n", flagStr, branch)
}
```

### Path Validation

**1. Worktree Path Construction**:
```go
worktreePath := filepath.Join("..", branch)
```

**2. Absolute Path Resolution**:
```go
absWorktreePath, err := filepath.Abs(worktreePath)
if err != nil {
    absWorktreePath = worktreePath
}
```

**3. Path Existence Check**:
```go
if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
    // Worktree doesn't exist
}
```

---

## Performance Considerations

### Git Command Execution

**Performance Factors**:
- **Repository size**: Larger repos = slower operations
- **Branch count**: More branches = slower listing
- **Commit history**: Long history = slower merge validation
- **Network**: Remote operations depend on network speed

**Optimizations Implemented**:
1. **Deduplication**: Branch name completion removes duplicates
2. **Prefix Filtering**: Only returns branches matching input
3. **Early Validation**: Fails fast on validation errors
4. **Parallel Operations**: Not currently implemented (future consideration)

### Memory Usage

- **Low Memory Footprint**: ~10-20MB typical usage
- **No Caching**: All data fetched fresh (no stale cache issues)
- **Temporary Data**: Minimal temporary data structures

### I/O Operations

**Disk I/O**:
- Worktree creation: Medium I/O (git clone, directory creation)
- Worktree deletion: Medium I/O (directory removal)
- Branch operations: Low I/O (git commands)

**Network I/O**:
- Clone operations: High I/O (full repository download)
- Remote branch operations: Medium I/O (fetch operations)
- Local operations: No network I/O

### Scalability Limits

| Repository Characteristic | Performance Impact | Mitigation |
|--------------------------|-------------------|------------|
| Large file count | Slow clone | Use shallow clone (future) |
| Many branches | Slow completion | Prefix filtering helps |
| Long history | Slow merge validation | Commit range limits (future) |
| Large files | High memory usage | Git LFS support (future) |

---

## Security Considerations

### Command Injection Protection

**Risk**: User-provided input used in shell commands.

**Mitigation**:
- Use `exec.Command()` with separate arguments (not shell strings)
- Avoid `exec.Command("git " + userInput)` - use `exec.Command("git", userInput)`
- Validate all user inputs

**Safe Example**:
```go
// SAFE: Arguments passed separately
cmd := exec.Command("git", "branch", branchName)

// UNSAFE: String concatenation (command injection risk)
cmd := exec.Command("git branch " + branchName)  // ❌ Don't do this
```

### Path Traversal Protection

**Risk**: User-provided branch names could contain path traversal sequences.

**Mitigation**:
- Use `filepath.Join()` for path construction (handles traversal)
- Validate branch names match expected patterns
- Use relative paths within repository

**Safe Path Construction**:
```go
// SAFE: Uses filepath.Join which handles traversal
worktreePath := filepath.Join("..", branchName)

// UNSAFE: Direct string concatenation
treePath := "../" + branchName  // ⚠️ Risk if branchName contains ".."
```

### File System Operations

**Risk**: Accidental file deletion or modification.

**Mitigation**:
- Dry-run mode for destructive operations
- Explicit confirmation required
- Only remove empty directories (safe deletion)
- Validate paths before operations

### Environment Security

**Risk**: Sensitive information in environment or commands.

**Mitigation**:
- No hardcoded credentials
- No logging of sensitive operations
- Use Git's built-in credential handling

---

## Extensibility and Future Architecture

### Potential Future Enhancements

| Feature | Implementation Complexity | Value |
|---------|-------------------------|-------|
| **Shallow Clone Support** | Low | Faster cloning for large repos |
| **Custom Worktree Paths** | Medium | User-defined directory structures |
| **Shell Selection** | Low | Choose shell for worktree (bash, zsh, etc.) |
| **Windows Support** | Medium | Full Windows compatibility |
| **GUI Integration** | High | Graphical interface for worktree management |
| **Configuration File** | Medium | Persistent settings and preferences |
| **Multiple Main Branches** | Medium | Support for different default branch names |
| **Worktree Templates** | Low | Pre-configured worktree setups |

### Refactoring Opportunities

**1. Code Organization**:
```
cmd/
├── worktree/
│   ├── main.go          # Entry point and command setup
│   ├── clone.go         # Clone command implementation
│   ├── create.go        # Worktree creation logic
│   ├── delete.go        # Deletion and cleanup logic
│   ├── list.go          # Branch listing
│   ├── merge.go         # Merge operations
│   └── utils/
│       ├── git.go       # Git command wrappers
│       ├── validation.go # Input validation
│       ├── paths.go     # Path manipulation
│       └── colors.go    # Color coding utilities
```

**2. Interface Abstraction**:
```go
// GitClient interface for testability
type GitClient interface {
    Branch(string) error
    WorktreeAdd(string, string) error
    Merge(string) error
    // ... etc
}

// RealGitClient implementation
type RealGitClient struct{}

// MockGitClient for testing
type MockGitClient struct{}
```

**3. Configuration System**:
```go
// Config struct
type Config struct {
    DefaultShell string
    MainBranches []string
    AutoCleanup  bool
    // ... etc
}

// Load from file
func LoadConfig() (*Config, error) { ... }
```

### Testing Infrastructure

**Current State**: No dedicated test suite

**Future Testing Approach**:
```go
// Table-driven tests for utility functions
func TestExtractProjectName(t *testing.T) {
    tests := []struct {
        url      string
        expected string
    }{
        {"https://github.com/user/repo.git", "repo"},
        {"git@github.com:user/repo.git", "repo"},
        // ... etc
    }
    
    for _, test := range tests {
        got := extractProjectName(test.url)
        if got != test.expected {
            t.Errorf("extractProjectName(%q) = %q, want %q", test.url, got, test.expected)
        }
    }
}

// Mock-based tests for command logic
func TestCreateWorktree(t *testing.T) {
    mockGit := &MockGitClient{}
    // Setup mock expectations
    
    err := createWorktree("test-branch", false, false)
    if err != nil {
        t.Errorf("createWorktree failed: %v", err)
    }
    
    // Verify mock was called correctly
    mockGit.AssertExpectations(t)
}
```

### Performance Improvements

**1. Caching**:
- Cache branch list for completion (with TTL)
- Cache remote branch information
- Cache repository state

**2. Parallel Execution**:
- Parallel worktree operations
- Concurrent Git command execution
- Background validation

**3. Lazy Loading**:
- Load branch information on demand
- Deferred repository analysis

---

## Summary

The worktree tool follows a **pragmatic, safety-first architecture** with the following key characteristics:

### Architectural Principles

✅ **Simplicity First**: Single file, minimal dependencies, straightforward logic  
✅ **Safety by Default**: Dry-run mode, explicit confirmation, comprehensive validation  
✅ **Git Native**: Uses Git's built-in features rather than reimplementing them  
✅ **User-Centric**: Clear feedback, helpful errors, intuitive workflows  
✅ **Extensible**: Designed for future growth and enhancement  

### Technical Stack Summary

- **Language**: Go 1.21+ (compiled, statically linked binaries)
- **Framework**: Cobra CLI (mature, well-supported)
- **Git Integration**: Native Git CLI (maximum compatibility)
- **Shell Integration**: Standard shells (Bash, Zsh, Fish, PowerShell)
- **Platform**: Unix-like systems (Linux, macOS)

### Code Quality

- **Readability**: Clear function and variable names
- **Error Handling**: Consistent error messages with context
- **Validation**: Comprehensive input and state validation
- **Organization**: Logical grouping of related functionality
- **Documentation**: Inline comments for complex logic

### Future Directions

The current architecture supports several evolution paths:
- **Code Organization**: Split into multiple files as complexity grows
- **Testing**: Add comprehensive test suite
- **Features**: Extend functionality while maintaining simplicity
- **Performance**: Optimize for large repositories and teams
- **Compatibility**: Broaden platform and shell support

For more information, see:
- [README.md](./README.md) - User documentation and usage examples
- [USAGE.md](./USAGE.md) - Comprehensive usage guide
- [WORKFLOW.md](./WORKFLOW.md) - Development workflow patterns
- [AGENTS.md](./AGENTS.md) - AI agent guidance