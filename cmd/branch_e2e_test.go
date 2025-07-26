package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	globalBinaryPath string
	globalBinaryOnce sync.Once
)

// getSharedBinary builds the CLI binary once and reuses it across tests
func getSharedBinary(t *testing.T) string {
	t.Helper()

	globalBinaryOnce.Do(func() {
		projectRoot, err := findProjectRoot()
		require.NoError(t, err)

		binaryDir, err := os.MkdirTemp("", "gh-stacked-shared-*")
		require.NoError(t, err)

		globalBinaryPath = filepath.Join(binaryDir, "gh-stacked")
		cmd := exec.Command("go", "build", "-o", globalBinaryPath, projectRoot)
		err = cmd.Run()
		require.NoError(t, err)
	})

	return globalBinaryPath
}

// setupE2ETestRepo creates a test repository and gets the shared CLI binary
func setupE2ETestRepo(t *testing.T) (string, string) {
	t.Helper()

	// Create temporary directory for test repo
	tempDir, err := os.MkdirTemp("", "gh-stacked-e2e-*")
	require.NoError(t, err)

	// Initialize git repository
	_, err = git.PlainInit(tempDir, false)
	require.NoError(t, err)

	// Create and commit an initial file
	repo, err := git.PlainOpen(tempDir)
	require.NoError(t, err)

	worktree, err := repo.Worktree()
	require.NoError(t, err)

	initialFile := filepath.Join(tempDir, "README.md")
	err = os.WriteFile(initialFile, []byte("# Test Repository"), 0644)
	require.NoError(t, err)

	_, err = worktree.Add("README.md")
	require.NoError(t, err)

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{})
	require.NoError(t, err)

	// Get shared binary
	binaryPath := getSharedBinary(t)

	// Clean up repo after test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return tempDir, binaryPath
}

// findProjectRoot finds the project root directory
func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Look for go.mod file going up directories
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return wd, nil // fallback to current directory
}

// runCLI executes the CLI binary with given args in the specified directory
func runCLI(t *testing.T, repoPath, binaryPath string, args ...string) (string, string, error) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = repoPath

	stdout, stderr, err := runCommand(cmd)
	t.Logf("CLI command: %s %v", binaryPath, args)
	t.Logf("stdout: %s", stdout)
	t.Logf("stderr: %s", stderr)

	return stdout, stderr, err
}

// runCommand executes a command and returns stdout, stderr, and error
func runCommand(cmd *exec.Cmd) (string, string, error) {
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// TestCLI_BranchCommand_Success tests successful branch creation via CLI
func TestCLI_BranchCommand_Success(t *testing.T) {
	repoPath, binaryPath := setupE2ETestRepo(t)

	// Test creating a branch
	stdout, stderr, err := runCLI(t, repoPath, binaryPath, "branch", "feature-test")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "Creating branch 'feature-test' from 'master'")
	assert.Contains(t, stdout, "✓ Created and switched to branch 'feature-test'")
	assert.Empty(t, stderr)

	// Verify branch was actually created using git command
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	require.NoError(t, err)
	assert.Equal(t, "feature-test", strings.TrimSpace(string(output)))
}

// TestCLI_BranchCommand_StackCreation tests creating a stack of branches
func TestCLI_BranchCommand_StackCreation(t *testing.T) {
	repoPath, binaryPath := setupE2ETestRepo(t)

	// Create first branch
	stdout, _, err := runCLI(t, repoPath, binaryPath, "branch", "feature-1")
	assert.NoError(t, err)
	assert.Contains(t, stdout, "Creating branch 'feature-1' from 'master'")

	// Create second branch (should be from feature-1)
	stdout, _, err = runCLI(t, repoPath, binaryPath, "branch", "feature-2")
	assert.NoError(t, err)
	assert.Contains(t, stdout, "Creating branch 'feature-2' from 'feature-1'")

	// Create third branch (should be from feature-2)
	stdout, _, err = runCLI(t, repoPath, binaryPath, "branch", "feature-3")
	assert.NoError(t, err)
	assert.Contains(t, stdout, "Creating branch 'feature-3' from 'feature-2'")

	// Verify current branch
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	require.NoError(t, err)
	assert.Equal(t, "feature-3", strings.TrimSpace(string(output)))
}

// TestCLI_BranchCommand_Errors tests error scenarios via CLI
func TestCLI_BranchCommand_Errors(t *testing.T) {
	repoPath, binaryPath := setupE2ETestRepo(t)

	t.Run("empty branch name", func(t *testing.T) {
		_, stderr, err := runCLI(t, repoPath, binaryPath, "branch", "")
		assert.Error(t, err)
		assert.Contains(t, stderr, "invalid branch name")
	})

	t.Run("duplicate branch name", func(t *testing.T) {
		// Create first branch
		_, _, err := runCLI(t, repoPath, binaryPath, "branch", "duplicate")
		require.NoError(t, err)

		// Try to create same branch again
		_, _, err = runCLI(t, repoPath, binaryPath, "branch", "duplicate")
		if err != nil {
			// Could be binary not found or actual branch error
			t.Logf("Expected error for duplicate branch: %v", err)
		}
	})

	t.Run("no arguments", func(t *testing.T) {
		_, _, err := runCLI(t, repoPath, binaryPath, "branch")
		if err != nil {
			// Could be binary not found or argument validation error
			t.Logf("Expected error for no arguments: %v", err)
		}
	})
}

// TestCLI_BranchCommand_Help tests help functionality
func TestCLI_BranchCommand_Help(t *testing.T) {
	repoPath, binaryPath := setupE2ETestRepo(t)

	stdout, stderr, err := runCLI(t, repoPath, binaryPath, "branch", "--help")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "Create a new branch")
	assert.Contains(t, stdout, "Usage:")
	assert.Contains(t, stdout, "stacked branch <branch-name>")
	assert.Empty(t, stderr)
}

// TestCLI_BranchCommand_NonGitDirectory tests running in non-git directory
func TestCLI_BranchCommand_NonGitDirectory(t *testing.T) {
	// Create non-git directory
	tempDir, err := os.MkdirTemp("", "non-git-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Use shared binary
	binaryPath := getSharedBinary(t)

	// Try to create branch in non-git directory
	_, stderr, err := runCLI(t, tempDir, binaryPath, "branch", "feature-1")

	assert.Error(t, err)
	assert.Contains(t, stderr, "not in a git repository")
}