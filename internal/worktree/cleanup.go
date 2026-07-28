package worktree

import (
	"os"
	"path/filepath"
)

// CleanupEmptyDirs removes empty parent directories starting from the given path.
// It walks up the directory tree until it finds a non-empty directory or reaches the root.
func CleanupEmptyDirs(path string) {
	cleanupPath := path
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
}
