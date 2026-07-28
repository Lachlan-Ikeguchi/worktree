package list

import (
	"fmt"
	"os"
	"os/exec"
)

// List displays all branches (local and remote) using git branch -a.
func List() error {
	gitArgs := []string{"branch", "-a"}
	cmdExec := exec.Command("git", gitArgs...)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("error running git branch: %v", err)
	}

	return nil
}
