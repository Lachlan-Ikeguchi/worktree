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
	Use:   "worktree",
	Short: "Create and manage Git worktrees and clone repositories",
	Long: `Create and manage Git worktrees at ../<branch> relative to the main repository.
Also clone Git repositories into project_name/[branch]/ structure.

Create a new branch and worktree:
  worktree <branch-name>

Create a worktree from an existing remote branch:
  worktree -r <branch-name>

Create a worktree from an existing local branch:
  worktree -e <branch-name>

List all local branches:
  worktree list

List all remote branches:
  worktree list -r

Delete a worktree directory:
  worktree -d <branch-name>

Merge a branch into main/master and clean up (dry-run by default):
  worktree --merge <branch-name>
  worktree --merge --confirm <branch-name>

Delete branch and worktree (dry-run by default):
  worktree --delete <branch-name>
  worktree --delete --confirm <branch-name>

Clone a repository:
  worktree clone <repository-url>`,
	SilenceUsage: true,
	Args:         cobra.ArbitraryArgs,
	ValidArgsFunction: branchNameCompletion,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip pre-run for help, completion, and completion-related calls
		if cmd.Name() == "completion" || (len(args) > 0 && args[0] == "completion") || (len(args) > 0 && args[0] == "__completeNoDesc") {
			return
		}
		// Note: Git repo checks are now done in individual command Run functions
		// This allows clone command to work from non-git directories
	},
}

// Flags for create mode
var (
	remoteFlag    bool
	existingFlag  bool
	deleteFlag    bool
	mergeMode     bool
	deleteMode    bool
	confirmFlag   bool
)

// Flags for list command
var listRemoteFlag bool

