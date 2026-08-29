package clone

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lachlan/worktree/internal/git"
)

// Clone performs the clone operation for a repository URL.
// It creates a project directory structure like: project_name/[branch]/
func Clone(repoURL string) error {
	// Extract project name from URL
	projectName := extractProjectName(repoURL)

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Create project directory
	projectPath := filepath.Join(currentDir, projectName)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Clone repository into the project directory
	cloneCmd := exec.Command("git", "clone", repoURL, projectName)
	cloneCmd.Dir = projectPath
	if err := cloneCmd.Run(); err != nil {
		// Clean up the directory we created
		os.RemoveAll(projectPath)
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	// Get the branch name from the cloned repo
	clonedRepoPath := filepath.Join(projectPath, projectName)
	branchName, err := git.GetCurrentBranch(clonedRepoPath)
	if err != nil {
		return fmt.Errorf("failed to get branch name: %v", err)
	}

	// Rename the cloned repo directory to the branch name
	branchPath := filepath.Join(projectPath, branchName)
	if err := os.Rename(clonedRepoPath, branchPath); err != nil {
		return fmt.Errorf("failed to rename directory: %v", err)
	}

	fmt.Printf("Successfully cloned: %s\n", projectName)

	return nil
}

// extractProjectName extracts the project name from a repository URL.
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
