package completion

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lachlan/worktree/internal/git"
)

// BranchNameCompletion returns available branch names for completion.
// It returns both local and remote branches to provide comprehensive completion.
func BranchNameCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	localBranches, err := git.GetLocalBranches()
	if err == nil {
		for _, b := range localBranches {
			if !seen[b] {
				seen[b] = true
				allBranches = append(allBranches, b)
			}
		}
	}

	// Get remote branches
	remoteBranches, err := git.GetRemoteBranches()
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

// GetCompletionCommand returns the cobra command for completion.
func GetCompletionCommand() *cobra.Command {
	return &cobra.Command{
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
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	}
}
