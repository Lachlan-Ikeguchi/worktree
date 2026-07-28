package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// IsWorktree checks if the current directory is a Git worktree.
// A worktree has .git as a file (symlink to the main repo), not a directory.
func IsWorktree() bool {
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

// IsGitRepo checks if the current directory is a Git repository.
// A Git repository has .git as a directory.
func IsGitRepo() bool {
	info, err := os.Stat(".git")
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetMainBranch returns the main branch name from the repository.
// It tries origin/HEAD first, then falls back to common branch names.
func GetMainBranch() (string, error) {
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

// BranchExists checks if a local branch exists.
func BranchExists(branch string) bool {
	cmd := exec.Command("git", "show-ref", "--quiet", "refs/heads/"+branch)
	return cmd.Run() == nil
}

// GetLocalBranches returns a list of all local branch names.
func GetLocalBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	// Filter out empty strings
	var result []string
	for _, b := range branches {
		if b != "" {
			result = append(result, b)
		}
	}
	return result, nil
}

// GetRemoteBranches returns a list of all remote branch names (without origin/ prefix).
func GetRemoteBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "-r", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	// Filter and remove origin/ prefix
	seen := make(map[string]bool)
	var result []string
	for _, b := range branches {
		if b != "" {
			// Remove origin/ prefix
			branch := strings.TrimPrefix(b, "origin/")
			// Skip invalid branch names and deduplicate
			if branch != "" && branch != "origin" && !seen[branch] {
				seen[branch] = true
				result = append(result, branch)
			}
		}
	}
	return result, nil
}

// GetCurrentBranch returns the current branch name for a given repository path.
func GetCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetDefaultBranchName returns the default branch name from git config
// or defaults to "master" if not configured.
func GetDefaultBranchName() string {
	// Try to get from git config
	cmd := exec.Command("git", "config", "--get", "init.defaultBranch")
	output, err := cmd.Output()
	if err == nil {
		branch := strings.TrimSpace(string(output))
		if branch != "" {
			return branch
		}
	}

	// Default to master if not configured
	return "master"
}

// GetAllBranches returns both local and remote branches combined and deduplicated.
func GetAllBranches() ([]string, error) {
	seen := make(map[string]bool)
	var allBranches []string

	// Get local branches
	localBranches, err := GetLocalBranches()
	if err != nil {
		return nil, err
	}
	for _, b := range localBranches {
		if !seen[b] {
			seen[b] = true
			allBranches = append(allBranches, b)
		}
	}

	// Get remote branches
	remoteBranches, err := GetRemoteBranches()
	if err != nil {
		return nil, err
	}
	for _, b := range remoteBranches {
		if !seen[b] {
			seen[b] = true
			allBranches = append(allBranches, b)
		}
	}

	return allBranches, nil
}

// GetBranchStatus returns the number of commits a branch is ahead and behind relative to another branch.
func GetBranchStatus(branch, mainBranch string) (ahead int, behind int, err error) {
	// Check commits ahead (branch has that main doesn't)
	cmd := exec.Command("git", "log", "--oneline", fmt.Sprintf("%s..%s", mainBranch, branch))
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("could not check commits ahead: %v", err)
	}
	ahead = len(strings.Split(strings.TrimSpace(string(output)), "\n"))

	// Check commits behind (main has that branch doesn't)
	cmd = exec.Command("git", "log", "--oneline", fmt.Sprintf("%s..%s", branch, mainBranch))
	output, err = cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("could not check commits behind: %v", err)
	}
	behind = len(strings.Split(strings.TrimSpace(string(output)), "\n"))

	return ahead, behind, nil
}

// IsBranchMerged checks if a branch is fully merged into the main branch.
func IsBranchMerged(branch, mainBranch string) (bool, error) {
	cmd := exec.Command("git", "branch", "--merged", mainBranch)
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("could not check merged branches: %v", err)
	}
	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, b := range branches {
		// Remove leading * and spaces
		b = strings.TrimLeft(b, " *")
		if b == branch {
			return true, nil
		}
	}
	return false, nil
}
