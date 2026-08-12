package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lachlan/worktree/internal/git"
)

// CreateWorktree creates a new worktree for the given branch.
// It auto-detects the branch source:
// 1. If remote branch exists (origin/<branch>), creates local tracking branch from remote
// 2. Else if local branch exists, uses the existing local branch
// 3. Otherwise, creates a new branch from HEAD
func CreateWorktree(branch string) error {
	worktreePath := filepath.Join("..", branch)

	// If worktree already exists, just start shell
	if _, err := os.Stat(worktreePath); err == nil {
		// Validate that the branch exists (local or remote)
		if !git.BranchExists(branch) && !git.RemoteBranchExists(branch) {
			return fmt.Errorf("branch '%s' does not exist locally or remotely", branch)
		}
	} else {
		// Auto-detect: check remote first, then local, then create new
		if git.RemoteBranchExists(branch) {
			// Create local tracking branch from remote
			cmd := exec.Command("git", "fetch", "origin", fmt.Sprintf("%s:%s", branch, branch))
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to create local tracking branch: %v", err)
			}
			// Set upstream tracking
			cmd = exec.Command("git", "branch", "--set-upstream-to=origin/"+branch, branch)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to set upstream tracking: %v", err)
			}
		} else if git.BranchExists(branch) {
			// Use existing local branch - ensure it has upstream tracking if possible
			if git.RemoteBranchExists(branch) {
				cmd := exec.Command("git", "branch", "--set-upstream-to=origin/"+branch, branch)
				cmd.Run() // Best effort - don't fail if this doesn't work
			}
		} else {
			// Create new branch from HEAD
			cmd := exec.Command("git", "branch", branch)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to create branch: %v", err)
			}
			fmt.Println()
			fmt.Println("Remember to push to create tracking branches in remote")
			fmt.Println()
		}

		// Create worktree
		cmd := exec.Command("git", "worktree", "add", worktreePath, branch)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create worktree: %v", err)
		}
	}

	// Start a new bash shell in the worktree
	cmd := exec.Command("bash", "-c", fmt.Sprintf("cd %q; exec bash", worktreePath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start shell: %v", err)
	}

	return nil
}
