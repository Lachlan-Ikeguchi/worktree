package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lachlan/worktree/internal/color"
	"github.com/lachlan/worktree/internal/git"
)

// MergeOrDelete handles both merge and delete operations with dry-run support.
// If confirm is false, it performs a dry run to check if operations are possible.
// If confirm is true, it executes the operations.
func MergeOrDelete(branch string, mergeMode, deleteMode, confirm bool) error {
	mainBranch, err := git.GetMainBranch()
	if err != nil {
		return fmt.Errorf("could not determine main branch: %v", err)
	}

	if !git.BranchExists(branch) {
		return fmt.Errorf("branch '%s' does not exist", branch)
	}

	worktreePath := filepath.Join("..", branch)

	if confirm {
		return executeMergeOrDelete(branch, mainBranch, worktreePath, mergeMode, deleteMode)
	}

	// Dry run - test if operations are actually possible
	return dryRunMergeOrDelete(branch, mainBranch, worktreePath, mergeMode, deleteMode)
}

// executeMergeOrDelete performs the actual merge and/or delete operations.
func executeMergeOrDelete(branch, mainBranch, worktreePath string, mergeMode, deleteMode bool) error {
	if mergeMode {
		// Perform merge
		cmd := exec.Command("git", "merge", branch)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("merge failed: %v", err)
		}
	}

	// Remove worktree
	cmd := exec.Command("git", "worktree", "remove", worktreePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not remove worktree: %v", err)
	}

	// Clean up empty parent directories
	CleanupEmptyDirs(filepath.Dir(worktreePath))

	// Delete local branch
	cmd = exec.Command("git", "branch", "-d", branch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete local branch: %v", err)
	}

	// Delete remote branch
	cmd = exec.Command("git", "push", "-d", "origin", branch)
	if err := cmd.Run(); err != nil {
		fmt.Println("Warning: remote branch may not exist or deletion failed.")
	}

	return nil
}

// dryRunMergeOrDelete tests if merge and/or delete operations are possible.
func dryRunMergeOrDelete(branch, mainBranch, worktreePath string, mergeMode, deleteMode bool) error {
	typeStr := "Delete"
	flagStr := "--delete"
	if mergeMode {
		typeStr = "Merge"
		flagStr = "--merge"
	}

	fmt.Println("=== DRY RUN - Testing if operations are possible ===")
	fmt.Println()

	allPossible := true

	// Check branch status relative to main (for both merge and delete modes)
	fmt.Printf("Checking: Branch %s status relative to %s\n", branch, mainBranch)

	// Check commits ahead
	ahead, behind, err := git.GetBranchStatus(branch, mainBranch)
	if err != nil {
		fmt.Printf("  %s: Could not check branch status: %v\n", color.Yellow("WARNING"), err)
	} else {
		if ahead > 0 {
			fmt.Printf("  %s: Branch %s has %d commit(s) ahead of %s\n", color.Blue("INFO"), branch, ahead, mainBranch)
		}
		if behind > 0 {
			fmt.Printf("  %s: Branch %s is %d commit(s) behind %s\n", color.Red("FAIL"), branch, behind, mainBranch)
			if mergeMode {
				allPossible = false
			}
		}
		if ahead == 0 && behind == 0 {
			fmt.Printf("  %s: Branch %s has no new commits (equal to %s)\n", color.Blue("INFO"), branch, mainBranch)
		}
	}

	if mergeMode {
		// Test merge
		fmt.Printf("Testing: %s %s onto %s\n", typeStr, branch, mainBranch)
		cmd := exec.Command("git", "merge", "--no-commit", "--no-ff", branch)
		if err := cmd.Run(); err != nil {
			fmt.Printf("  %s: Merge would fail: %v\n", color.Red("FAIL"), err)
			allPossible = false
		} else {
			// Reset the merge (we only tested with --no-commit)
			exec.Command("git", "merge", "--abort").Run() // Ignore error, might not have started
			fmt.Printf("  %s: Merge is possible\n", color.Green("PASS"))
		}
	}

	// Test worktree removal
	fmt.Printf("Testing: Delete worktree at %s\n", worktreePath)
	// Get absolute path for comparison
	absWorktreePath, err := filepath.Abs(worktreePath)
	if err != nil {
		absWorktreePath = worktreePath
	}

	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		fmt.Printf("  %s: Worktree directory not found at %s\n", color.Red("FAIL"), worktreePath)
		allPossible = false
	} else {
		// Check if worktree is valid and can be removed
		cmd := exec.Command("git", "worktree", "list")
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("  %s: Could not list worktrees: %v\n", color.Red("FAIL"), err)
			allPossible = false
		} else {
			worktreeList := string(output)
			// Check both relative and absolute paths
			if !strings.Contains(worktreeList, worktreePath) && !strings.Contains(worktreeList, absWorktreePath) {
				fmt.Printf("  %s: Worktree at %s is not registered\n", color.Red("FAIL"), worktreePath)
				allPossible = false
			} else {
				fmt.Printf("  %s: Worktree can be removed\n", color.Green("PASS"))
			}
		}
	}

	// Test local branch deletion
	fmt.Printf("Testing: Delete local branch %s\n", branch)
	if !git.BranchExists(branch) {
		fmt.Printf("  %s: Local branch '%s' does not exist\n", color.Red("FAIL"), branch)
		allPossible = false
	} else {
		// Check if branch is fully merged to main branch
		merged, err := git.IsBranchMerged(branch, mainBranch)
		if err != nil {
			fmt.Printf("  %s: Could not check merged branches: %v\n", color.Red("FAIL"), err)
			allPossible = false
		} else {
			if merged {
				fmt.Printf("  %s: Local branch can be deleted (fully merged)\n", color.Green("PASS"))
			} else {
				// Branch not merged, but can be force deleted
				fmt.Printf("  %s: Local branch can be deleted (with force -D)\n", color.Green("PASS"))
			}
		}
	}

	// Test remote branch deletion
	fmt.Printf("Testing: Delete remote branch origin/%s\n", branch)
	// First check if remote branch exists
	cmd := exec.Command("git", "show-ref", "--quiet", "refs/remotes/origin/"+branch)
	if err := cmd.Run(); err != nil {
		// Remote branch doesn't exist - that's okay, just warn
		fmt.Printf("  %s: Remote branch origin/%s does not exist (will be skipped)\n", color.Yellow("WARNING"), branch)
	} else {
		// Remote branch exists, test deletion
		cmd = exec.Command("git", "push", "-d", "--dry-run", "origin", branch)
		if err := cmd.Run(); err != nil {
			fmt.Printf("  %s: Remote branch deletion would fail: %v\n", color.Red("FAIL"), err)
			allPossible = false
		} else {
			fmt.Printf("  %s: Remote branch can be deleted\n", color.Green("PASS"))
		}
	}

	fmt.Println()
	if allPossible {
		fmt.Println(color.Green("All operations are possible!"))
	} else {
		fmt.Println(color.Red("Some operations would fail - see above for details"))
	}
	fmt.Println()
	fmt.Printf("To execute, run: worktree %s --confirm %s\n", flagStr, branch)
	fmt.Println("Warning: ensures no one else is working on this branch - it will be deleted")

	return nil
}
