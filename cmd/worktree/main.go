package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/lachlan/worktree/internal/clone"
	"github.com/lachlan/worktree/internal/completion"
	"github.com/lachlan/worktree/internal/git"
	"github.com/lachlan/worktree/internal/list"
	"github.com/lachlan/worktree/internal/project"
	"github.com/lachlan/worktree/internal/worktree"
)

// Global flags for create/delete/merge operations
var (
	remoteFlag    bool
	existingFlag  bool
	deleteFlag    bool
	mergeMode     bool
	deleteMode    bool
	confirmFlag   bool
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

List all branches (local and remote):
  worktree list

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
	ValidArgsFunction: completion.BranchNameCompletion,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip pre-run for help, completion, and completion-related calls
		if cmd.Name() == "completion" || (len(args) > 0 && args[0] == "completion") || (len(args) > 0 && args[0] == "__completeNoDesc") {
			return
		}
		// Note: Git repo checks are now done in individual command Run functions
		// This allows clone command to work from non-git directories
	},
}

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
		Long:  "List all branches (local and remote)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if we're in the main repo (has .git/) not a worktree (has .git file)
			if git.IsWorktree() {
				return fmt.Errorf("must be called from the main repository, not a worktree")
			}

			if !git.IsGitRepo() {
				return fmt.Errorf("not a git repository")
			}

			return list.List()
		},
	}
	rootCmd.AddCommand(listCmd)

	// Clone command
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return clone.Clone(args[0])
		},
	}
	rootCmd.AddCommand(cloneCmd)

	// Init command
	initCmd := &cobra.Command{
		Use:   "init <project-name>",
		Short: "Initialize a new project with worktree structure",
		Long: `Initialize a new project directory with the worktree structure.

Creates a project directory with a subdirectory named after the main branch
(master/main/trunk) and initializes a git repository in that subdirectory.

The resulting structure will be:
  <current-dir>/
  └── <project-name>/
      └── <main-branch>/
          └── .git/ (newly initialized git repository)

Examples:
  worktree init myproject
  worktree init test_project

The main branch name is determined from git config (init.defaultBranch) or
defaults to 'master' if not configured.`,
		SilenceUsage: true,
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return project.Init(args[0])
		},
	}
	rootCmd.AddCommand(initCmd)

	// Completion command
	rootCmd.AddCommand(completion.GetCompletionCommand())

	// Main run function for create/delete/merge operations
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		// Check if we're in the main repo (has .git/) not a worktree (has .git file)
		if git.IsWorktree() {
			fmt.Fprintln(os.Stderr, "WARNING: must be called from the main repository, not a worktree")
			os.Exit(1)
		}

		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "WARNING: not a git repository")
			os.Exit(1)
		}

		if len(args) == 0 {
			cmd.Help()
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
			if err := worktree.MergeOrDelete(branch, mergeMode, deleteMode, confirmFlag); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: %v\n", err)
				os.Exit(1)
			}
			return
		}

		if deleteFlag {
			if err := worktree.DeleteWorktreeOnly(branch); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// CREATE MODE
		if err := worktree.CreateWorktree(branch, remoteFlag, existingFlag); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: %v\n", err)
			os.Exit(1)
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
