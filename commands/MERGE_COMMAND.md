# `gh stacked merge` Command Specification

## Overview
Merges the current branch's PR and automatically rebases any dependent branches. Only works for branches that are direct children of main/master. For branches in the middle of a dependency tree, requires parent PRs to be merged first.

## Command Signature
```bash
gh stacked merge [flags]
```

## Arguments
None - always operates on current branch

## Flags
- `--dry-run`: Show merge plan without actually merging

## Behavior

### Core Functionality

#### 1. Current Branch Analysis
- Use discovery logic (see [STATUS_COMMAND.md](STATUS_COMMAND.md#discovery-process)) to determine current branch's position in dependency tree
- Check if current branch is a direct child of main/master
- Identify any dependent branches that will need rebasing

#### 2. Merge Validation
- **Direct child of main**: Proceed with merge
- **Middle of tree**: Block with error, require parent to be merged first
- Verify current PR is ready to merge (approved, CI passing, no conflicts)

#### 3. Merge and Cascade Cleanup
- Merge current PR to main/master
- Switch to main and execute `gh stacked cascade` to handle orphaned dependents (see [CASCADE_COMMAND.md](CASCADE_COMMAND.md))
- Rebase all dependent branches onto main/master
- Update dependent PR base branches

## Examples

### Successful Merge (Root Branch)
```bash
# On feature-auth branch (depends on main)
gh stacked merge

✓ Analyzing current branch: feature-auth
✓ Branch is direct child of main
✓ Found 2 dependent branches: feature-auth-tests, feature-auth-cleanup

Checking merge readiness:
✓ #123 feature-auth: ready to merge (approved, CI passed)

This will:
1. Merge #123 feature-auth → main
2. Delete feature-auth branch
3. Switch to main and run 'gh stacked cascade'
4. All dependent branches rebased onto main with updated PR bases

? Proceed with merge? (y/N) y

✓ Merging #123 feature-auth → main
✓ Deleted feature-auth branch
✓ Switched to main
✓ Running cascade to rebase orphaned dependents...
  ├─ Rebasing feature-auth-tests onto main ✓
  └─ Rebasing feature-auth-cleanup onto main ✓
✓ Updated PR bases: #124, #126 → main

Merge completed successfully.
1 PR merged, 2 dependents rebased via cascade.
```

### Blocked Merge (Middle of Tree)
```bash
# On feature-auth-tests branch (depends on feature-auth)
gh stacked merge

✓ Analyzing current branch: feature-auth-tests
❌ Cannot merge feature-auth-tests

Reason: Branch is not a direct child of main
Current dependency chain: main ← feature-auth ← feature-auth-tests

To merge feature-auth-tests:
1. Switch to feature-auth: git checkout feature-auth
2. Merge parent first: gh stacked merge
3. Return and merge: git checkout feature-auth-tests && gh stacked merge

Or merge the entire chain from the root:
git checkout feature-auth && gh stacked merge
```

### Not Ready to Merge
```bash
gh stacked merge

✓ Analyzing current branch: feature-auth
✓ Branch is direct child of main

❌ Cannot merge feature-auth

PR #123 is not ready:
├─ ❌ CI failing (2 checks failed)
├─ ❌ Not approved
└─ ✓ No merge conflicts

Resolution:
1. Fix failing CI checks
2. Request approval from reviewers
3. Re-run merge command when ready
```

### No Dependents
```bash
# On standalone branch
gh stacked merge

✓ Analyzing current branch: feature-standalone  
✓ Branch is direct child of main
✓ No dependent branches found

Checking merge readiness:
✓ #127 feature-standalone: ready to merge

? Merge #127 feature-standalone → main? (y/N) y

✓ Merging #127 feature-standalone → main
✓ Cleaning up feature-standalone branch

Merge completed successfully.
```

### Dry Run
```bash
gh stacked merge --dry-run

Merge plan for feature-auth:

Target: #123 feature-auth → main
├─ Status: ✓ ready (approved, CI passed)
├─ Dependencies: 2 branches will be rebased
│  ├─ feature-auth-tests → main (PR #124)
│  └─ feature-auth-cleanup → main (PR #126)
└─ Actions: merge PR, cleanup branch, rebase dependents

Ready to proceed. Run without --dry-run to execute.
```

## Merge Position Validation

### 1. Dependency Analysis
```go
func AnalyzeBranchPosition(currentBranch string) (BranchPosition, error) {
    // 1. Use discovery to find parent and children
    // 2. Determine if branch is direct child of main/master
    // 3. Return position info and dependent branches
}

type BranchPosition struct {
    IsRootBranch     bool     // Direct child of main/master
    ParentBranch     string   // Parent branch name
    ParentPR         *int     // Parent PR number (if exists)
    DependentBranches []string // Branches that depend on this one
}
```

### 2. Merge Eligibility Rules
- ✅ **Root branch**: Can merge (direct child of main/master)
- ❌ **Middle branch**: Cannot merge (has unmerged parent)
- ✅ **Leaf branch**: Can merge if it's also a root branch

## GitHub Integration

### 1. PR Merge
```bash
# Merge current PR via GitHub API
PUT /repos/{owner}/{repo}/pulls/{pr_number}/merge
{
  "commit_title": "Merge pull request #123 from feature-auth",
  "merge_method": "merge"  # Respect repo settings
}
```

### 2. Automatic Cleanup
After successful merge, triggers [CLEANUP_COMMAND.md](CLEANUP_COMMAND.md#merged-branch-with-dependents-automatic-rebase) logic:
- Delete merged branch
- Rebase all dependent branches onto main
- Update dependent PR base branches
- Handle conflicts interactively (see [CASCADE_COMMAND.md](CASCADE_COMMAND.md#interactive-conflict-resolution))

## Validation Rules

### Pre-conditions
- ✅ Must be in a Git repository with GitHub remote
- ✅ Must have GitHub authentication configured
- ✅ Current branch must have an open PR
- ✅ Current branch must be part of stacked workflow

### Merge Eligibility
- ✅ Current branch must be direct child of main/master
- ✅ Current PR must be approved and have passing CI
- ✅ Current PR must have no merge conflicts
- ✅ User must confirm merge operation

### Post-conditions
- ✅ Current PR merged to main/master
- ✅ Current branch deleted locally and remotely
- ✅ All dependent branches rebased onto main/master
- ✅ Dependent PR base branches updated

## Error Handling

### Position Validation Errors
- **Middle of tree**: Clear error with steps to merge parent first
- **No PR found**: Error if current branch doesn't have associated PR
- **Already merged**: Handle case where PR was already merged externally

### Merge Readiness Errors
- **CI failing**: Show specific failing checks with links
- **Not approved**: List required approvers
- **Merge conflicts**: Guide user to resolve conflicts first

### Rebase Conflicts (Dependents)
- Uses interactive conflict resolution from [CASCADE_COMMAND.md](CASCADE_COMMAND.md#interactive-conflict-resolution)
- Stop and prompt user to resolve conflicts
- Continue after conflicts resolved

## Output Formats

### Success Output
```bash
✓ Merge validation passed
✓ Merging #123 feature-auth → main
✓ Branch cleanup: deleted feature-auth
✓ Cascading to 2 dependent branches...
  ├─ Rebasing feature-auth-tests onto main ✓
  └─ Rebasing feature-auth-cleanup onto main ✓
✓ Updated PR bases: #124, #126 → main

Merge completed successfully.
1 PR merged, 2 dependents rebased.
```

### Position Blocking
```bash
❌ Cannot merge feature-auth-tests

Position: Middle of dependency tree
Chain: main ← #123 feature-auth ← #124 feature-auth-tests

Required: Merge parent PR first
1. git checkout feature-auth
2. gh stacked merge  # Merge #123 first
3. git checkout feature-auth-tests  
4. gh stacked merge  # Then merge #124

Tip: Start from root (feature-auth) to merge entire chain.
```

### Readiness Blocking
```bash
❌ Cannot merge feature-auth

PR #123 readiness issues:
├─ ❌ CI: 2 checks failing
│  ├─ build (failed) - Fix compilation errors
│  └─ test (failed) - Fix failing unit tests  
├─ ❌ Reviews: Needs 1 approval
│  └─ Request review from: @reviewer1, @reviewer2
└─ ✓ Conflicts: None

Fix these issues and try again.
```

## Integration Points

### Status Command Integration
- Reuse dependency discovery and tree building
- Share branch position analysis logic  
- Use same GitHub API client

### Cleanup Command Integration
- Automatically trigger cleanup after successful merge
- Reuse safety validation and rebase logic
- Share conflict resolution patterns

### GitHub Integration
- Use `go-gh` REST client for merge operations
- Respect repository merge method configuration
- Handle authentication and rate limiting

## Technical Implementation Notes

### Key Operations
1. **Position analysis** - Determine if branch can be merged
2. **Readiness validation** - Check PR is ready for merge
3. **Merge execution** - Merge PR and trigger cleanup
4. **Dependent rebase** - Automatically rebase all dependents
5. **Error handling** - Clear guidance for blocked scenarios

### Data Structures
```go
type MergeOperation struct {
    CurrentBranch     string
    PRNumber         int
    IsEligible       bool
    BlockingReason   string
    DependentBranches []string
    MergeTarget      string // Usually "main"
}
```

### Core Git/GitHub Operations
```bash
# Validate branch position
git merge-base feature-auth main
git log --grep="gh-stacked:" feature-auth

# Merge PR
gh api repos/{owner}/{repo}/pulls/{pr}/merge

# Cleanup (via existing cleanup logic)
# Rebase dependents (via existing cascade logic)
```

## Success Criteria
- ✅ Successfully merge root branches (direct children of main)
- ✅ Block middle branches with clear guidance to merge parents first
- ✅ Automatic cascade rebase of dependent branches after merge
- ✅ Proper validation of PR readiness (approved, CI passing, no conflicts)
- ✅ Clear error messages for all blocking scenarios
- ✅ Integration with existing cleanup and cascade rebase logic