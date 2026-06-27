package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Parse flags
	mergeMode := flag.Bool("merge", false, "Merge the branch into main and clean up")
	deleteMode := flag.Bool("delete", false, "Delete the branch, remote branch, and worktree")
	confirm := flag.Bool("confirm", false, "Confirm the merge or delete operation")
	remote := flag.Bool("r", false, "Create local tracking branch from origin/<branch>")
	existing := flag.Bool("e", false, "Create a worktree from an existing local branch")
	delete := flag.Bool("d", false, "Delete the worktree directory")
	help := flag.Bool("h", false, "Show help message")
	
	flag.BoolVar(help, "help", false, "Show help message")
	
	flag.Usage = func() {
		printUsage()
	}

	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) == 0 {
		printUsage()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Error: provide branch name")
		os.Exit(1)
	}

	// Handle list command
	if args[0] == "list" {
		// Check for -r flag in the remaining arguments
		var useRemote bool
		if len(args) > 1 {
			for i := 1; i < len(args); i++ {
				if args[i] == "-r" {
					useRemote = true
				} else {
					// Unknown flag for list command
					fmt.Fprintf(os.Stderr, "Error: unknown flag '%s' for list command\n", args[i])
					printUsage()
					os.Exit(1)
				}
			}
		}

		var gitArgs []string
		if useRemote {
			gitArgs = []string{"branch", "-r"}
		} else {
			gitArgs = []string{"branch"}
		}
		cmd := exec.Command("git", gitArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running git branch: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	branch := args[0]

	// Check if we're in the main repo (has .git/) not a worktree (has .git file)
	if isWorktree() {
		fmt.Fprintln(os.Stderr, "Error: must be called from the main repository, not a worktree")
		os.Exit(1)
	}

	if !isGitRepo() {
		fmt.Fprintln(os.Stderr, "Error: not a git repository")
		os.Exit(1)
	}

	if *mergeMode && *deleteMode {
		fmt.Fprintln(os.Stderr, "Error: cannot use --merge and --delete together")
		os.Exit(1)
	}

	if *mergeMode || *deleteMode {
		if *deleteMode && !*confirm {
			fmt.Fprintln(os.Stderr, "Error: --delete requires --confirm")
			os.Exit(1)
		}

		if *remote || *existing || *delete {
			fmt.Fprintln(os.Stderr, "Error: cannot use -r, -e, or -d with --merge or --delete")
			os.Exit(1)
		}

		mergeOrDelete(branch, *mergeMode, *deleteMode, *confirm)
		os.Exit(0)
	}

	if *delete {
		if *remote || *existing {
			fmt.Fprintln(os.Stderr, "Error: cannot use -r or -e with -d")
			os.Exit(1)
		}

		if *confirm {
			fmt.Fprintln(os.Stderr, "Error: --confirm is only used with --merge or --continue")
			os.Exit(1)
		}

		deleteWorktree(branch)
		os.Exit(0)
	}

	// CREATE MODE
	if *remote && *existing {
		fmt.Fprintln(os.Stderr, "Error: cannot use switches -r and -e at the same time")
		os.Exit(1)
	}

	if *delete {
		fmt.Fprintln(os.Stderr, "Error: -d cannot be used with create mode")
		os.Exit(1)
	}

	if *confirm {
		fmt.Fprintln(os.Stderr, "Error: --confirm is only used with --merge or --continue")
		os.Exit(1)
	}

	createWorktree(branch, *remote, *existing)
}

func printUsage() {
	fmt.Println(`Usage: worktree [OPTIONS] <branch>
Usage: worktree list

Create a branch and a worktree linked to it, delete a worktree, or merge and clean up.

Options:
  -r             Create local tracking branch from origin/<branch> (requires fetch).
                 Mutually exclusive with -e and -d
  -e             Create a worktree from an existing local branch.
                 Mutually exclusive with -r and -d
  -d             Delete the worktree directory at ../<branch> and clean up empty parent directories.
                 Mutually exclusive with -r, -e, --merge, --delete
  --merge        Merges the branch into master and clean up worktree (requires --confirm)
  --delete       Delete the branch, remote branch, and worktree (requires --confirm)
  --confirm      Confirm the merge or delete operation (used with --merge, --delete)
  -h             Show this help message

Commands:
  list           List all worktrees (runs 'git branch')
  list -r        List all remote branches (runs 'git branch -r')

Without switches: creates a new local branch from current HEAD.

Create mode (default):
  1. Creates a local branch (-e skips this step)
  2. Creates a worktree at ../<branch>
  3. Starts a new bash shell in the new worktree
    - If worktree already exists, enters it in a new shell

Delete mode (-d):
  1. Deletes the worktree at ../<branch>
  2. Cleans up empty parent directories

Merge mode (--merge):
  1. Merges <branch> onto master/main/trunk
  2. Deletes the worktree directory at ../<branch>
  3. Cleans up empty parent directories
  4. Deletes the local branch
  5. Deletes the remote branch (origin/<branch>)
  Note: dry-run by default, requires --confirm to actually perform actions

Delete mode (--delete):
  1. Deletes the worktree directory at ../<branch> and empty parent directories
  2. Deletes the local branch
  3. Deletes the remote branch (origin/<branch>)
  Note: Requires --confirm.`) 
}

func isWorktree() bool {
	// Check if .git is a file (indicates worktree)
	if _, err := os.Stat(".git"); err == nil {
		// .git exists, check if it's a file
		info, err := os.Stat(".git")
		if err == nil && !info.IsDir() {
			return true
		}
	}
	return false
}

func isGitRepo() bool {
	info, err := os.Stat(".git")
	if err != nil {
		return false
	}
	return info.IsDir()
}

func getMainBranch() (string, error) {
	// Try to get the main branch from origin/HEAD
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.Output()
	if err == nil {
		ref := strings.TrimSpace(string(output))
		// Extract branch name from refs/remotes/origin/<branch>
		parts := strings.Split(ref, "/")
		if len(parts) >= 4 {
			return parts[3], nil
		}
	}

	// Fall back to checking common branch names
	for _, candidate := range []string{"main", "master", "trunk"} {
		cmd := exec.Command("git", "show-ref", "--quiet", "refs/heads/"+candidate)
		if err := cmd.Run(); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("cannot determine main branch (tried origin/HEAD, main, master, trunk)")
}

func branchExists(branch string) bool {
	cmd := exec.Command("git", "show-ref", "--quiet", "refs/heads/"+branch)
	return cmd.Run() == nil
}

func mergeOrDelete(branch string, mergeMode, deleteMode, confirm bool) {
	mainBranch, err := getMainBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !branchExists(branch) {
		fmt.Fprintf(os.Stderr, "Error: branch '%s' does not exist\n", branch)
		os.Exit(1)
	}

	worktreePath := filepath.Join("..", branch)

	if confirm {
		if mergeMode {
			// Perform merge
			cmd := exec.Command("git", "merge", branch)
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: merge failed: %v\n", err)
				os.Exit(1)
			}
		}

		// Remove worktree
		cmd := exec.Command("git", "worktree", "remove", worktreePath)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not remove worktree: %v\n", err)
			os.Exit(1)
		}

		// Clean up empty parent directories
		cleanupPath := filepath.Dir(worktreePath)
		for cleanupPath != "." && cleanupPath != ".." {
			// Only remove if directory is empty (equivalent to rmdir --ignore-fail-on-non-empty)
			if _, err := os.Stat(cleanupPath); os.IsNotExist(err) {
				break
			}
			entries, err := os.ReadDir(cleanupPath)
			if err != nil {
				break
			}
			if len(entries) == 0 {
				if err := os.Remove(cleanupPath); err != nil {
					break
				}
			} else {
				break
			}
			cleanupPath = filepath.Dir(cleanupPath)
		}

		// Delete local branch
		cmd = exec.Command("git", "branch", "-d", branch)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to delete local branch: %v\n", err)
			os.Exit(1)
		}

		// Delete remote branch
		cmd = exec.Command("git", "push", "-d", "origin", branch)
		if err := cmd.Run(); err != nil {
			fmt.Println("Warning: remote branch may not exist or deletion failed.")
		}
	} else {
		// Dry run
		typeStr := "Delete"
		flagStr := "--delete"
		if mergeMode {
			typeStr = "Merge"
			flagStr = "--merge"
		}

		// Validate worktree exists
		if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: worktree directory not found at %s\n", worktreePath)
			os.Exit(1)
		}

		fmt.Println("=== DRY RUN - No changes will be made ===")
		fmt.Println()
		fmt.Println("Would perform the following:")
		if deleteMode {
			fmt.Printf("  1. Delete worktree at %s and clean up empty parent directories\n", worktreePath)
			fmt.Printf("  2. Delete local branch %s\n", branch)
			fmt.Printf("  3. Delete remote branch origin/%s\n", branch)
		} else {
			fmt.Printf("  1. %s %s onto %s\n", typeStr, branch, mainBranch)
			fmt.Printf("  2. Delete worktree at %s\n", worktreePath)
			fmt.Printf("  3. Clean up empty parent directories\n")
			fmt.Printf("  4. Delete local branch %s\n", branch)
			fmt.Printf("  5. Delete remote branch origin/%s\n", branch)
		}
		fmt.Println()
		fmt.Printf("To execute, run: worktree %s --confirm %s\n", flagStr, branch)
		fmt.Println("Warning: ensures no one else is working on this branch - it will be deleted")
	}
}

