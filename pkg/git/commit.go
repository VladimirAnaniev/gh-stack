package git

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/vladimir-ananiev/gh-stacked/pkg/stack"
)

func AmendCommitWithPRNumber(commit Commit, pr int, parentPr *int) error {
	repo, err := getRepo()
	if err != nil {
		return stack.ErrNotInRepository
	}

	// Get the commit object
	commitObj, err := repo.CommitObject(plumbing.NewHash(commit.Hash))
	if err != nil {
		return fmt.Errorf("error getting commit: %w", err)
	}

	// Create new commit message with PR metadata
	newMessage := addPRMetadataToMessage(commitObj.Message, pr, parentPr)

	// Get the current working tree
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree: %w", err)
	}

	// Amend the commit
	_, err = worktree.Commit(newMessage, &git.CommitOptions{
		Amend: true,
		Author: &object.Signature{
			Name:  commitObj.Author.Name,
			Email: commitObj.Author.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("error amending commit: %w", err)
	}

	return nil
}

func addPRMetadataToMessage(message string, pr int, parentPr *int) string {
	message = trimMetadata(message)
	return fmt.Sprintf("%s\n\n%s", message, metadataString(pr, parentPr))
}

func trimMetadata(message string) string {
	return strings.Split(message, "gh-stacked:")[0]
}

func metadataString(pr int, parentPr *int) string {
	metadata := fmt.Sprintf("gh-stacked: pr=%d", pr)
	if parentPr != nil {
		metadata += fmt.Sprintf(" parent-pr=%d", *parentPr)
	}
	return metadata
}

type PRMetadata struct {
	PR       int
	ParentPR *int
}

func ParsePRMetadataFromMessage(message string) *PRMetadata {
	re := regexp.MustCompile(`gh-stacked:\s*pr=(\d+)(?:\s+parent-pr=(\d+))?`)
	matches := re.FindStringSubmatch(message)
	if len(matches) < 2 {
		return nil
	}

	pr, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}

	metadata := &PRMetadata{PR: pr}
	
	if len(matches) > 2 && matches[2] != "" {
		parentPR, err := strconv.Atoi(matches[2])
		if err == nil {
			metadata.ParentPR = &parentPR
		}
	}

	return metadata
}
