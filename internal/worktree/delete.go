package worktree

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// DeleteWorktree removes a worktree directory and cleans up empty parent directories.
func DeleteWorktree(branch string) error {
	worktreePath := filepath.Join("..", branch)

	cmd := exec.Command("git", "worktree", "remove", worktreePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not remove worktree: %v", err)
	}

	// Clean up empty parent directories
	CleanupEmptyDirs(filepath.Dir(worktreePath))

	fmt.Printf("Deleted worktree at %s\n", worktreePath)

	return nil
}

// DeleteWorktreeOnly removes just the worktree directory without cleaning up branches.
// This is used when the -d flag is used.
func DeleteWorktreeOnly(branch string) error {
	worktreePath := filepath.Join("..", branch)

	cmd := exec.Command("git", "worktree", "remove", worktreePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not remove worktree: %v", err)
	}

	// Clean up empty parent directories
	CleanupEmptyDirs(filepath.Dir(worktreePath))

	fmt.Printf("Deleted worktree at %s\n", worktreePath)

	return nil
}
