package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: provide a repository\n")
		os.Exit(1)
	}

	repoURL := os.Args[1]

	// Extract project name from URL
	projectName := extractProjectName(repoURL)

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to get current directory: %v\n", err)
		os.Exit(1)
	}

	// Create project directory
	projectPath := filepath.Join(currentDir, projectName)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create directory: %v\n", err)
		os.Exit(1)
	}

	// Clone repository into the project directory
	cloneCmd := exec.Command("git", "clone", repoURL, projectName)
	cloneCmd.Dir = projectPath
	if err := cloneCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to clone repository: %v\n", err)
		// Clean up the directory we created
		os.RemoveAll(projectPath)
		os.Exit(1)
	}

	// Get the branch name from the cloned repo
	clonedRepoPath := filepath.Join(projectPath, projectName)
	branchName, err := getCurrentBranch(clonedRepoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to get branch name: %v\n", err)
		os.Exit(1)
	}

	// Rename the cloned repo directory to the branch name
	branchPath := filepath.Join(projectPath, branchName)
	if err := os.Rename(clonedRepoPath, branchPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to rename directory: %v\n", err)
		os.Exit(1)
	}
}

func extractProjectName(url string) string {
	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")

	// Handle SSH URLs like git@github.com:owner/repo
	// Split by : or / and take the last element, similar to IFS=':/' read -r -a url
	if strings.Contains(url, "@") && strings.Contains(url, ":") {
		parts := strings.Split(url, ":")
		if len(parts) >= 2 {
			// Get the last part after the colon
			lastPart := parts[len(parts)-1]
			// Further split by / to handle git@host:owner/repo
			pathParts := strings.Split(lastPart, "/")
			if len(pathParts) > 0 {
				return pathParts[len(pathParts)-1]
			}
			return lastPart
		}
	}

	// Handle HTTPS URLs like https://github.com/owner/repo
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return url
}

func getCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
