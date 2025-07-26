# `gh stacked cleanup` Command Specification

## Overview
Safely removes merged/closed PR branches and allows abandoning unwanted branches, with safety checks to prevent deletion of branches with dependent PRs.

## Command Signature
```bash
gh stacked cleanup [<branch-name> | --merged | --closed]
```

## Arguments
- `<branch-name>` (optional): Specific branch to abandon (closes PR and deletes branch)

## Flags
- `--merged`: Remove all local branches for merged PRs
- `--closed`: Remove all local branches for closed PRs (merged + closed without merge)

## Behavior

### Core Functionality

#### 1. Discovery Phase
- Reuse discovery logic from `status` command to find all stacked PRs
- Query GitHub API for current PR states (merged, closed, open)
- Identify local branches that correspond to merged/closed PRs

#### 2. Safety Validation
- **No children deletion** - Block deletion if branch has dependent PRs
- **Current branch protection** - Cannot delete the branch you're currently on
- **Confirmation prompts** - Show what will be deleted before proceeding

#### 3. Cleanup Operations
- **Auto cleanup** (`--merged`/`--closed`) - Delete multiple branches safely
- **Manual abandon** (`<branch-name>`) - Close PR and delete specific branch
- **Local branch deletion** - Remove branches from local Git repository

## Examples

### Auto Cleanup Merged Branches
```bash
gh stacked cleanup --merged

Found merged PRs with local branches:
✓ #123 feature-auth (merged 2 days ago)
✓ #124 feature-auth-tests (merged 1 day ago)

This will delete 2 local branches:
- feature-auth
- feature-auth-tests

? Proceed with cleanup? (y/N) y

✓ Deleted branch 'feature-auth'
✓ Deleted branch 'feature-auth-tests'

Cleaned up 2 merged branches.
```

### Auto Cleanup Closed Branches
```bash
gh stacked cleanup --closed

Found closed PRs with local branches:
✓ #123 feature-auth (merged)
❌ #125 feature-experiment (closed without merge)

This will delete 2 local branches:
- feature-auth  
- feature-experiment

? Proceed with cleanup? (y/N) y

✓ Deleted branch 'feature-auth'
✓ Deleted branch 'feature-experiment'

Cleaned up 2 closed branches.
```

### Manual Branch Abandonment
```bash
gh stacked cleanup feature-experiment

Branch: feature-experiment
PR: #125 (open)
Status: This will close the PR and delete the local branch

? Close PR #125 and delete branch 'feature-experiment'? (y/N) y

✓ Closed PR #125
✓ Deleted branch 'feature-experiment'

Branch abandoned successfully.
```

### Merged Branch with Dependents (Automatic Rebase)
```bash
gh stacked cleanup --merged

Found merged PRs with dependents:
✓ #123 feature-auth (merged to main, has 2 dependents)

This will:
1. Delete merged branch 'feature-auth'
2. Switch to main
3. Run 'gh stacked cascade' to rebase orphaned dependents
4. All dependent branches rebased onto main with updated PR bases

? Proceed with cleanup and rebase? (y/N) y

✓ Deleted branch 'feature-auth'
✓ Switched to main
✓ Running cascade to rebase orphaned dependents...
  ├─ Rebasing feature-auth-tests onto main ✓
  └─ Rebasing feature-auth-docs onto main ✓
✓ Updated PR bases: #124, #126 → main

Cleanup completed. 1 branch deleted, 2 dependents rebased via cascade.

**Note**: Rebasing delegated to [CASCADE_COMMAND.md](CASCADE_COMMAND.md) for consistency.
```

### Safety Blocking (Open PR with Dependents)
```bash
gh stacked cleanup feature-auth

❌ Cannot delete branch 'feature-auth'

Reason: Branch has open PR (#123) with dependent PRs
Dependencies:
├─ #124 feature-auth-tests
└─ #126 feature-auth-docs

Merge or close PR #123 first, then run cleanup again.
(Cleanup will automatically rebase dependents when parent PR is merged)
```

## Discovery and Safety Logic

