package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clone <repository-url>",
	Short: "Clone a Git repository into project_name/[branch]/",
	Long: `Clone a Git repository into a directory structure like:
  project_name/[branch]/

Where [branch] is the repository's default branch (main, master, or trunk).

Examples:
  clone https://github.com/user/repo.git
  clone git@github.com:user/repo.git

The command:
1. Extracts the project name from the repository URL
2. Creates a project directory with the project name
3. Clones the repository into that directory
4. Determines the default branch from the cloned repository
5. Renames the cloned repository directory to the branch name

The resulting structure will be:
  <current-dir>/
  └── project_name/
      └── [branch-name]/
          └── (repository contents)`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires a repository URL argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]

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

		fmt.Printf("Successfully cloned %s into %s\n", repoURL, filepath.Join(projectName, branchName))
	},
}

// completionCmd generates shell completion scripts
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion scripts for your shell",
	Long: `To load completions:

Bash:
  $ source <(clone completion bash)
  # To load completions for each session, execute once:
  $ clone completion bash > /etc/bash_completion.d/clone

Zsh:
  $ source <(clone completion zsh)
  # To load completions for each session, execute once:
  $ clone completion zsh > "${fpath[1]}/_clone"

Fish:
  $ clone completion fish | source
  # To load completions for each session, execute once:
  $ clone completion fish > ~/.config/fish/completions/clone.fish

PowerShell:
  $ clone completion powershell | Out-String | Invoke-Expression
  # To load completions for each session, add to your profile:
  $ clone completion powershell > $PROFILE.CurrentUserCurrentHost`,
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Args:     cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletion(os.Stdout)
		default:
			return fmt.Errorf("unsupported shell: %s", args[0])
		}
	},
}

func init() {
	// Add completion command
	rootCmd.AddCommand(completionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
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