# WORKFLOW.md

This document describes the desired development workflow when using the worktree tool. It provides step-by-step guides for common development scenarios and explains the philosophy behind the tool's design.

## Table of Contents

- [Worktree Philosophy](#worktree-philosophy)
- [Core Workflow Principles](#core-workflow-principles)
- [Getting Started Workflow](#getting-started-workflow)
- [Feature Development Workflow](#feature-development-workflow)
- [Bug Fix Workflow](#bug-fix-workflow)
- [Documentation Workflow](#documentation-workflow)
- [Team Collaboration Workflow](#team-collaboration-workflow)
- [Project Setup Workflow](#project-setup-workflow)
- [Cleanup and Maintenance Workflow](#cleanup-and-maintenance-workflow)
- [Advanced Workflows](#advanced-workflows)
- [Workflow Visualization](#workflow-visualization)
- [Best Practices by Scenario](#best-practices-by-scenario)

---

## Worktree Philosophy

The worktree tool is built around several core principles that guide its design and intended usage:

### 1. Directory Structure as Documentation

**Principle**: Your directory structure should tell the story of your project.

The tool creates organized directory structures that make it immediately obvious:
- What projects you're working on
- What branches exist
- What features are in development
- What bugs are being fixed

```
projects/
├── myapp/
│   ├── main/           # Production-ready code
│   ├── feat/
│   │   ├── login/      # User login feature in development
│   │   └── payment/    # Payment processing feature in development
│   └── fix/
│       └── security/   # Security vulnerability fix in progress
└── api/
    ├── master/         # Main API code
    └── feat/
        └── v2/         # API v2 development
```

### 2. Isolated Working Environments

**Principle**: Each worktree provides a completely isolated working environment.

- Changes in one worktree don't affect others
- Each worktree can be on a different branch
- Independent git state (index, working directory)
- No accidental cross-contamination between features

### 3. Safety First

**Principle**: Destructive operations should be safe and predictable.

- Dry-run mode for all destructive operations
- Explicit confirmation required for execution
- Comprehensive validation before execution
- Clear, color-coded feedback

### 4. Developer Productivity

**Principle**: Reduce cognitive overhead and friction in development.

- Shell autocompletion for branch names
- Automatic directory creation and management
- Minimal command syntax
- Intelligent defaults

### 5. Team Collaboration

**Principle**: Make team workflows visible and manageable.

- Clear branch naming conventions
- Easy to see what others are working on
- Safe cleanup of completed work
- Predictable directory structures

---

## Core Workflow Principles

### The Golden Rules

1. **Always work from the main repository**
   - Commands must be run where `.git` is a directory, not a file
   - This ensures you have the complete repository context

2. **One feature, one branch, one worktree**
   - Each feature/bugfix gets its own isolated environment
   - Clear separation of concerns

3. **Dry-run before execute**
   - Always test destructive operations first
   - Review what will happen before committing

4. **Merge, then delete**
   - Completed work should be merged to main
   - Abandoned work can be deleted entirely

5. **Clean as you go**
   - Remove old worktrees regularly
   - Keep only active branches
   - Maintain a clean workspace

### The Worktree Lifecycle

The worktree lifecycle follows a clear flow from idea to completion or abandonment:

1. **Idea/Requirement**: A new feature, bug fix, or task is identified that requires isolated development

2. **Setup**: Create the development environment
   - For new projects: `worktree clone <url>` creates the repository structure
   - For existing projects: `worktree <name>` creates the worktree with auto-detection

3. **Development in Worktree**: Active development happens in the isolated worktree environment
   - Make code changes
   - Commit changes to the branch
   - Push changes to remote (for collaboration)

4. **Completion Path** (for successful work):
   - When work is complete and tested: `worktree --merge --confirm <name>`
   - This merges changes into main/master, removes the worktree directory, cleans up empty parent directories, deletes local branch, and deletes remote branch
   - Result: Changes are merged to main

5. **Abandonment Path** (for discarded work):
   - When work is abandoned or no longer needed: `worktree --delete --confirm <name>`
   - This removes the worktree directory, cleans up empty parent directories, deletes local branch, and deletes remote branch
   - Result: Branch and worktree are completely removed, no merge occurs

---

## Getting Started Workflow

### New Project Setup

**Scenario**: You're starting work on a new project.

#### Step 1: Clone the Repository

```bash
# Navigate to your projects directory
cd ~/Projects

# Clone the repository with structured organization
worktree clone https://github.com/company/new-project.git

# This creates:
# ~/Projects/new-project/main/
```

#### Step 2: Explore the Project

```bash
# Navigate to the main repository
cd new-project/main

# Check the current state
git status
worktree list

# View all branches (local and remote)
worktree list
```

#### Step 3: Set Up Development Environment

```bash
# Install dependencies (project-specific)
npm install  # if Node.js
pip install -r requirements.txt  # if Python
# etc.

# Configure IDE/editors
code .
```

**Workflow Summary**:
```
1. Clone repository → structured directory
2. Navigate to main → explore project
3. Set up environment → ready to develop
```

### New Project Initialization

**Scenario**: You're starting a new project from scratch (not cloning an existing repository).

#### Step 1: Initialize Project with Worktree Structure

```bash
# Navigate to your projects directory
cd ~/Projects

# Initialize a new project with worktree structure
worktree init my-new-project

# This creates:
# ~/Projects/my-new-project/master/
#   └── .git/ (newly initialized git repository)
```

#### Step 2: Set Up the Project

```bash
# Navigate to the main repository
cd my-new-project/master

# Initialize project files (create README, etc.)
echo "# My New Project" > README.md

# Add and commit initial files
git add README.md
git commit -m "Initial commit"
```

#### Step 3: Set Up Development Environment

```bash
# Install dependencies (project-specific)
npm init -y  # if Node.js
# etc.

# Configure IDE/editors
code .
```

**Workflow Summary**:
```
1. Initialize project → creates project/[branch]/ with .git
2. Navigate to main → set up initial files
3. Commit initial state → ready to develop
```

**Note**: Use `worktree init` for new projects, `worktree clone` for existing repositories.

---

## Feature Development Workflow

### New Feature Development

**Scenario**: You need to implement a new user authentication feature.

#### Step 1: Create Feature Branch and Worktree

```bash
# From the main repository
cd ~/Projects/new-project/main

# Create new feature branch with worktree
worktree feat/user-authentication

# This:
# 1. Creates branch 'feat/user-authentication' from current HEAD
# 2. Creates worktree at '../feat/user-authentication'
# 3. Starts bash shell in the new worktree
```

#### Step 2: Develop the Feature

```bash
# You're now in the worktree: ~/Projects/new-project/feat/user-authentication

# Make changes to files
echo "Authentication logic" > auth.js

# Commit changes
git add auth.js
git commit -m "Add user authentication logic"

# Continue development...
```

#### Step 3: Push to Remote

```bash
# Push the new branch to remote (remember the message from step 1)
git push -u origin feat/user-authentication

# Continue development, pushing as needed
```

#### Step 4: Merge When Complete

```bash
# Exit worktree (Ctrl+D or exit)
# Back in main repository: ~/Projects/new-project/main

# First, validate the merge is possible
worktree --merge feat/user-authentication

# Review the dry-run output:
# - Branch status relative to main
# - Merge possibility
# - Worktree removal validation
# - Branch deletion validation

# If all checks pass, execute the merge
worktree --merge --confirm feat/user-authentication

# This:
# 1. Merges feat/user-authentication into main
# 2. Removes the worktree directory
# 3. Cleans up empty parent directories (../feat/ if empty)
# 4. Deletes the local branch
# 5. Deletes the remote branch
```

**Workflow Summary**:
```
1. Create feature branch → isolated worktree
2. Develop feature → commit changes
3. Push to remote → share with team
4. Merge when complete → clean up automatically
```

### Feature with Team Collaboration

**Scenario**: Multiple team members are working on the same feature.

#### Step 1: Create from Remote Branch

```bash
# Someone else created the branch and pushed to remote
# You can create a worktree from the remote branch (auto-detected)
worktree feat/team-feature

# This:
# 1. Creates local tracking branch from origin/feat/team-feature (with upstream tracking)
# 2. Creates worktree at '../feat/team-feature'
# 3. Starts bash shell in the worktree
```

#### Step 2: Collaborate

```bash
# Work in the worktree
# Make changes, commit, push as usual

# Pull changes from teammates
git pull origin feat/team-feature
```

#### Step 3: Final Merge

```bash
# Back in main repository
worktree --merge --confirm feat/team-feature
```

**Key Difference**: Auto-detection will find the remote branch and create a local tracking branch with upstream.

---

## Bug Fix Workflow

### Urgent Bug Fix

**Scenario**: A critical bug needs to be fixed immediately.

#### Step 1: Create Bug Fix Branch

```bash
# From main repository
worktree fix/login-bug-123

# Auto-detection will find existing remote branch if it exists
```

#### Step 2: Fix the Bug

```bash
# In the worktree
# Reproduce the issue
# Implement the fix
# Test the fix

git add .
git commit -m "Fix login bug #123"
```

#### Step 3: Hotfix Merge

```bash
# Back in main
# Validate and merge quickly
worktree --merge --confirm fix/login-bug-123

# Or if you need to deploy immediately
git checkout main
git merge fix/login-bug-123
git push origin main

# Then clean up
worktree -d fix/login-bug-123  # Just remove worktree
# Or
worktree --delete --confirm fix/login-bug-123  # Remove everything
```

**Workflow Summary**:
```
1. Create bug fix branch → isolated environment
2. Fix the issue → test thoroughly
3. Merge quickly → deploy if needed
4. Clean up → maintain organization
```

### Bug Fix with Main Branch Behind

**Scenario**: Main branch has moved forward since you started the bug fix.

#### Step 1: Rebase Before Merging

```bash
# In the worktree
git checkout fix/my-bug
git rebase main

# Resolve any conflicts
# Test the fix again

git push -f origin fix/my-bug  # Force push after rebase
```

#### Step 2: Merge

```bash
# Back in main
worktree --merge --confirm fix/my-bug
```

**Note**: If the branch is behind main, the dry-run will show this and the merge will fail. Rebase first.

---

## Documentation Workflow

### New Documentation

**Scenario**: Adding new documentation to the project.

#### Step 1: Create Documentation Branch

```bash
# From main repository
worktree docs/api-guide

# Or for README updates
worktree docs/readme-update
```

#### Step 2: Write Documentation

```bash
# In the worktree
# Write documentation
# Commit changes
git add docs/api.md
git commit -m "Add API documentation"
git push -u origin docs/api-guide
```

#### Step 3: Merge Documentation

```bash
# Back in main
worktree --merge --confirm docs/api-guide
```

**Best Practice**: Use `docs/` prefix for documentation branches to keep them organized.

---

## Team Collaboration Workflow

### Daily Development Flow

#### Morning: Check Current Work

```bash
# List all branches
worktree list

# List all branches (local and remote) to see what teammates are working on
worktree list

# See active worktrees
git worktree list
```

#### During Development: Pull and Update

```bash
# In your worktree
# Regularly pull from remote
git pull origin <your-branch>

# If main has moved forward, rebase
git fetch origin
git rebase origin/main
```

#### End of Day: Push Progress

```bash
# In your worktree
git add .
git commit -m "End of day progress"
git push origin <your-branch>
```

### Code Review Workflow

#### Step 1: Create Feature Branch

```bash
worktree feat/amazing-feature
```

#### Step 2: Develop and Push

```bash
# Develop feature
git add .
git commit -m "Implement amazing feature"
git push -u origin feat/amazing-feature
```

#### Step 3: Create Pull Request

```bash
# Create PR through GitHub/GitLab UI
# Point PR to feat/amazing-feature → main
```

#### Step 4: Address Review Comments

```bash
# In the worktree
# Make changes based on feedback
git add .
git commit -m "Address review comments"
git push origin feat/amazing-feature
```

#### Step 5: Merge After Approval

```bash
# Back in main
worktree --merge --confirm feat/amazing-feature
```

**Note**: The merge will fail if the branch is behind main. Make sure to rebase or merge main into your branch first.

### Team Member Onboarding

#### Step 1: Clone Project

```bash
worktree clone https://github.com/company/existing-project.git
```

#### Step 2: Set Up All Current Branches

```bash
cd existing-project/main

# See what branches exist (local and remote)
worktree list

# Create worktrees for branches you need to work on (auto-detected)
worktree feat/current-feature
worktree fix/urgent-bug
```

#### Step 3: Start Contributing

```bash
# Use existing workflow for your contributions
worktree feat/my-first-contribution
```

---

## Project Setup Workflow

### Setting Up Multiple Related Projects

**Scenario**: You need to work on a frontend and backend project together.

#### Step 1: Clone Both Projects

```bash
# Clone frontend
worktree clone https://github.com/company/frontend.git

# Clone backend
worktree clone https://github.com/company/backend.git

# Directory structure:
# projects/
# ├── frontend/
# │   └── main/
# └── backend/
#     └── main/
```

#### Step 2: Create Corresponding Feature Branches

```bash
# Frontend
cd frontend/main
worktree feat/new-api-integration

# Backend (in another terminal)
cd backend/main
worktree feat/new-api-endpoints
```

#### Step 3: Work on Both Simultaneously

```bash
# Terminal 1: Frontend worktree
# Terminal 2: Backend worktree
# Work on both aspects of the feature in parallel
```

#### Step 4: Coordinate Merges

```bash
# Merge backend first (if backend changes are needed for frontend)
cd backend/main
worktree --merge --confirm feat/new-api-endpoints

# Then merge frontend
cd frontend/main
worktree --merge --confirm feat/new-api-integration
```

### Monorepo Setup

**Scenario**: Your project is a monorepo with multiple components.

#### Step 1: Clone Monorepo

```bash
worktree clone https://github.com/company/monorepo.git
```

#### Step 2: Create Component-Specific Branches

```bash
cd monorepo/main

# Feature affecting multiple components
worktree feat/shared-auth

# Bug fix in specific component
worktree fix/component-b-bug
```

#### Step 3: Work on Specific Components

```bash
# In the worktree
cd packages/component-a
# Make changes to component A

cd packages/component-b
# Make changes to component B

# Commit all changes together
git add .
git commit -m "Update multiple components"
```

---

## Cleanup and Maintenance Workflow

### Regular Cleanup

#### Step 1: Identify Completed Work

```bash
# List all branches
worktree list

# Check which branches have been merged to main
git branch --merged main
```

#### Step 2: Clean Up Completed Features

```bash
# For each completed feature
worktree --merge --confirm feat/completed-feature

# This removes worktree, local branch, and remote branch
```

#### Step 3: Clean Up Abandoned Work

```bash
# For abandoned features
worktree --delete --confirm feat/abandoned-feature

# This removes worktree, local branch, and remote branch
# No merge - changes are discarded
```

#### Step 4: Clean Up Old Worktrees

```bash
# Remove worktrees for branches that were deleted other ways
worktree -d old-branch-name

# This only removes the worktree directory
# Use when branch was deleted manually or through other means
```

### Monthly Maintenance

#### Step 1: Prune Remote References

```bash
# Clean up remote tracking branches
git remote prune origin
```

#### Step 2: Optimize Repository

```bash
# In main repository
git gc
git repack
```

#### Step 3: Organize Directory Structure

```bash
# Check for empty directories
find .. -type d -empty -delete

# Verify all worktrees are registered
git worktree list
```

### Emergency Cleanup

**Scenario**: Something went wrong and you need to reset.

#### Step 1: Check Current State

```bash
# List all worktrees
git worktree list

# Check branch status
git status
worktree list
```

#### Step 2: Safe Cleanup

```bash
# For each problematic worktree
worktree --delete --confirm branch-name

# Or just remove worktree directories
worktree -d branch-name
```

#### Step 3: Reset Main Repository

```bash
# In main repository
git reset --hard origin/main
git clean -fd
```

---

## Advanced Workflows

### Cross-Branch Testing

**Scenario**: You need to test changes from one branch in the context of another.

#### Step 1: Create Both Branches

```bash
worktree feat/feature-a
worktree feat/feature-b
```

#### Step 2: Test Integration

```bash
# In feature-a worktree
git checkout feat/feature-b -- paths/to/test

# Or cherry-pick specific commits
git cherry-pick <commit-hash>
```

### Temporary Experimental Branches

**Scenario**: Quick experimentation without polluting the main workflow.

#### Step 1: Create Experimental Branch

```bash
worktree exp/test-idea
```

#### Step 2: Experiment

```bash
# Try things out
# Make changes
# Test concepts
```

#### Step 3: Discard Experiment

```bash
# Back in main
worktree --delete --confirm exp/test-idea

# All traces removed
```

### Branch Rename Workflow

**Scenario**: You need to rename a branch.

#### Step 1: Rename Local Branch

```bash
# In main repository
git branch -m old-name new-name
```

#### Step 2: Update Worktree

```bash
# The worktree is still at ../old-name
# Need to manually move it
mv ../old-name ../new-name

# Update git worktree registration
git worktree remove ../old-name
git worktree add ../new-name new-name
```

#### Step 3: Rename Remote Branch

```bash
# Push new name
git push origin new-name

# Delete old name
git push origin --delete old-name
```

**Note**: Branch renaming is not directly supported by worktree tool. Use git commands directly.

### Submodule Workflow

**Scenario**: Your project uses git submodules.

#### Step 1: Clone with Submodules

```bash
worktree clone https://github.com/company/project-with-submodules.git
cd project-with-submodules/main

# Initialize submodules
git submodule update --init --recursive
```

#### Step 2: Create Worktree

```bash
worktree feat/new-feature

# Submodules work normally in worktrees
```

---

## Workflow Visualization

### Directory Structure Evolution

```
Initial State:
└── projects/

After cloning two projects:
└── projects/
    ├── project-a/
    │   └── main/
    │       └── .git/ (directory)
    └── project-b/
        └── master/
            └── .git/ (directory)

After creating feature branches:
└── projects/
    ├── project-a/
    │   ├── main/
    │   │   └── .git/
    │   └── feat/
    │       ├── login/
    │       │   └── .git (file - points to main repo)
    │       └── payment/
    │           └── .git (file)
    └── project-b/
        ├── master/
        │   └── .git/
        └── fix/
            └── security/
                └── .git (file)

After merging login feature:
└── projects/
    ├── project-a/
    │   ├── main/
    │   │   └── .git/
    │   └── feat/
    │       └── payment/
    │           └── .git (file)
    └── project-b/
        ├── master/
        │   └── .git/
        └── fix/
            └── security/
                └── .git (file)
```

### Command Flow Description

For the clone command `worktree clone <url>`, the processing flow is:

1. User provides repository URL as input

2. **Create project directory**: The tool creates a directory named after the project extracted from the URL

3. **Clone repo into temp directory**: The repository is cloned into a temporary directory within the project directory

4. Two parallel operations occur:
   - **Determine project name**: The `extractProjectName(url)` function parses the URL to extract the repository name
   - **Get default branch name**: The tool determines the repository's default branch by first checking origin/HEAD, then falling back to common names (main, master, trunk)

5. **Rename directory**: The cloned repository directory is renamed from its temporary name to the detected branch name

6. **Success output**: The tool prints "SUCCESS: Cloned into project/branch/" confirming the operation completed

---

## Best Practices by Scenario

### Solo Developer

| Scenario | Recommended Workflow |
|----------|---------------------|
| New feature | `worktree feat/name` → develop → `--merge --confirm` |
| Bug fix | `worktree fix/name` → fix → `--merge --confirm` |
| Experiment | `worktree exp/name` → test → `--delete --confirm` |
| Documentation | `worktree docs/name` → write → `--merge --confirm` |

### Team Lead

| Scenario | Recommended Workflow |
|----------|---------------------|
| Assign work | Team member: `worktree feat/assigned` |
| Review PR | Check out PR branch locally |
| Coordinate merge | `worktree --merge --confirm` after approval |
| Clean up old branches | Regular `worktree --delete --confirm` |

### Open Source Maintainer

| Scenario | Recommended Workflow |
|----------|---------------------|
| Review contribution | `worktree pr-branch` → test → provide feedback |
| Merge contribution | `worktree --merge --confirm` after CI passes |
| Triage issues | Create branches for issue investigation |

### DevOps Engineer

| Scenario | Recommended Workflow |
|----------|---------------------|
| Hotfix | `worktree fix/hotfix` → test → `--merge --confirm` → deploy |
| Infrastructure | `worktree infra/terraform` → update → `--merge --confirm` |
| CI/CD updates | `worktree ci/new-pipeline` → update → `--merge --confirm` |

---

## Summary

The worktree tool is designed to support efficient, organized, and safe Git workflows. By following the patterns and principles outlined in this document, you can:

✅ **Maintain organized codebases** with clear directory structures  
✅ **Work efficiently** with isolated development environments  
✅ **Collaborate effectively** with team members  
✅ **Ensure safety** with dry-run validation and explicit confirmation  
✅ **Scale your workflow** from solo development to team collaboration  

**Remember the core workflow**:
1. **Create** worktrees for focused development
2. **Develop** in isolated environments
3. **Validate** with dry-run operations
4. **Execute** with explicit confirmation
5. **Clean up** to maintain organization

For more details, see:
- [README.md](./README.md) - Main documentation with command reference
- [USAGE.md](./USAGE.md) - Comprehensive usage guide with examples
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Technical implementation details
- [AGENTS.md](./AGENTS.md) - AI agent guidance