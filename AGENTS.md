# AGENTS.md

This file provides guidance for AI agents (like Mistral Vibe) when working with the worktree repository.

## AGENTS.md Maintenance (Living Document)

**This is a LIVING DOCUMENT** maintained by AI agents. You MUST update this file based on user feedback and new learnings.

### Agent Responsibilities

- **MUST**: Update this file after every user interaction that provides new information or feedback
- **MUST**: Record all user preferences, constraints, and learnings
- **MUST**: Remove outdated, incorrect, or superseded information
- **MUST**: Review this file at the start of each session to understand current requirements
- **SHOULD**: Add entries to the Update Log when making changes

### Purpose

This file serves as:
- **Persistent memory** of user preferences and requirements across agent sessions
- **Process documentation** for how agents should work in this repository
- **Knowledge base** of what works, what doesn't, and what the user expects
- **Feedback loop** to continuously improve agent behavior

### Update Process

When the user provides feedback:
1. Analyze what the user is communicating (preference, constraint, correction, etc.)
2. Update the relevant section of this AGENTS.md file
3. Add an entry to the Update Log at the bottom of this file
4. Commit the changes with a descriptive message

## Repository Overview

**Name**: worktree  
**Purpose**: A Go-based CLI tool that wraps git to provide efficient worktree management and repository cloning with structured directory organization.  
**Primary Use Case**: Developer productivity tool for managing multiple Git branches simultaneously using worktrees.  
**Language**: Go 1.21+  
**Framework**: [Cobra](https://github.com/spf13/cobra) CLI framework

## Repository Purpose and Scope

The worktree tool solves the following problems:
- **Complex branch management**: Manages multiple feature/bugfix branches simultaneously
- **Repository organization**: Clones repositories into predictable directory structures (`project_name/[branch]/`)
- **Worktree automation**: Creates and manages Git worktrees at `../[branch]` relative to the main repository
- **Safety and validation**: Provides dry-run capabilities for destructive operations (merge, delete)
- **Developer experience**: Shell autocompletion for branch names and intelligent error handling

## Technical Stack

- **Language**: Go 1.21 (as specified in go.mod)
- **CLI Framework**: Cobra v1.7.0
- **Dependencies**: 
  - `github.com/spf13/cobra` - CLI command structure
  - `github.com/spf13/pflag` - Flag parsing (indirect)
  - `github.com/inconshreveable/mousetrap` - Cross-platform trap handling (indirect)
- **Build System**: Standard Go modules
- **Installation**: Custom bash scripts (`install.bash`, `uninstall.bash`)

## Build and Development Procedures

### Building the Tool

```bash
# From repository root
go build -o worktree ./cmd/worktree

# Or use the provided install script
./install.bash
```

The install script:
- Builds the Go binary
- Installs to `~/bin/` (ensure this is in your PATH)
- Creates necessary directory structure

### Testing

**Current State**: The repository does not have a dedicated test suite. Agents should:
- Test manually by running commands in a safe environment
- Use `--dry-run` or `--help` flags where available
- Validate functionality against the README.md examples

**Recommended Testing Approach**:
1. Create a temporary test repository
2. Run worktree commands and verify output
3. Check directory structure matches expected patterns
4. Test edge cases (non-existent branches, invalid URLs, etc.)

### Common Development Commands

```bash
# Build
go build ./cmd/worktree

# Run directly
go run ./cmd/worktree --help

# Install completion
go run ./cmd/worktree completion bash > /tmp/worktree_completion

# Clean
go clean
```

### Debug Build Cleanup

**Important**: Always clean up debug builds before committing:

```bash
# Remove debug binary from repository root
rm -f worktree

# Or use git clean to remove all untracked files
git clean -fd
```

**Agent Responsibilities**:
- **MUST** remove debug binaries (`worktree`) from the repository before committing
- **MUST** check `git status` before committing to ensure no unintended files are included
- **SHOULD** use `git clean -fd` or explicitly remove debug builds after testing

## Code Structure and Conventions

```
.
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums
├── README.md           # Main user documentation
├── AGENTS.md           # This file
├── install.bash        # Installation script
├── uninstall.bash      # Uninstallation script
├── install_completion.bash  # Completion installation helper
└── cmd/
    └── worktree/
        └── main.go      # Main application source
```

### Source Code Organization

**`cmd/worktree/main.go`** contains the entire application logic:
- **Lines 1-20**: Imports and color constant definitions
- **Lines 22-84**: Color helper functions and root command setup
- **Lines 86-97**: Global flag definitions
- **Lines 99-226**: Command initialization and subcommand definitions
- **Lines 228-306**: Main run function with argument parsing and validation
- **Lines 309-352**: Completion command for shell integration
- **Lines 354-860**: Helper functions and core logic

### Naming Conventions

- **Commands**: lowercase, hyphen-separated (e.g., `worktree clone`)
- **Flags**: single letter short forms with long form equivalents (e.g., `-r, --remote`)
- **Functions**: camelCase for Go functions (e.g., `createWorktree`, `deleteWorktree`)
- **Variables**: camelCase with descriptive names
- **Constants**: UPPER_CASE for color codes and configuration values

## Agent Constraints and Guidelines

### DO's

✅ **DO** use the `--help` flag to discover command capabilities  
✅ **DO** test destructive operations (`--merge`, `--delete`) in dry-run mode first  
✅ **DO** run commands from the main repository, not from worktree directories  
✅ **DO** validate that you're in a Git repository before running worktree commands  
✅ **DO** use the provided autocompletion for branch name discovery  
✅ **DO** respect the directory structure conventions (`project_name/[branch]/`)  
✅ **DO** check for existing worktrees before creating new ones  
✅ **DO** use appropriate error handling and user feedback  

### DON'Ts

❌ **DON'T** run destructive operations without `--confirm` flag in production environments  
❌ **DON'T** modify the `.git` directory directly - use git commands  
❌ **DON'T** assume branch names - always validate or use completion  
❌ **DON'T** run worktree commands from within a worktree directory  
❌ **DON'T** delete the main branch worktree  
❌ **DON'T** modify existing test repositories without explicit permission  
❌ **DON'T** push changes to protected branches (main, master, trunk) without approval  
❌ **DON'T** change the go.mod or dependency versions without testing  

## Git Workflow Expectations

The worktree tool follows and enforces these Git workflow patterns:

### Standard Development Flow

1. **Clone**: `worktree clone <repository-url>`
   - Creates `project_name/[default-branch]/` structure
   - Default branch is determined from origin/HEAD (main, master, or trunk)

2. **Create Feature Branch**: `worktree <feature-name>`
   - Creates new branch from current HEAD
   - Creates worktree at `../feature-name/`
   - Starts bash shell in the new worktree

3. **Create from Remote**: `worktree -r <branch-name>`
   - Creates local tracking branch from `origin/<branch-name>`
   - Creates worktree at `../<branch-name>/`

4. **Create from Existing Local**: `worktree -e <branch-name>`
   - Creates worktree from existing local branch
   - Does not create new branch

5. **List Branches**: `worktree list`
   - Lists all branches (local and remote)

6. **Merge and Cleanup**: `worktree --merge --confirm <branch-name>`
   - Merges branch into main/master
   - Removes worktree directory
   - Cleans up empty parent directories
   - Requires `--confirm` for execution

7. **Delete Everything**: `worktree --delete --confirm <branch-name>`
   - Deletes local branch
   - Deletes remote branch
   - Removes worktree directory
   - Requires `--confirm` for execution

8. **Delete Worktree Only**: `worktree -d <branch-name>`
   - Removes worktree directory only
   - Does not delete branches

### Branch Naming Conventions

The tool does not enforce branch naming, but follows these patterns:
- `feat/*` - Feature branches
- `fix/*` - Bug fix branches
- `docs/*` - Documentation updates
- `chore/*` - Maintenance tasks
- `refactor/*` - Code refactoring
- `test/*` - Test-related branches

## Directory Structure Conventions

The worktree tool creates and expects the following structure:

```
current-directory/
└── project_name/
    ├── main/               # Main branch (or master/trunk)
    │   ├── .git/           # Git repository (directory, not file)
    │   └── ...            # Project files
    ├── feat/
    │   ├── feature-a/      # Feature branch worktree
    │   │   └── ...         # Worktree files
    │   └── feature-b/
    │       └── ...
    ├── fix/
    │   └── bug-fix/        # Bug fix worktree
    │       └── ...
    └── docs/
        └── readme/         # Documentation worktree
            └── ...
```

**Key Insight**: You are working in a worktree checked out to `project/[branch name]`. The actual git repository is at `project/[master/main]` which contains the `.git` **directory**. Worktrees have `.git` as a **file** (symlink to main repo), while the main repo has `.git` as a **directory**. This distinction is critical for understanding the worktree tool's behavior.

## Agent Learnings and Guidelines

### Documentation Style

Based on user feedback, agents must follow these documentation guidelines:

✅ **DO**: 
- Use **text descriptions** for complex concepts, relationships, and flows
- Explain diagrams and architecture using clear, descriptive paragraphs
- Use **bullet points** and **numbered lists** for step-by-step processes
- Keep **simple tree structures** (├──, └──, │) as they provide useful hierarchy visualization

❌ **DON'T**: 
- Create **ASCII art** or **box diagrams** in documentation
- Use complex visual diagrams that rely on box-drawing characters (┌, ┐, ┘, │, ├, etc.)
- Assume readers can parse visual layouts - use text descriptions instead

### Git Workflow

✅ **DO**: 
- Create **conventional commits** with each **logical unit of work**
- Use descriptive commit messages that explain the what and why
- Make atomic commits that represent a single logical change
- Commit frequently to maintain a clean history

❌ **DON'T**: 
- Bundle unrelated changes into a single commit
- Use vague commit messages like "fixed stuff" or "updated files"

## Shell Autocompletion

The tool provides comprehensive shell autocompletion:

### Supported Shells
- Bash
- Zsh  
- Fish
- PowerShell

### Completion Features
- **Branch name completion**: Tab completion for branch names
- **Both local and remote branches**: Includes all available branches
- **Prefix filtering**: Only shows branches matching the current input
- **Deduplication**: Removes duplicate branch names
- **Context-aware**: Only completes branch names, not flags or other arguments

### Completion Command
```bash
# Show completion help
worktree completion --help

# Generate completion script
worktree completion bash  # or zsh, fish, powershell
```

## Error Handling and User Feedback

The tool uses color-coded output for different message types:

| Color | Meaning | Usage |
|-------|---------|-------|
| Green | PASS/SUCCESS | Successful operations, validation passed |
| Red | FAIL/ERROR | Failed operations, validation errors |
| Yellow | WARNING | Non-critical issues, user attention needed |
| Blue | INFO | Informational messages, status updates |

### Dry-Run Output Format

When using `--merge` or `--delete` without `--confirm`, the tool outputs:

```
=== DRY RUN - Testing if operations are possible ===

Checking: Branch <branch> status relative to <main>
  [COLOR]: [Status message]

Testing: [Operation description]
  [COLOR]: [Result message]

[Summary line]

To execute, run: worktree [--merge|--delete] --confirm <branch>
```

## Environment and Dependencies

### Required Tools
- **Git**: Must be installed and available in PATH
- **Go 1.21+**: For building the tool
- **Bash**: For installation scripts (though the tool itself is shell-agnostic)

### Environment Variables
The tool does not currently use environment variables, but agents should be aware of:
- `GOPATH`: May affect Go module behavior
- `GOBIN`: May affect where binaries are installed
- `PATH`: Must include `~/bin` for installed tool to be accessible

## Common Pitfalls and Solutions

### Issue: "WARNING: must be called from the main repository, not a worktree"
**Solution**: Run the command from the main repository directory (where `.git` is a directory, not a file)

### Issue: "WARNING: not a git repository"
**Solution**: Ensure you're in a directory with a `.git` directory or file, or use `worktree clone` to create a new repository

### Issue: "WARNING: cannot use -r and -e at the same time"
**Solution**: Choose either `-r` (for remote branch) or `-e` (for existing local branch), not both

### Issue: Branch completion not working
**Solution**: Ensure you're in a Git repository and the completion script is properly sourced

### Issue: Merge would fail due to conflicts
**Solution**: Run the dry-run first (`worktree --merge <branch>`) to identify conflicts, resolve them manually, then run with `--confirm`

## Code Review and Quality Guidelines

### For Agent-Generated Code Changes

When making changes to the codebase:

1. **Maintain existing style**: Match the existing code formatting and patterns
2. **Add appropriate error handling**: Use the established error message patterns
3. **Use color coding**: For user-facing output, use the existing color helper functions
4. **Validate inputs**: Check for nil/empty values and invalid states
5. **Test thoroughly**: Manually test changes before committing
6. **Document changes**: Update relevant documentation if functionality changes

### Code Quality Standards

- Follow Go idioms and best practices
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and single-purpose
- Handle edge cases gracefully
- Use consistent error messages (start with "WARNING:" or descriptive text)

## Security Considerations

⚠️ **Important Security Notes**:

- The tool executes shell commands using `exec.Command()` - be cautious of command injection
- Repository URLs are user-provided - validate before use
- Branch names are used in file paths - ensure they're sanitized
- The tool has access to the entire file system based on the current directory

**Agent Responsibilities**:
- Never hardcode sensitive information
- Validate all user inputs before use
- Be cautious with file system operations
- Respect file permissions and ownership

## Performance Considerations

- The tool makes multiple git calls for complex operations
- Dry-run mode can be expensive for repositories with many branches/commits
- Branch name completion may be slow in large repositories
- Worktree creation involves file system operations that may be I/O intensive

## Documentation Standards

When updating documentation:

1. **Keep examples accurate**: Test all code examples before including them
2. **Use consistent terminology**: Match existing README.md terminology
3. **Include all flags**: Document all available flags for each command
4. **Show both short and long forms**: Include both `-r` and `--remote` in examples
5. **Add warnings**: Note any destructive operations or potential pitfalls
6. **Cross-reference**: Link to related commands or documentation

## Agent-Specific Configuration

For Mistral Vibe and similar AI agents:

```yaml
# Suggested agent configuration for this repository
repo:
  name: worktree
  language: go
  framework: cobra
  
constraints:
  - "Never run destructive operations without --confirm in production"
  - "Always test in dry-run mode first"
  - "Validate git repository state before operations"
  - "Respect the directory structure conventions"
  
preferences:
  shell: bash
  editor: none  # Use direct file operations
  testing: manual  # No automated test suite
```

## Getting Help

If you're an AI agent and encounter issues:

1. **Check the existing code**: The main.go file contains all logic and examples
2. **Run with --help**: Each command has built-in help
3. **Test in isolation**: Create a test repository to experiment
4. **Review README.md**: Contains user-focused documentation
5. **Check dry-run output**: Provides detailed information about what would happen

## Version Information

- **Current Version**: As per git history, latest commit is `21816b6` (docs)
- **Go Version**: 1.21 (from go.mod)
- **Cobra Version**: v1.7.0
- **Repository**: github.com/lachlan/worktree

---

## Update Log

This section records changes made to AGENTS.md to maintain a history of learning and refinement.

| Date | Change | Based On | Commit |
|------|--------|----------|--------|
| 2026-06-30 | Initial AGENTS.md creation with comprehensive agent guidance | User request | 21816b6 |
| 2026-06-30 | Added living document instructions - AGENTS.md now self-maintaining | User feedback: "AGENTS.md should be a living document" | 72d3030 |
| 2026-06-30 | Added documentation style guidelines - no ASCII art, use text descriptions | User feedback: "get rid of those ascii boarders... put it in the AGENTS.md to not create those" | 72d3030 |
| 2026-06-30 | Added git workflow guidelines - conventional commits per logical work unit | User feedback: "add that it should create a conventional commit with each logical units of work" | 72d3030 |
| 2026-06-30 | Added worktree structure explanation - project/[branch name] with project/[master/main] as actual repo | User feedback: "Explain that they are in a worktree checked out to the project/[branch name] with project/[master/main] being the actual repo" | 72d3030 |
| 2026-06-30 | Replaced all ASCII diagrams in ARCHITECTURE.md with text descriptions | User feedback: "get rid of those ascii boarders... replace with text explaing it in detail" | 72d3030 |
| 2026-06-30 | Replaced all ASCII diagrams in WORKFLOW.md with text descriptions | User feedback: "get rid of those ascii boarders... replace with text explaing it in detail" | 72d3030 |
| 2026-06-30 | Added missing `worktree init` command documentation to USAGE.md, WORKFLOW.md, ARCHITECTURE.md | User feedback: "is the documentation up to date?" → No, init command was missing | 205b110 |
| 2026-07-19 | Added debug build cleanup guidelines - MUST remove debug binaries before committing | User feedback: "Take necessary measures to remember to delete debug builds" | de3cb7c |
| 2026-07-25 | Added ahead/behind commit check to --delete dry-run mode | User request: "make it so that it checks by how much the branch is ahead on the --delete flag dry-run" | HEAD |

---

*This AGENTS.md file provides comprehensive guidance for AI agents working with the worktree repository. This is a living document - agents MUST update it based on user feedback.*