package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lachlan/worktree/internal/git"
)

// CreateWorktree creates a new worktree for the given branch.
// If remote is true, it creates a local tracking branch from origin/<branch>.
// If existing is true, it uses an existing local branch.
// Otherwise, it creates a new branch from HEAD.
func CreateWorktree(branch string, remote, existing bool) error {
	worktreePath := filepath.Join("..", branch)

	// If worktree already exists, skip creation
	if _, err := os.Stat(worktreePath); err == nil {
		// Branch validation still needed for -e flag
		if existing {
			if !git.BranchExists(branch) {
				return fmt.Errorf("branch '%s' does not exist locally", branch)
			}
		}
	} else {
		// Create the branch
		if remote {
			cmd := exec.Command("git", "fetch", "origin", fmt.Sprintf("%s:%s", branch, branch))
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to create local tracking branch: %v", err)
			}
		} else if !existing {
			cmd := exec.Command("git", "branch", branch)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to create branch: %v", err)
			}
			fmt.Println()
			fmt.Println("Remember to push to create tracking branches in remote")
			fmt.Println()
		} else {
			if !git.BranchExists(branch) {
				return fmt.Errorf("branch '%s' does not exist locally", branch)
			}
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
