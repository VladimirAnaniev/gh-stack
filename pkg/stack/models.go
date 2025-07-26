package stack

import "errors"

// BranchNode represents a single node in the branch tree
type BranchNode struct {
	Name      string
	Status    BranchStatus
	Parent    *BranchNode
	Children  []*BranchNode
	CommitSHA string
}

// BranchTree represents the complete stack as a tree structure
type BranchTree struct {
	Root *BranchNode // Usually main/master
}

// BranchStatus represents the current state of a branch
type BranchStatus string

const (
	StatusLocal  BranchStatus = "local"
	StatusPushed BranchStatus = "pushed"
	StatusPR     BranchStatus = "pr"
	StatusMerged BranchStatus = "merged"
)

// BranchService interface for branch operations
type BranchService interface {
	CreateBranch(name, parentBranch string) error
	GetCurrentBranch() (string, error)
}

// Common errors
var (
	ErrBranchExists    = errors.New("branch already exists")
	ErrInvalidBranch   = errors.New("invalid branch name")
	ErrNotInRepository = errors.New("not in a git repository")
	ErrParentNotFound  = errors.New("parent branch not found")
)
