package git

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/vladimir-ananiev/gh-stacked/pkg/stack"
)

// BranchService handles Git branch operations
type BranchService struct {
	repoPath string
}

// NewBranchService creates a new BranchService
func NewBranchService(repoPath string) *BranchService {
	return &BranchService{
		repoPath: repoPath,
	}
}

// CreateBranch creates a new branch from the specified parent branch and switches to it
func (s *BranchService) CreateBranch(name, parentBranch string) error {
	if strings.TrimSpace(name) == "" {
		return stack.ErrInvalidBranch
	}

	repo, err := git.PlainOpen(s.repoPath)
	if err != nil {
		return stack.ErrNotInRepository
	}

	// Check if branch already exists
	_, err = repo.Reference(plumbing.NewBranchReferenceName(name), true)
	if err == nil {
		return stack.ErrBranchExists
	}

	// Get the parent branch reference
	parentRef, err := repo.Reference(plumbing.NewBranchReferenceName(parentBranch), true)
	if err != nil {
		return stack.ErrParentNotFound
	}

	// Create new branch reference
	branchRef := plumbing.NewBranchReferenceName(name)
	ref := plumbing.NewHashReference(branchRef, parentRef.Hash())
	err = repo.Storer.SetReference(ref)
	if err != nil {
		return err
	}

	// Checkout the new branch
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	return worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
	})
}

// GetCurrentBranch returns the name of the currently checked out branch
func (s *BranchService) GetCurrentBranch() (string, error) {
	repo, err := git.PlainOpen(s.repoPath)
	if err != nil {
		return "", stack.ErrNotInRepository
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	if !head.Name().IsBranch() {
		return "", stack.ErrNotInRepository
	}

	return head.Name().Short(), nil
}
