package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lachlan/worktree/internal/git"
)

// Init creates a new project with the worktree structure.
// It creates: <current-dir>/<project-name>/<main-branch>/.git/
func Init(projectName string) error {
	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Create project directory
	projectPath := filepath.Join(currentDir, projectName)
	if _, err := os.Stat(projectPath); err == nil {
		return fmt.Errorf("project directory '%s' already exists", projectName)
	}

	// Get default branch name
	branchName := git.GetDefaultBranchName()
	fmt.Printf("Using default branch: %s\n", branchName)

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	// Create branch directory
	branchPath := filepath.Join(projectPath, branchName)
	if err := os.MkdirAll(branchPath, 0755); err != nil {
		// Clean up project directory
		os.RemoveAll(projectPath)
		return fmt.Errorf("failed to create branch directory: %v", err)
	}

	// Initialize git repository
	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = branchPath
	if err := gitCmd.Run(); err != nil {
		// Clean up
		os.RemoveAll(projectPath)
		return fmt.Errorf("failed to initialize git repository: %v", err)
	}

	fmt.Printf("Successfully initialized project '%s' with structure:\n", projectName)
	fmt.Printf("  %s/\n", projectName)
	fmt.Printf("  └── %s/\n", branchName)
	fmt.Printf("      └── .git/\n")

	return nil
}
