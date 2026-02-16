package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// CheckDirty returns true if the working tree has uncommitted changes
func CheckDirty() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("git status failed: %v, output: %s", err, string(out))
	}
	return len(strings.TrimSpace(string(out))) > 0, nil
}

// CreateTaskBranch creates and checkouts a new branch with a unique name
func CreateTaskBranch() (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	branchName := fmt.Sprintf("task-%s", timestamp)
	
	cmd := exec.Command("git", "checkout", "-b", branchName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git checkout failed: %v, output: %s", err, string(out))
	}
	
	return branchName, nil
}

// CommitChanges commits all changes in the current branch
func CommitChanges(message string) error {
	addCmd := exec.Command("git", "add", ".")
	if out, err := addCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %v, output: %s", err, string(out))
	}
	
	commitCmd := exec.Command("git", "commit", "-m", message)
	if out, err := commitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit failed: %v, output: %s", err, string(out))
	}
	
	return nil
}