### 1. Discovery Phase
Reuses discovery logic from [STATUS_COMMAND.md](STATUS_COMMAND.md#discovery-process) to:
- Find all stacked PRs and build dependency tree
- Query GitHub API for current PR states 
- Map branches to PR numbers and states

### 2. Safety Checks
```go
type SafetyCheck struct {
    CanDelete    bool
    BlockReason  string
    Dependencies []int  // PR numbers that depend on this branch
}

func CheckSafety(branchName string, tree *PRNode) SafetyCheck {
    // 1. Check if current branch
    // 2. Check for dependent children
    // 3. Check if PR is merged (if merged with dependents, allow cleanup with rebase)
    // 4. Return safety status
}
```

### 3. Branch-to-PR Mapping
```go
// From discovery, map local branches to PR numbers
type BranchPRMapping struct {
    BranchName string
    PRNumber   int
    PRState    string  // "open", "merged", "closed"
    HasChildren bool
}
```

## Cleanup Operations

### Auto Cleanup (--merged/--closed)
```go
func AutoCleanup(state string) error {
    // 1. Discover all branches with matching PR state
    // 2. Filter out branches with safety issues
    // 3. Show confirmation prompt with list
    // 4. Delete approved branches
    // 5. Report results
}
```

### Manual Abandon (branch-name)
```go
func AbandonBranch(branchName string) error {
    // 1. Find PR for branch
    // 2. Run safety checks
    // 3. Show confirmation with PR details
    // 4. Close PR via GitHub API
    // 5. Delete local branch
}
```

### Git Branch Deletion
```bash
# Delete local branch safely
git branch -D <branch-name>

# Verify branch is gone
git branch --list <branch-name>
```

## Validation Rules

### Pre-conditions
- ✅ Must be in a Git repository with GitHub remote
- ✅ Must have GitHub authentication configured
- ✅ Must have at least one argument or flag specified

### Safety Requirements
- ✅ Cannot delete current branch
- ✅ Cannot delete branches with dependent PRs
- ✅ Must confirm destructive operations
- ✅ Branches must exist locally before deletion

### Post-conditions
- ✅ Specified branches removed from local Git repository
- ✅ PRs closed if using manual abandon mode
- ✅ Dependency tree remains valid after cleanup

## Error Handling

### Input Validation
- **No arguments provided**: Show usage help
- **Branch doesn't exist**: Clear error message
- **Invalid flags combination**: Error for conflicting flags

### Safety Violations
- **Current branch deletion**: Error with instruction to switch branches
- **Has dependent PRs**: Error with dependency list and suggestions
- **Branch not found locally**: Warning but continue with other branches

### Git Operations
- **Branch deletion fails**: Show Git error and continue with others
- **Permission issues**: Clear error with troubleshooting steps

### GitHub API
- **PR not found**: Warning but continue with local branch deletion
- **Close PR fails**: Show error but offer to continue with local deletion
- **Authentication issues**: Clear instructions for gh auth setup

## Output Formats

### Confirmation Prompts
```bash
Found merged PRs with local branches:
✓ #123 feature-auth (merged 2 days ago)
✓ #124 feature-auth-tests (merged 1 day ago)

This will delete 2 local branches:
- feature-auth
- feature-auth-tests

? Proceed with cleanup? (y/N)
```

### Success Output
```bash
✓ Closed PR #125
✓ Deleted branch 'feature-experiment'
✓ Deleted branch 'feature-auth'

Cleaned up 2 branches.
```

### Safety Block Output
```bash
❌ Cannot delete branch 'feature-auth'

Reason: Branch has dependent PRs
Dependencies:
├─ #124 feature-auth-tests  
└─ #126 feature-auth-docs

Suggestion: Use 'gh stacked rebase' to move children to main, then retry cleanup.
```

### No Matches Output
```bash
gh stacked cleanup --merged

No merged branches found to clean up.

All local branches are either:
- Still open PRs
- Already cleaned up
- Not managed by gh-stacked
```

## Integration Points

### Status Command Integration
- Reuse discovery logic from status command
- Use same PRNode tree structure
- Share GitHub API client and caching

### Git Integration
- Use `go-git` library for branch operations
- Respect Git configuration and repository state
- Handle branch deletion safely

### GitHub Integration
- Use `go-gh` REST client for PR operations
- Close PRs when abandoning branches
- Query PR states for auto-cleanup decisions

## Technical Implementation Notes

### Key Operations
1. **Discovery reuse** - Leverage status command's discovery logic
2. **Safety validation** - Check for dependencies and current branch
3. **Confirmation prompts** - Interactive confirmation for destructive operations
4. **Batch operations** - Handle multiple branch deletions efficiently
5. **Error recovery** - Continue processing other branches if one fails

### Data Structures
```go
type CleanupCandidate struct {
    BranchName  string
    PRNumber   int
    PRState    string
    CanDelete  bool
    BlockReason string
    Dependencies []int
}

type CleanupResult struct {
    Deleted   []string
    Failed    []string
    Skipped   []string
    PRsClosed []int
}
```

### GitHub API Operations
```bash
# Close PR when abandoning
PATCH /repos/{owner}/{repo}/pulls/{pr_number}
{"state": "closed"}

# Query PR state for auto-cleanup
GET /repos/{owner}/{repo}/pulls/{pr_number}
```

## Success Criteria
- ✅ Safely remove merged/closed branches without breaking dependency trees
- ✅ Prevent accidental deletion of branches with dependent PRs
- ✅ Clear confirmation prompts for all destructive operations
- ✅ Robust error handling that doesn't break on partial failures
- ✅ Integration with existing discovery logic from status command
- ✅ Support both batch cleanup and individual branch abandonment