func deleteWorktree(branch string) {
	worktreePath := filepath.Join("..", branch)

	cmd := exec.Command("git", "worktree", "remove", worktreePath)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not remove worktree: %v\n", err)
		os.Exit(1)
	}

	// Clean up empty parent directories
	cleanupPath := filepath.Dir(worktreePath)
	for cleanupPath != "." && cleanupPath != ".." {
		// Only remove if directory is empty (equivalent to rmdir --ignore-fail-on-non-empty)
		if _, err := os.Stat(cleanupPath); os.IsNotExist(err) {
			break
		}
		entries, err := os.ReadDir(cleanupPath)
		if err != nil {
			break
		}
		if len(entries) == 0 {
			if err := os.Remove(cleanupPath); err != nil {
				break
			}
		} else {
			break
		}
		cleanupPath = filepath.Dir(cleanupPath)
	}

	fmt.Printf("Deleted worktree at %s\n", worktreePath)
}

func createWorktree(branch string, remote, existing bool) {
	worktreePath := filepath.Join("..", branch)

	// If worktree already exists, skip creation
	if _, err := os.Stat(worktreePath); err == nil {
		// Branch validation still needed for -e flag
		if existing {
			if !branchExists(branch) {
				fmt.Fprintf(os.Stderr, "Error: branch '%s' does not exist locally\n", branch)
				os.Exit(1)
			}
		}
	} else {
		// Create the branch
		if remote {
			cmd := exec.Command("git", "fetch", "origin", fmt.Sprintf("%s:%s", branch, branch))
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to create local tracking branch: %v\n", err)
				os.Exit(1)
			}
		} else if !existing {
			cmd := exec.Command("git", "branch", branch)
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to create branch: %v\n", err)
				os.Exit(1)
			}
			fmt.Println()
			fmt.Println("Remember to push to create tracking branches in remote")
			fmt.Println()
		} else {
			if !branchExists(branch) {
				fmt.Fprintf(os.Stderr, "Error: branch '%s' does not exist locally\n", branch)
				os.Exit(1)
			}
		}

		// Create worktree
		cmd := exec.Command("git", "worktree", "add", worktreePath, branch)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create worktree: %v\n", err)
			os.Exit(1)
		}
	}

	// Start a new bash shell in the worktree
	cmd := exec.Command("bash", "-c", fmt.Sprintf("cd %q; exec bash", worktreePath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to start shell: %v\n", err)
		os.Exit(1)
	}
}
