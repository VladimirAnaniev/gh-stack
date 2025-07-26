package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vladimir-ananiev/gh-stacked/pkg/stack"
)

func setupTestRepo(t *testing.T) string {
	t.Helper()

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "gh-stacked-test-*")
	require.NoError(t, err)

	// Initialize git repository
	_, err = git.PlainInit(tempDir, false)
	require.NoError(t, err)

	// Create and commit an initial file so we have a main branch
	repo, err := git.PlainOpen(tempDir)
	require.NoError(t, err)

	worktree, err := repo.Worktree()
	require.NoError(t, err)

	// Create initial file
	initialFile := filepath.Join(tempDir, "README.md")
	err = os.WriteFile(initialFile, []byte("# Test Repository"), 0644)
	require.NoError(t, err)

	// Add and commit
	_, err = worktree.Add("README.md")
	require.NoError(t, err)

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{})
	require.NoError(t, err)

	// Clean up after test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return tempDir
}

func TestBranchService_CreateBranch_ValidBranchName_ShouldCreateAndSwitchToBranch(t *testing.T) {
	repoPath := setupTestRepo(t)
	service := NewBranchService(repoPath)

	err := service.CreateBranch("feature-1", "master")
	assert.NoError(t, err)

	// Should switch to the newly created branch
	currentBranch, err := service.GetCurrentBranch()
	assert.NoError(t, err)
	assert.Equal(t, "feature-1", currentBranch)
}

func TestBranchService_CreateBranch_EmptyBranchName_ShouldReturnError(t *testing.T) {
	repoPath := setupTestRepo(t)
	service := NewBranchService(repoPath)

	err := service.CreateBranch("", "master")
	assert.Equal(t, stack.ErrInvalidBranch, err)
}

func TestBranchService_CreateBranch_ExistingBranch_ShouldReturnError(t *testing.T) {
	repoPath := setupTestRepo(t)
	service := NewBranchService(repoPath)

	err := service.CreateBranch("feature-1", "master") // Create first
	require.NoError(t, err)

	err = service.CreateBranch("feature-1", "master") // Try to create again
	assert.Equal(t, stack.ErrBranchExists, err)
}

func TestBranchService_GetCurrentBranch_ShouldReturnExpectedBranchName(t *testing.T) {
	repoPath := setupTestRepo(t)
	service := NewBranchService(repoPath)

	// Should start on master branch (default for go-git)
	branch, err := service.GetCurrentBranch()
	assert.NoError(t, err)
	assert.Equal(t, "master", branch)
}
