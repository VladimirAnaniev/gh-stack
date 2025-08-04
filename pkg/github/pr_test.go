package github

import (
	"reflect"
	"testing"
)

func TestBuildDependencyTree(t *testing.T) {
	tests := []struct {
		name     string
		prs      []*PR
		expected []*TreeNode
	}{
		{
			name:     "empty PR list",
			prs:      []*PR{},
			expected: []*TreeNode{},
		},
		{
			name: "single PR targeting base branch",
			prs: []*PR{
				{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
			},
			expected: []*TreeNode{
				{
					PR:       &PR{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
					Children: []*TreeNode{},
				},
			},
		},
		{
			name: "linear chain of PRs",
			prs: []*PR{
				{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
				{Number: 2, HeadRefName: "feature-2", BaseRefName: "feature-1"},
				{Number: 3, HeadRefName: "feature-3", BaseRefName: "feature-2"},
			},
			expected: []*TreeNode{
				{
					PR: &PR{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
					Children: []*TreeNode{
						{
							PR: &PR{Number: 2, HeadRefName: "feature-2", BaseRefName: "feature-1"},
							Children: []*TreeNode{
								{
									PR:       &PR{Number: 3, HeadRefName: "feature-3", BaseRefName: "feature-2"},
									Children: []*TreeNode{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "branching tree structure",
			prs: []*PR{
				{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
				{Number: 2, HeadRefName: "feature-2", BaseRefName: "feature-1"},
				{Number: 3, HeadRefName: "feature-3", BaseRefName: "feature-1"},
				{Number: 4, HeadRefName: "feature-4", BaseRefName: "feature-2"},
			},
			expected: []*TreeNode{
				{
					PR: &PR{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
					Children: []*TreeNode{
						{
							PR: &PR{Number: 2, HeadRefName: "feature-2", BaseRefName: "feature-1"},
							Children: []*TreeNode{
								{
									PR:       &PR{Number: 4, HeadRefName: "feature-4", BaseRefName: "feature-2"},
									Children: []*TreeNode{},
								},
							},
						},
						{
							PR:       &PR{Number: 3, HeadRefName: "feature-3", BaseRefName: "feature-1"},
							Children: []*TreeNode{},
						},
					},
				},
			},
		},
		{
			name: "multiple base branches",
			prs: []*PR{
				{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
				{Number: 2, HeadRefName: "feature-2", BaseRefName: "main"},
				{Number: 3, HeadRefName: "feature-3", BaseRefName: "develop"},
			},
			expected: []*TreeNode{
				{
					PR:       &PR{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
					Children: []*TreeNode{},
				},
				{
					PR:       &PR{Number: 2, HeadRefName: "feature-2", BaseRefName: "main"},
					Children: []*TreeNode{},
				},
				{
					PR:       &PR{Number: 3, HeadRefName: "feature-3", BaseRefName: "develop"},
					Children: []*TreeNode{},
				},
			},
		},
		{
			name: "complex mixed structure",
			prs: []*PR{
				{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
				{Number: 2, HeadRefName: "feature-2", BaseRefName: "feature-1"},
				{Number: 3, HeadRefName: "feature-3", BaseRefName: "develop"},
				{Number: 4, HeadRefName: "feature-4", BaseRefName: "feature-3"},
				{Number: 5, HeadRefName: "feature-5", BaseRefName: "main"},
			},
			expected: []*TreeNode{
				{
					PR: &PR{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
					Children: []*TreeNode{
						{
							PR:       &PR{Number: 2, HeadRefName: "feature-2", BaseRefName: "feature-1"},
							Children: []*TreeNode{},
						},
					},
				},
				{
					PR: &PR{Number: 3, HeadRefName: "feature-3", BaseRefName: "develop"},
					Children: []*TreeNode{
						{
							PR:       &PR{Number: 4, HeadRefName: "feature-4", BaseRefName: "feature-3"},
							Children: []*TreeNode{},
						},
					},
				},
				{
					PR:       &PR{Number: 5, HeadRefName: "feature-5", BaseRefName: "main"},
					Children: []*TreeNode{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildDependencyTree(tt.prs)

			if !equalTreeNodes(result, tt.expected) {
				t.Errorf("BuildDependencyTree() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFindCurrentBranchTree(t *testing.T) {
	roots := []*TreeNode{
		{
			PR: &PR{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
			Children: []*TreeNode{
				{
					PR: &PR{Number: 2, HeadRefName: "feature-2", BaseRefName: "feature-1"},
					Children: []*TreeNode{
						{
							PR:       &PR{Number: 3, HeadRefName: "feature-3", BaseRefName: "feature-2"},
							Children: []*TreeNode{},
						},
					},
				},
			},
		},
		{
			PR:       &PR{Number: 4, HeadRefName: "feature-4", BaseRefName: "develop"},
			Children: []*TreeNode{},
		},
	}

	tests := []struct {
		name          string
		currentBranch string
		expectedPR    *PR
	}{
		{
			name:          "find root branch",
			currentBranch: "feature-1",
			expectedPR:    &PR{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
		},
		{
			name:          "find nested branch",
			currentBranch: "feature-3",
			expectedPR:    &PR{Number: 1, HeadRefName: "feature-1", BaseRefName: "main"},
		},
		{
			name:          "find branch in different tree",
			currentBranch: "feature-4",
			expectedPR:    &PR{Number: 4, HeadRefName: "feature-4", BaseRefName: "develop"},
		},
		{
			name:          "branch not found",
			currentBranch: "non-existent",
			expectedPR:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindCurrentBranchTree(roots, tt.currentBranch)

			if tt.expectedPR == nil {
				if result != nil {
					t.Errorf("FindCurrentBranchTree() = %v, want nil", result)
				}
			} else {
				if result == nil || !reflect.DeepEqual(result.PR, tt.expectedPR) {
					t.Errorf("FindCurrentBranchTree() = %v, want %v", result.PR, tt.expectedPR)
				}
			}
		})
	}
}

func TestGetStatusIcon(t *testing.T) {
	tests := []struct {
		name     string
		pr       *PR
		expected string
	}{
		{
			name:     "draft PR",
			pr:       &PR{IsDraft: true},
			expected: "📝",
		},
		{
			name:     "conflicting PR",
			pr:       &PR{Mergeable: "CONFLICTING"},
			expected: "⚠️",
		},
		{
			name:     "approved PR",
			pr:       &PR{ReviewDecision: "APPROVED"},
			expected: "✅",
		},
		{
			name:     "changes requested",
			pr:       &PR{ReviewDecision: "CHANGES_REQUESTED"},
			expected: "❌",
		},
		{
			name:     "default status",
			pr:       &PR{},
			expected: "🔄",
		},
		{
			name:     "draft takes precedence over approval",
			pr:       &PR{IsDraft: true, ReviewDecision: "APPROVED"},
			expected: "📝",
		},
		{
			name:     "conflict takes precedence over approval",
			pr:       &PR{Mergeable: "CONFLICTING", ReviewDecision: "APPROVED"},
			expected: "⚠️",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStatusIcon(tt.pr)
			if result != tt.expected {
				t.Errorf("getStatusIcon() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatPRNode(t *testing.T) {
	pr := &PR{
		Number:      123,
		HeadRefName: "feature-branch",
		Title:       "Add new feature",
	}

	tests := []struct {
		name          string
		currentBranch string
		expectCurrent bool
	}{
		{
			name:          "current branch",
			currentBranch: "feature-branch",
			expectCurrent: true,
		},
		{
			name:          "different branch",
			currentBranch: "other-branch",
			expectCurrent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPRNode(pr, tt.currentBranch)

			if tt.expectCurrent {
				if !containsString(result, "current") {
					t.Errorf("formatPRNode() should contain 'current' for current branch")
				}
			} else {
				if containsString(result, "current") {
					t.Errorf("formatPRNode() should not contain 'current' for non-current branch")
				}
			}

			if !containsString(result, "#123") {
				t.Errorf("formatPRNode() should contain PR number")
			}

			if !containsString(result, "Add new feature") {
				t.Errorf("formatPRNode() should contain PR title")
			}
		})
	}
}

func TestFormatPRNodeTitleTruncation(t *testing.T) {
	pr := &PR{
		Number:      123,
		HeadRefName: "feature-branch",
		Title:       "This is a very long title that should be truncated because it exceeds the maximum length allowed",
	}

	result := formatPRNode(pr, "other-branch")

	if !containsString(result, "...") {
		t.Errorf("formatPRNode() should truncate long titles with '...'")
	}
}

// Helper functions

func equalTreeNodes(a, b []*TreeNode) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !equalTreeNode(a[i], b[i]) {
			return false
		}
	}

	return true
}

func equalTreeNode(a, b *TreeNode) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if !reflect.DeepEqual(a.PR, b.PR) {
		return false
	}

	return equalTreeNodes(a.Children, b.Children)
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
