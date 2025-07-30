package git

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	gh "github.com/cli/go-gh/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Commit struct {
	Hash    string
	Message string
}

func getRepo() (*git.Repository, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	repo, err := git.PlainOpenWithOptions(pwd, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, fmt.Errorf("error opening git repository: %w", err)
	}
	return repo, nil
}

// GetCurrentBranch returns the name of the currently checked out branch
func GetCurrentBranch(ctx context.Context) (string, error) {
	repo, err := getRepo()
	if err != nil {
		return "", fmt.Errorf("not in git repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	if !head.Name().IsBranch() {
		return "", fmt.Errorf("not on a branch")
	}

	return head.Name().Short(), nil
}

// GetDefaultBranch returns the default branch using go-gh
func GetDefaultBranch(ctx context.Context) (string, error) {
	output, _, err := gh.ExecContext(ctx, "repo", "view", "--json", "defaultBranchRef")
	if err != nil {
		return "", fmt.Errorf("failed to get default branch: %w", err)
	}
	
	var result struct {
		DefaultBranchRef struct {
			Name string `json:"name"`
		} `json:"defaultBranchRef"`
	}
	
	if err := json.Unmarshal(output.Bytes(), &result); err != nil {
		return "", fmt.Errorf("failed to parse default branch response: %w", err)
	}
	
	return result.DefaultBranchRef.Name, nil
}

// CheckoutBranch checks out a specific branch
func CheckoutBranch(ctx context.Context, branch string) error {
	repo, err := getRepo()
	if err != nil {
		return fmt.Errorf("not in git repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		return fmt.Errorf("failed to checkout %s: %w", branch, err)
	}

	return nil
}

// CheckoutAndPull checks out a branch and pulls latest changes
func CheckoutAndPull(ctx context.Context, branch string) error {
	// First checkout the branch
	if err := CheckoutBranch(ctx, branch); err != nil {
		return err
	}

	// Then pull using git command
	cmd := exec.CommandContext(ctx, "git", "pull")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull %s: %w", branch, err)
	}

	return nil
}

// RebaseOnto rebases current branch onto target branch using git command
func RebaseOnto(ctx context.Context, target string) error {
	cmd := exec.CommandContext(ctx, "git", "rebase", target)
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Rebase conflict detected on %s\n", target)
		fmt.Println("   Please resolve conflicts manually and run 'git rebase --continue'")
		fmt.Println("   Then re-run 'gh stack cascade' to continue")
		return fmt.Errorf("rebase conflict - manual resolution needed")
	}
	return nil
}

// PushBranch pushes current branch to remote with force-with-lease
func PushBranch(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "push", "--force-with-lease")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push branch: %w", err)
	}
	return nil
}