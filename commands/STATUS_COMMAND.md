# `gh stacked status` Command Specification

## Overview
Discovers and displays dependency trees by analyzing commit metadata and GitHub PR information, showing the current state of stacked PRs as an ASCII tree.

## Command Signature
```bash
gh stacked status
```

## Arguments
None - analyzes all branches in repository

## Flags
None - always displays as tree visualization

## Behavior

### Core Functionality

#### 1. Metadata Discovery via Git Log
- Use single `git log` command to find all commits with stacked metadata
- Parse PR numbers and parent-PR relationships from commit messages
- Extract branch information from commit references

#### 2. GitHub Integration
- Query GitHub API for PR details using discovered PR numbers
- Get current PR status (open, draft, merged, closed)
- Fetch PR titles and merge status

#### 3. Tree Construction
- Build dependency tree with main/master as root
- Create PRNode hierarchy from parent-PR relationships
- Identify current branch and mark in tree

#### 4. Tree Rendering
- Display dependency relationships as ASCII tree rooted at main
- Highlight current branch with visual indicators
- Show PR status with icons and helpful context

## Examples

### Single Dependency Tree
```bash
gh stacked status

main
├─ #123 feature-auth ✅ merged
   ├─ #124 feature-auth-tests 🔄 open, ready to merge
   │  └─ #125 feature-auth-cleanup 📝 draft *
   └─ #126 feature-auth-docs ⚠️ conflicts
└─ #127 feature-ui 🔄 open

* = current branch
✅ = merged, 🔄 = open, 📝 = draft, ⚠️ = conflicts

Current: feature-auth-cleanup (#125)
Next: Complete work and run 'gh stacked push'
```

### No Dependencies Found
```bash
gh stacked status

No stacked PRs found.

ℹ To get started:
  gh stacked branch feature-name
  # make commits
  gh stacked push
```

## Discovery Process

### 1. Git Log Discovery
```bash
# Single command to find all stacked commits
git log --all --grep="stacked-gh:" --format="%H %s %D"

# Returns:
# abc123 feat: add auth stacked-gh: pr=123 (origin/feature-auth, feature-auth)
# def456 test: auth tests stacked-gh: pr=124 parent-pr=123 (origin/feature-auth-tests)
```

### 2. Metadata Parsing
```go
// Parse commit messages for stack metadata
type StackMetadata struct {
    CommitHash string
    PR         int  // Own PR number  
    ParentPR   *int // Parent PR number (nil for root)
    BranchName string
}

// Example parsing:
"stacked-gh: pr=124 parent-pr=123" → {PR: 124, ParentPR: &123}
"stacked-gh: pr=123"               → {PR: 123, ParentPR: nil}
```

### 3. GitHub PR Resolution
```bash
# For each discovered PR number, query GitHub API
GET /repos/{owner}/{repo}/pulls/{pr_number}

# Get essential PR details:
# - state, draft status, mergeable status
# - head/base branch names for validation
```

### 4. Tree Building
```go
// Simplified data structure
type PRNode struct {
    PR              *PRInfo  // nil for main branch
    BranchName      string   // branch name
    Children        []*PRNode
    IsCurrentBranch bool
}

// Build tree with main as root
func BuildTree(metadata []StackMetadata, prDetails map[int]PRInfo) *PRNode {
    // 1. Create main branch node as root
    // 2. Add root PRs (no parent-pr) as children of main
    // 3. Recursively add child PRs based on parent-pr relationships
    // 4. Mark current branch in tree
}
```

## Status Indicators

### PR State Icons
- ✅ **Merged** - PR successfully merged
- 🔄 **Open** - PR open and ready for review
- 📝 **Draft** - PR in draft state
- ⚠️ **Conflicts** - PR has merge conflicts
- ❌ **Closed** - PR closed without merging

### Branch Indicators
- `*` - Current branch
- No extra indicators for simplicity

## Tree Rendering Format

### ASCII Tree Structure
```bash
main                           # Root (always main/master)
├─ #123 feature-auth ✅        # Root PR (parent-pr is null)
   ├─ #124 feature-auth-tests 🔄 # Child PR (parent-pr=123)
   │  └─ #125 feature-cleanup 📝 * # Grandchild (parent-pr=124, current)
   └─ #126 feature-docs ⚠️     # Another child (parent-pr=123)
└─ #127 feature-ui 🔄          # Another root PR
```

### Information Layout
- **PR number and branch name** for identification
- **Status icon** for quick visual status
- **Current branch marker** (`*`) for context
- **Tree connectors** (`├─`, `│`, `└─`) for visual hierarchy

## Git Log Command Details

### Basic Discovery
```bash
git log --all --grep="stacked-gh:" --format="%H %s %D"
```

### Performance Optimization (if needed)
```bash
# Add time limit if repositories become large
git log --all --grep="stacked-gh:" --since="60 days ago" --format="%H %s %D"
```

### Parsing Output
```bash
# Example output line:
abc123 feat: add authentication stacked-gh: pr=123 (origin/feature-auth, feature-auth)

# Extract:
# - Commit hash: abc123
# - Message with metadata: "feat: add authentication stacked-gh: pr=123"
# - Branch refs: origin/feature-auth, feature-auth
```

## Validation Rules

### Pre-conditions
- ✅ Must be in a Git repository
- ✅ Must have GitHub remote configured
- ✅ Must have GitHub authentication configured

### Discovery Requirements
- ✅ At least one commit with stacked metadata found
- ✅ GitHub PR numbers must be valid and accessible
- ✅ Dependency relationships must form valid trees (no cycles)

## Error Handling

### Git Operations
- **No Git repository**: Clear error message with setup instructions
- **No commits found**: Show "No stacked PRs found" message
- **Git command fails**: Show underlying Git error

### GitHub API
- **Authentication failed**: Clear instructions for gh auth setup
- **PR not found**: Show warning but continue with other PRs
- **Rate limiting**: Handle with exponential backoff

### Metadata Issues
- **Invalid metadata format**: Show warning, skip invalid entries
- **Circular dependencies**: Detect and report error
- **Orphaned PRs**: Handle PRs whose parents no longer exist

## Integration Points

### Git Integration
- Single `git log` command for all discovery
- Use `go-git` library or shell out to git command
- Parse commit message format and branch references

### GitHub Integration
- Use `go-gh` REST client for PR API calls
- Batch API requests for efficiency
- Handle GitHub API errors gracefully

### UI Integration
- Use Bubbletea/Bubbles for tree visualization
- Apply consistent styling with Lipgloss
- Ensure readable output in various terminal widths

## Technical Implementation Notes

### Key Operations
1. **Git log execution** - Single command to find all relevant commits
2. **Metadata parsing** - Extract PR numbers and relationships from commit messages
3. **GitHub API calls** - Fetch PR details for discovered PR numbers
4. **Tree building** - Construct PRNode hierarchy with main as root
5. **Tree rendering** - Display ASCII tree with status indicators

### Performance Considerations
- **Single Git command** - Much faster than scanning branches individually
- **Lazy GitHub API calls** - Only fetch PRs that were discovered locally
- **Simple data structures** - PRNode tree is lightweight and fast to traverse

### Current Branch Detection
```bash
# Get current branch name
git branch --show-current

# Mark corresponding node in tree during building
```

## Success Criteria
- ✅ Fast discovery using single optimized Git command
- ✅ Accurate dependency tree construction with main as root
- ✅ Clear ASCII tree visualization showing PR hierarchy and status
- ✅ Robust error handling for Git and GitHub API operations
- ✅ Simple, focused interface with no unnecessary complexity
- ✅ Current branch context clearly highlighted in output