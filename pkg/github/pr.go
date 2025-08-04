package github

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	gh "github.com/cli/go-gh/v2"
	"github.com/vladimir-ananiev/gh-stack/pkg/git"
)

const (
	maxTitleLength = 50
)

type PR struct {
	Number         int    `json:"number"`
	Title          string `json:"title"`
	HeadRefName    string `json:"headRefName"`
	BaseRefName    string `json:"baseRefName"`
	State          string `json:"state"`
	IsDraft        bool   `json:"isDraft"`
	Mergeable      string `json:"mergeable"`
	ReviewDecision string `json:"reviewDecision,omitempty"`
}

type TreeNode struct {
	PR       *PR
	Children []*TreeNode
}

var (
	branchStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	baseBranchStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	currentStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	numberStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// GetOpenPRs gets all open PRs for the current repository authored by the current user
func GetOpenPRs(ctx context.Context) ([]*PR, error) {
	output, _, err := gh.ExecContext(ctx, "pr", "list",
		"--json", "number,title,headRefName,baseRefName,state,isDraft,mergeable,reviewDecision",
		"--state", "open",
		"--author", "@me")
	if err != nil {
		return nil, fmt.Errorf("failed to get PRs: %w", err)
	}

	var prs []*PR
	if err := json.Unmarshal(output.Bytes(), &prs); err != nil {
		return nil, fmt.Errorf("failed to parse PR data: %w", err)
	}

	return prs, nil
}

// BuildDependencyTree builds a tree structure from PRs based on branch relationships using topological sorting
func BuildDependencyTree(prs []*PR) []*TreeNode {
	if len(prs) == 0 {
		return []*TreeNode{}
	}

	branchToPR := make(map[string]*PR)
	for _, pr := range prs {
		branchToPR[pr.HeadRefName] = pr
	}

	visited := make(map[string]bool)
	var roots []*TreeNode

	for _, pr := range prs {
		if visited[pr.HeadRefName] {
			continue
		}

		// Check if this PR's base branch is external (not a head branch of any PR)
		if _, exists := branchToPR[pr.BaseRefName]; !exists {
			// This is a root - build its tree
			root := buildSubtree(pr, branchToPR, visited)
			roots = append(roots, root)
		}
	}

	return roots
}

// buildSubtree builds a tree rooted at the given PR using DFS
func buildSubtree(pr *PR, branchToPR map[string]*PR, visited map[string]bool) *TreeNode {
	visited[pr.HeadRefName] = true
	node := &TreeNode{PR: pr}

	var childPRs []*PR
	for _, childPR := range branchToPR {
		if childPR.BaseRefName == pr.HeadRefName && !visited[childPR.HeadRefName] {
			childPRs = append(childPRs, childPR)
		}
	}

	slices.SortFunc(childPRs, func(a, b *PR) int {
		return a.Number - b.Number
	})

	for _, childPR := range childPRs {
		childNode := buildSubtree(childPR, branchToPR, visited)
		node.Children = append(node.Children, childNode)
	}

	return node
}

// PrintTree prints the dependency tree with base branches as roots using lipgloss tree
func PrintTree(roots []*TreeNode, currentBranch string) {
	if len(roots) == 0 {
		fmt.Println("No open PRs found")
		return
	}

	// Group roots by base branch
	branchGroups := make(map[string][]*TreeNode)
	for _, root := range roots {
		baseBranch := root.PR.BaseRefName
		branchGroups[baseBranch] = append(branchGroups[baseBranch], root)
	}

	// Sort base branches for deterministic output
	var sortedBranches []string
	for baseBranch := range branchGroups {
		sortedBranches = append(sortedBranches, baseBranch)
	}
	sort.Strings(sortedBranches)

	// Print each base branch group in sorted order
	for i, baseBranch := range sortedBranches {
		// Add spacing between sections (except for first one)
		if i > 0 {
			fmt.Println()
		}

		// Print base branch without indentation
		fmt.Println(baseBranchStyle.Render(baseBranch))

		// Create tree for PRs under this base branch
		t := tree.New()
		for _, root := range branchGroups[baseBranch] {
			addPRNodeToTree(t, root, currentBranch)
		}

		fmt.Println(t)
	}
}

func formatPRNode(pr *PR, currentBranch string) string {
	status := getStatusIcon(pr)

	branchText := pr.HeadRefName
	if pr.HeadRefName == currentBranch {
		branchText = currentStyle.Render(pr.HeadRefName + " ← current")
	} else {
		branchText = branchStyle.Render(pr.HeadRefName)
	}

	numberText := numberStyle.Render(fmt.Sprintf("#%d", pr.Number))
	title := pr.Title
	if len(title) > maxTitleLength {
		title = title[:maxTitleLength-3] + "..."
	}

	return fmt.Sprintf("%s %s %s %s", status, branchText, numberText, title)
}

func addPRNodeToTree(t *tree.Tree, node *TreeNode, currentBranch string) {
	nodeText := formatPRNode(node.PR, currentBranch)

	if len(node.Children) == 0 {
		t.Child(nodeText)
	} else {
		childTree := tree.Root(nodeText)
		for _, child := range node.Children {
			addPRNodeToTree(childTree, child, currentBranch)
		}
		t.Child(childTree)
	}
}

func getStatusIcon(pr *PR) string {
	if pr.IsDraft {
		return "📝"
	}
	if pr.Mergeable == "CONFLICTING" {
		return "⚠️"
	}
	if pr.ReviewDecision == "APPROVED" {
		return "✅"
	}
	if pr.ReviewDecision == "CHANGES_REQUESTED" {
		return "❌"
	}
	return "🔄"
}

// FindCurrentBranchTree finds the tree containing the current branch
func FindCurrentBranchTree(roots []*TreeNode, currentBranch string) *TreeNode {
	for _, root := range roots {
		if found := findBranchInNode(root, currentBranch); found != nil {
			return found
		}
	}
	return nil
}

func findBranchInNode(node *TreeNode, targetBranch string) *TreeNode {
	if node.PR.HeadRefName == targetBranch {
		return node
	}

	for _, child := range node.Children {
		if found := findBranchInNode(child, targetBranch); found != nil {
			return node // Return the root of the tree containing the target
		}
	}

	return nil
}

// ProcessSingleTreeRebase processes a single tree in dependency order for cascading rebase
func ProcessSingleTreeRebase(ctx context.Context, root *TreeNode) error {
	return processNodeRebase(ctx, root)
}

func processNodeRebase(ctx context.Context, node *TreeNode) error {
	fmt.Printf("🔄 Processing %s -> %s...\n",
		branchStyle.Render(node.PR.HeadRefName),
		branchStyle.Render(node.PR.BaseRefName))

	if err := git.CheckoutBranch(ctx, node.PR.HeadRefName); err != nil {
		return fmt.Errorf("failed to checkout %s: %w", node.PR.HeadRefName, err)
	}

	if err := git.RebaseOnto(ctx, node.PR.BaseRefName); err != nil {
		return err // Error already formatted in RebaseOnto
	}

	if err := git.PushBranch(ctx); err != nil {
		return fmt.Errorf("failed to push %s: %w", node.PR.HeadRefName, err)
	}

	fmt.Printf("✅ Completed %s\n", branchStyle.Render(node.PR.HeadRefName))

	for _, child := range node.Children {
		if err := processNodeRebase(ctx, child); err != nil {
			return err
		}
	}

	return nil
}