func init() {
	// Global flags for create/delete/merge operations
	rootCmd.PersistentFlags().BoolVarP(&remoteFlag, "remote", "r", false, "Create local tracking branch from origin/<branch>")
	rootCmd.PersistentFlags().BoolVarP(&existingFlag, "existing", "e", false, "Create a worktree from an existing local branch")
	rootCmd.PersistentFlags().BoolVarP(&deleteFlag, "delete-worktree", "d", false, "Delete the worktree directory")
	rootCmd.PersistentFlags().BoolVar(&mergeMode, "merge", false, "Merge the branch into main and clean up")
	rootCmd.PersistentFlags().BoolVar(&deleteMode, "delete", false, "Delete the branch, remote branch, and worktree")
	rootCmd.PersistentFlags().BoolVar(&confirmFlag, "confirm", false, "Confirm the merge or delete operation")

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all worktrees",
		Long:  "List all local or remote branches\n\nList all local branches:\n  worktree list\n\nList all remote branches:\n  worktree list -r",
		Run: func(cmd *cobra.Command, args []string) {
			// Check if we're in the main repo (has .git/) not a worktree (has .git file)
			if isWorktree() {
				fmt.Fprintln(os.Stderr, "WARNING: must be called from the main repository, not a worktree")
				os.Exit(1)
			}

			if !isGitRepo() {
				fmt.Fprintln(os.Stderr, "WARNING: not a git repository")
				os.Exit(1)
			}

			var gitArgs []string
			if listRemoteFlag {
				gitArgs = []string{"branch", "-r"}
			} else {
				gitArgs = []string{"branch"}
			}
			cmdExec := exec.Command("git", gitArgs...)
			cmdExec.Stdout = os.Stdout
			cmdExec.Stderr = os.Stderr
			if err := cmdExec.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running git branch: %v\n", err)
				os.Exit(1)
			}
		},
	}
	listCmd.Flags().BoolVarP(&listRemoteFlag, "remote", "r", false, "List all remote branches")
	rootCmd.AddCommand(listCmd)

	// Clone command - moved from separate clone script
	cloneCmd := &cobra.Command{
		Use:   "clone <repository-url>",
		Short: "Clone a Git repository into project_name/[branch]/",
		Long: `Clone a Git repository into a directory structure like:
  project_name/[branch]/

Where [branch] is the repository's default branch (main, master, or trunk).

Examples:
  worktree clone https://github.com/user/repo.git
  worktree clone git@github.com:user/repo.git

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
	rootCmd.AddCommand(cloneCmd)

	// Completion command (automatically generated by Cobra)
	rootCmd.AddCommand(completionCmd)

	// Main run function for create/delete/merge operations
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}

		// Check if we're in the main repo (has .git/) not a worktree (has .git file)
		if isWorktree() {
			fmt.Fprintln(os.Stderr, "WARNING: must be called from the main repository, not a worktree")
			os.Exit(1)
		}

		if !isGitRepo() {
			fmt.Fprintln(os.Stderr, "WARNING: not a git repository")
			os.Exit(1)
		}

		// Validate flag combinations
		if mergeMode && deleteMode {
			fmt.Fprintln(os.Stderr, "WARNING: cannot use --merge and --delete together")
			os.Exit(1)
		}

		if (mergeMode || deleteMode) && (remoteFlag || existingFlag || deleteFlag) {
			fmt.Fprintln(os.Stderr, "WARNING: cannot use -r, -e, or -d with --merge or --delete")
			os.Exit(1)
		}

		if deleteFlag && (remoteFlag || existingFlag) {
			fmt.Fprintln(os.Stderr, "WARNING: cannot use -r or -e with -d")
			os.Exit(1)
		}

		if deleteFlag && confirmFlag {
			fmt.Fprintln(os.Stderr, "WARNING: --confirm is only used with --merge or --delete")
			os.Exit(1)
		}

		if confirmFlag && !mergeMode && !deleteMode {
			fmt.Fprintln(os.Stderr, "WARNING: --confirm is only used with --merge or --delete")
			os.Exit(1)
		}

		if remoteFlag && existingFlag {
			fmt.Fprintln(os.Stderr, "WARNING: cannot use switches -r and -e at the same time")
			os.Exit(1)
		}

		if deleteFlag && (mergeMode || deleteMode) {
			fmt.Fprintln(os.Stderr, "WARNING: -d cannot be used with create mode")
			os.Exit(1)
		}

		branch := args[0]

		if mergeMode || deleteMode {
			mergeOrDelete(branch, mergeMode, deleteMode, confirmFlag)
			return
		}

		if deleteFlag {
			deleteWorktree(branch)
			return
		}

		// CREATE MODE
		createWorktree(branch, remoteFlag, existingFlag)
	}
}

// completionCmd generates shell completion scripts
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion scripts for your shell",
	Long: `To load completions:

Bash:
  $ source <(worktree completion bash)
  # To load completions for each session, execute once:
  $ worktree completion bash > /etc/bash_completion.d/worktree

Zsh:
  $ source <(worktree completion zsh)
  # To load completions for each session, execute once:
  $ worktree completion zsh > "${fpath[1]}/_worktree"

Fish:
  $ worktree completion fish | source
  # To load completions for each session, execute once:
  $ worktree completion fish > ~/.config/fish/completions/worktree.fish

PowerShell:
  $ worktree completion powershell | Out-String | Invoke-Expression
  # To load completions for each session, add to your profile:
  $ worktree completion powershell > $PROFILE.CurrentUserCurrentHost`,
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Args:     cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Completion command doesn't need git repo checks
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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

// getLocalBranches returns a list of all local branch names
func getLocalBranches() ([]string, error) {
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

// getRemoteBranches returns a list of all remote branch names (without origin/ prefix)
func getRemoteBranches() ([]string, error) {
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

// branchNameCompletion returns available branch names for completion
// Returns both local and remote branches to provide comprehensive completion
func branchNameCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// We only complete the first positional argument (branch name)
	// If there are already positional args, don't provide completions
	if len(args) > 0 {
		// Check if the last arg is a flag (starts with -)
		// If so, we might be completing a branch name after a flag
		lastArg := args[len(args)-1]
		if strings.HasPrefix(lastArg, "-") {
			// Last argument is a flag, so we're completing the branch name
			// This is fine, continue
		} else {
			// We already have a positional argument, don't complete
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
	}

	// Get both local and remote branches
	// This provides the most comprehensive completion experience
	seen := make(map[string]bool)
	var allBranches []string

	// Get local branches
	localBranches, err := getLocalBranches()
	if err == nil {
		for _, b := range localBranches {
			if !seen[b] {
				seen[b] = true
				allBranches = append(allBranches, b)
			}
		}
	}

	// Get remote branches
	remoteBranches, err := getRemoteBranches()
	if err == nil {
		for _, b := range remoteBranches {
			if !seen[b] {
				seen[b] = true
				allBranches = append(allBranches, b)
			}
		}
	}

	if len(allBranches) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Filter branches based on the prefix being completed
	var completions []string
	for _, branch := range allBranches {
		if strings.HasPrefix(branch, toComplete) {
			completions = append(completions, branch)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func mergeOrDelete(branch string, mergeMode, deleteMode, confirm bool) {
	mainBranch, err := getMainBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: %v\n", err)
		os.Exit(1)
	}

	if !branchExists(branch) {
		fmt.Fprintf(os.Stderr, "WARNING: branch '%s' does not exist\n", branch)
		os.Exit(1)
	}

	worktreePath := filepath.Join("..", branch)

	if confirm {
		if mergeMode {

			// Perform merge
			cmd := exec.Command("git", "merge", branch)
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: merge failed: %v\n", err)
				os.Exit(1)
			}
		}

		// Remove worktree
		cmd := exec.Command("git", "worktree", "remove", worktreePath)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: could not remove worktree: %v\n", err)
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
			fmt.Fprintf(os.Stderr, "WARNING: failed to delete local branch: %v\n", err)
			os.Exit(1)
		}

		// Delete remote branch
		cmd = exec.Command("git", "push", "-d", "origin", branch)
		if err := cmd.Run(); err != nil {
			fmt.Println("Warning: remote branch may not exist or deletion failed.")
		}
	} else {
		// Dry run - test if operations are actually possible
		typeStr := "Delete"
		flagStr := "--delete"
		if mergeMode {
			typeStr = "Merge"
			flagStr = "--merge"
		}

		fmt.Println("=== DRY RUN - Testing if operations are possible ===")
		fmt.Println()
		
		allPossible := true

		if mergeMode {
			// Check branch status relative to main
			fmt.Printf("Checking: Branch %s status relative to %s\n", branch, mainBranch)
			cmd := exec.Command("git", "log", "--oneline", fmt.Sprintf("%s..%s", mainBranch, branch))
			output, err := cmd.Output()
			if err != nil {
				fmt.Printf("  WARNING: Could not check branch status: %v\n", err)
			} else {
				commitOutput := strings.TrimSpace(string(output))
				if commitOutput == "" {
					fmt.Printf("  INFO: Branch %s has no new commits (may be behind or equal to %s)\n", branch, mainBranch)
				} else {
					commits := strings.Split(commitOutput, "\n")
					fmt.Printf("  INFO: Branch %s has %d commit(s) ahead of %s\n", branch, len(commits), mainBranch)
				}
			}

			// Check if branch is behind main
			cmd = exec.Command("git", "log", "--oneline", fmt.Sprintf("%s..%s", branch, mainBranch))
			output, err = cmd.Output()
			if err != nil {
				fmt.Printf("  WARNING: Could not check if branch is behind: %v\n", err)
			} else {
				commitOutput := strings.TrimSpace(string(output))
				if commitOutput != "" {
					commits := strings.Split(commitOutput, "\n")
					fmt.Fprintf(os.Stderr, "WARNING: Branch %s is %d commit(s) behind %s\n", branch, len(commits), mainBranch)
				} else {
					fmt.Printf("  PASS: Branch %s is up to date with %s\n", branch, mainBranch)
				}
			}

			// Test merge
			fmt.Printf("Testing: %s %s onto %s\n", typeStr, branch, mainBranch)
			cmd = exec.Command("git", "merge", "--no-commit", "--no-ff", branch)
			if err := cmd.Run(); err != nil {
				fmt.Printf("  FAIL: Merge would fail: %v\n", err)
				allPossible = false
			} else {
				// Reset the merge (we only tested with --no-commit)
				cmd := exec.Command("git", "merge", "--abort")
				cmd.Run() // Ignore error, might not have started
				fmt.Printf("  PASS: Merge is possible\n")
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
			fmt.Printf("  FAIL: Worktree directory not found at %s\n", worktreePath)
			allPossible = false
		} else {
			// Check if worktree is valid and can be removed
			cmd := exec.Command("git", "worktree", "list")
			output, err := cmd.Output()
			if err != nil {
				fmt.Printf("  FAIL: Could not list worktrees: %v\n", err)
				allPossible = false
			} else {
				worktreeList := string(output)
				// Check both relative and absolute paths
				if !strings.Contains(worktreeList, worktreePath) && !strings.Contains(worktreeList, absWorktreePath) {
					fmt.Printf("  FAIL: Worktree at %s is not registered\n", worktreePath)
					allPossible = false
				} else {
					fmt.Printf("  PASS: Worktree can be removed\n")
				}
			}
		}

		// Test local branch deletion
		fmt.Printf("Testing: Delete local branch %s\n", branch)
		if !branchExists(branch) {
			fmt.Printf("  FAIL: Local branch '%s' does not exist\n", branch)
			allPossible = false
		} else {
			// Check if branch is fully merged to main branch
			// A branch can be safely deleted with -d if it's merged
			// Otherwise, we need -D (force)
			cmd := exec.Command("git", "branch", "--merged", mainBranch)
			output, err := cmd.Output()
			if err != nil {
				fmt.Printf("  FAIL: Could not check merged branches: %v\n", err)
				allPossible = false
			} else {
				branches := strings.Split(strings.TrimSpace(string(output)), "\n")
				merged := false
				for _, b := range branches {
					// Remove leading * and spaces
					b = strings.TrimLeft(b, " *")
					if b == branch {
						merged = true
						break
					}
				}
				if merged {
					fmt.Printf("  PASS: Local branch can be deleted (fully merged)\n")
				} else {
					// Branch not merged, but can be force deleted
					fmt.Printf("  PASS: Local branch can be deleted (with force -D)\n")
				}
			}
		}

		// Test remote branch deletion
		fmt.Printf("Testing: Delete remote branch origin/%s\n", branch)
		// First check if remote branch exists
		cmd := exec.Command("git", "show-ref", "--quiet", "refs/remotes/origin/"+branch)
		if err := cmd.Run(); err != nil {
			// Remote branch doesn't exist - that's okay, just warn
			fmt.Printf("  WARNING: Remote branch origin/%s does not exist (will be skipped)\n", branch)
		} else {
			// Remote branch exists, test deletion
			cmd := exec.Command("git", "push", "-d", "--dry-run", "origin", branch)
			if err := cmd.Run(); err != nil {
				fmt.Printf("  FAIL: Remote branch deletion would fail: %v\n", err)
				allPossible = false
			} else {
				fmt.Printf("  PASS: Remote branch can be deleted\n")
			}
		}

		fmt.Println()
		if allPossible {
			fmt.Println("All operations are possible!")
		} else {
			fmt.Println("Some operations would fail - see above for details")
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
		fmt.Fprintf(os.Stderr, "WARNING: could not remove worktree: %v\n", err)
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
				fmt.Fprintf(os.Stderr, "WARNING: branch '%s' does not exist locally\n", branch)
				os.Exit(1)
			}
		}
	} else {
		// Create the branch
		if remote {
			cmd := exec.Command("git", "fetch", "origin", fmt.Sprintf("%s:%s", branch, branch))
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: failed to create local tracking branch: %v\n", err)
				os.Exit(1)
			}
		} else if !existing {
			cmd := exec.Command("git", "branch", branch)
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: failed to create branch: %v\n", err)
				os.Exit(1)
			}
			fmt.Println()
			fmt.Println("Remember to push to create tracking branches in remote")
			fmt.Println()
		} else {
			if !branchExists(branch) {
				fmt.Fprintf(os.Stderr, "WARNING: branch '%s' does not exist locally\n", branch)
				os.Exit(1)
			}
		}

		// Create worktree
		cmd := exec.Command("git", "worktree", "add", worktreePath, branch)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to create worktree: %v\n", err)
			os.Exit(1)
		}
	}

	// Start a new bash shell in the worktree
	cmd := exec.Command("bash", "-c", fmt.Sprintf("cd %q; exec bash", worktreePath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to start shell: %v\n", err)
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
