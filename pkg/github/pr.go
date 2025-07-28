package github

import (
	"context"
	"encoding/json"
	"fmt"

	gh "github.com/cli/go-gh/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/vladimir-ananiev/gh-stacked/pkg/git"
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
	branchStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	currentStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	numberStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// GetOpenPRs gets all open PRs for the current repository
func GetOpenPRs(ctx context.Context) ([]*PR, error) {
	output, _, err := gh.ExecContext(ctx, "pr", "list", 
		"--json", "number,title,headRefName,baseRefName,state,isDraft,mergeable,reviewDecision",
		"--state", "open")
	if err != nil {
		return nil, fmt.Errorf("failed to get PRs: %w", err)
	}
	
	var prs []*PR
	if err := json.Unmarshal(output.Bytes(), &prs); err != nil {
		return nil, fmt.Errorf("failed to parse PR data: %w", err)
	}
	
	return prs, nil
}

// BuildDependencyTree builds a tree structure from PRs based on branch relationships
func BuildDependencyTree(prs []*PR) []*TreeNode {
	branchToPR := make(map[string]*PR)
	for _, pr := range prs {
		branchToPR[pr.HeadRefName] = pr
	}
	
	var roots []*TreeNode
	processed := make(map[string]bool)
	
	for _, pr := range prs {
		if processed[pr.HeadRefName] {
			continue
		}
		
		if pr.BaseRefName == "main" || pr.BaseRefName == "master" {
			node := &TreeNode{PR: pr}
			buildChildren(node, branchToPR, processed)
			roots = append(roots, node)
		}
	}
	
	for _, pr := range prs {
		if !processed[pr.HeadRefName] {
			node := &TreeNode{PR: pr}
			buildChildren(node, branchToPR, processed)
			roots = append(roots, node)
		}
	}
	
	return roots
}

func buildChildren(node *TreeNode, branchToPR map[string]*PR, processed map[string]bool) {
	processed[node.PR.HeadRefName] = true
	
	for _, pr := range branchToPR {
		if pr.BaseRefName == node.PR.HeadRefName && !processed[pr.HeadRefName] {
			child := &TreeNode{PR: pr}
			buildChildren(child, branchToPR, processed)
			node.Children = append(node.Children, child)
		}
	}
}

// PrintTree prints the dependency tree using lipgloss tree
func PrintTree(roots []*TreeNode, currentBranch string) {
	if len(roots) == 0 {
		fmt.Println("No open PRs found")
		return
	}
	
	t := tree.New()
	
	for _, root := range roots {
		addNodeToTree(t, root, currentBranch)
	}
	
	fmt.Println(t)
}

func addNodeToTree(t *tree.Tree, node *TreeNode, currentBranch string) {
	status := getStatusIcon(node.PR)
	
	branchText := node.PR.HeadRefName
	if node.PR.HeadRefName == currentBranch {
		branchText = currentStyle.Render(node.PR.HeadRefName + " ← current")
	} else {
		branchText = branchStyle.Render(node.PR.HeadRefName)
	}
	
	numberText := numberStyle.Render(fmt.Sprintf("#%d", node.PR.Number))
	title := node.PR.Title
	if len(title) > 50 {
		title = title[:47] + "..."
	}
	
	nodeText := fmt.Sprintf("%s %s %s %s", status, branchText, numberText, title)
	
	if len(node.Children) == 0 {
		t.Child(nodeText)
	} else {
		childTree := tree.Root(nodeText)
		for _, child := range node.Children {
			addChildToTree(childTree, child, currentBranch)
		}
		t.Child(childTree)
	}
}

func addChildToTree(t *tree.Tree, node *TreeNode, currentBranch string) {
	status := getStatusIcon(node.PR)
	
	branchText := node.PR.HeadRefName
	if node.PR.HeadRefName == currentBranch {
		branchText = currentStyle.Render(node.PR.HeadRefName + " ← current")
	} else {
		branchText = branchStyle.Render(node.PR.HeadRefName)
	}
	
	numberText := numberStyle.Render(fmt.Sprintf("#%d", node.PR.Number))
	title := node.PR.Title
	if len(title) > 50 {
		title = title[:47] + "..."
	}
	
	nodeText := fmt.Sprintf("%s %s %s %s", status, branchText, numberText, title)
	
	if len(node.Children) == 0 {
		t.Child(nodeText)
	} else {
		childTree := tree.Root(nodeText)
		for _, child := range node.Children {
			addChildToTree(childTree, child, currentBranch)
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

// ProcessCascadeRebase processes the tree in dependency order for cascading rebase
func ProcessCascadeRebase(ctx context.Context, roots []*TreeNode) error {
	for _, root := range roots {
		if err := processNodeRebase(ctx, root); err != nil {
			return err
		}
	}
	return nil
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