# `gh stacked cascade` Command Specification

## Overview
Cascades local changes from the current branch down through all dependent branches using depth-first traversal. Rebases each dependent branch, resolves PR ID changes, and pushes to remote. Blocks interactively when conflicts occur until user resolves them.

## Command Signature
```bash
gh stacked cascade
```

## Arguments
None - always operates from the current branch as the starting point

## Flags
None - simple command with interactive conflict resolution

## Behavior

### Core Algorithm

#### 1. Build Dependency Tree
- Start from current branch
- Find all branches that depend on current branch (directly or transitively)
- Build complete dependency tree structure

#### 2. Depth-First Traversal
- Traverse dependency tree using DFS
- Process each branch in dependency order (parents before children)

#### 3. For Each Branch: Rebase → Update → Push
- **Rebase**: `git rebase <immediate-parent>`
- **Update PR IDs**: Update any parent-pr references if PR numbers changed
- **Push**: `git push --force-with-lease origin <branch>`

#### 4. Interactive Conflict Resolution
- On conflict: describe the issue and block with interactive prompt
- Wait for user to manually resolve conflicts using standard Git workflow
- User confirms resolution, then continue to next branch
- **Note**: This conflict resolution pattern is reused by other commands that perform rebasing

## Examples

### Successful Cascade
```bash
gh stacked cascade

✓ Building dependency tree from feature-auth...
✓ Found 3 dependent branches:
  ├─ feature-auth-tests (depends on feature-auth)
  ├─ feature-auth-docs (depends on feature-auth-tests)
  └─ feature-auth-cleanup (depends on feature-auth)

✓ Processing branches (depth-first)...

[1/3] feature-auth-tests
  ├─ Rebasing onto feature-auth... ✓
  ├─ Updating PR metadata... ✓
  └─ Pushing to origin... ✓

[2/3] feature-auth-cleanup  
  ├─ Rebasing onto feature-auth... ✓
  ├─ Updating PR metadata... ✓
  └─ Pushing to origin... ✓

[3/3] feature-auth-docs
  ├─ Rebasing onto feature-auth-tests... ✓
  ├─ Updating PR metadata... ✓
  └─ Pushing to origin... ✓

✅ Cascade completed successfully!
3 branches updated and pushed.
```

### Cascade with Conflicts
```bash
gh stacked cascade

✓ Building dependency tree from feature-auth...
✓ Found 2 dependent branches

✓ Processing branches (depth-first)...

[1/2] feature-auth-tests
  ├─ Rebasing onto feature-auth...
❌ Rebase conflicts detected

Conflicts in files:
  • tests/auth_test.go (lines 45-67)
  • src/helpers.go (lines 12-23)

Please resolve these conflicts:
1. Edit the conflicted files above
2. Stage your changes: git add <files>
3. Complete the rebase: git rebase --continue

? Have you resolved the conflicts and completed the rebase? (y/N) 
```

### After User Resolves Conflicts
```bash
? Have you resolved the conflicts and completed the rebase? (y/N) y

✓ Continuing with feature-auth-tests...
  ├─ Rebase completed ✓
  ├─ Updating PR metadata... ✓
  └─ Pushing to origin... ✓

[2/2] feature-auth-docs
  ├─ Rebasing onto feature-auth-tests... ✓
  ├─ Updating PR metadata... ✓
  └─ Pushing to origin... ✓

✅ Cascade completed successfully!
2 branches updated and pushed.
```

### No Dependents Found
```bash
gh stacked cascade

ℹ No dependent branches found for feature-auth.

This branch is either:
- A leaf node (no other branches depend on it)
- Not part of a stacked PR workflow
- The only branch in your stack

Nothing to cascade.
```

## Algorithm Implementation

### 1. Dependency Tree Building
```go
type BranchNode struct {
    Name     string
    PRNumber int
    Children []*BranchNode
    Parent   *BranchNode
}

func BuildDependencyTree(startBranch string) (*BranchNode, error) {
    // 1. Parse all gh-stacked metadata from git log
    // 2. Build tree structure starting from startBranch
    // 3. Return root node with all descendants
}
```

### 2. Depth-First Processing
```go
func ProcessTreeDFS(node *BranchNode) error {
    // Process all children first (depth-first)
    for _, child := range node.Children {
        if err := ProcessBranch(child); err != nil {
            return err
        }
        if err := ProcessTreeDFS(child); err != nil {
            return err
        }
    }
    return nil
}
```

### 3. Branch Processing
```go
func ProcessBranch(branch *BranchNode) error {
    // 1. git checkout <branch.Name>
    // 2. git rebase <branch.Parent.Name>
    // 3. Handle conflicts if any (interactive prompt)
    // 4. Update PR metadata if parent PR changed
    // 5. git push --force-with-lease origin <branch.Name>
}
```

### 4. Interactive Conflict Resolution
```go
func ResolveConflicts(branch string, conflictedFiles []string) error {
    // 1. Display conflict information
    // 2. Show resolution steps
    // 3. Prompt user and wait for confirmation
    // 4. Validate rebase was completed
    // 5. Continue processing
}
```

## PR Metadata Updates

Cascade preserves existing metadata during rebase operations. Metadata is only updated in specific scenarios:

- **Normal cascade**: Metadata remains unchanged 
- **Parent PR merged**: Handled by cleanup command (see [CLEANUP_COMMAND.md](CLEANUP_COMMAND.md#merged-branch-with-dependents-automatic-rebase))

For metadata annotation rules and format, see [PUSH_COMMAND.md](PUSH_COMMAND.md#commit-annotation).

## Conflict Resolution Details

### Conflict Detection
```bash
# After git rebase command
if rebase_exit_code != 0:
    conflicts = git status --porcelain | grep "^UU\|^AA\|^DD"
    display_conflicts(conflicts)
    prompt_user_resolution()
```

### Interactive Resolution Flow
```go
func HandleRebaseConflicts(branch string) error {
    conflicts := DetectConflictedFiles()
    
    fmt.Printf("❌ Rebase conflicts detected in %s\n\n", branch)
    fmt.Println("Conflicts in files:")
    for _, file := range conflicts {
        fmt.Printf("  • %s\n", file)
    }
    
    fmt.Println("\nPlease resolve these conflicts:")
    fmt.Println("1. Edit the conflicted files above")
    fmt.Println("2. Stage your changes: git add <files>")
    fmt.Println("3. Complete the rebase: git rebase --continue")
    
    // Block until user confirms resolution
    for {
        confirmed := PromptConfirmation("Have you resolved the conflicts and completed the rebase?")
        if confirmed {
            if ValidateRebaseCompleted() {
                return nil
            } else {
                fmt.Println("❌ Rebase is not yet completed. Please run 'git rebase --continue' first.")
            }
        } else {
            if PromptConfirmation("Would you like to abort the rebase?") {
                exec.Command("git", "rebase", "--abort").Run()
                return fmt.Errorf("rebase aborted by user")
            }
        }
    }
}
```

## Validation Rules

### Pre-conditions
- ✅ Must be in a Git repository with GitHub remote
- ✅ Current branch must be part of stacked workflow (have gh-stacked metadata)
- ✅ Current branch must have local changes to cascade
- ✅ All dependent branches must exist locally with clean working directories

### Processing Safety
- ✅ Validate each branch exists before processing
- ✅ Check working directory is clean before each rebase
- ✅ Use `--force-with-lease` for safe force pushes
- ✅ Validate rebase completion before continuing

### Post-conditions
- ✅ All dependent branches rebased with current branch's changes
- ✅ All PR metadata updated to reflect any parent PR changes
- ✅ All successful branches pushed to remote
- ✅ Dependency relationships preserved through the cascade

## Error Handling

### Missing Branches
- **Local branch missing**: Clear error if dependent branch doesn't exist locally
- **Recovery guidance**: Show how to fetch or create missing branches

### Rebase Failures
- **Conflict resolution**: Interactive prompts with clear guidance
- **Abort option**: Allow user to abort rebase and stop cascade
- **Validation**: Ensure rebase is actually completed before continuing

### Push Failures
- **Force push rejection**: Handle remote rejection with clear error message
- **Network issues**: Retry logic for temporary network failures
- **Authentication**: Clear guidance for auth issues

### Git Repository State
- **Dirty working directory**: Clear error with cleanup guidance
- **Detached HEAD**: Validate proper branch state before processing
- **Repository corruption**: Handle Git repository issues gracefully

## Output Formats

### Progress Indicator
```bash
✓ Processing branches (depth-first)...

[1/3] feature-auth-tests
  ├─ Rebasing onto feature-auth... ✓ (2 commits applied)
  ├─ Updating PR metadata... ✓ (parent-pr updated: 123→125)  
  └─ Pushing to origin... ✓

[2/3] feature-auth-cleanup
  ├─ Rebasing onto feature-auth... ✓ (2 commits applied)
  ├─ Updating PR metadata... ✓ (no changes needed)
  └─ Pushing to origin... ✓
```

### Conflict Resolution
```bash
❌ Rebase conflicts detected in feature-auth-tests

Conflicts in files:
  • tests/auth_test.go (merge conflict on lines 45-67)
  • src/helpers.go (both modified lines 12-23)

Please resolve these conflicts:
1. Edit the conflicted files above  
2. Stage your changes: git add <files>
3. Complete the rebase: git rebase --continue

? Have you resolved the conflicts and completed the rebase? (y/N) _
```

### Final Summary
```bash
✅ Cascade completed successfully!

Summary:
  • 3 branches processed
  • 2 branches had PR metadata updates
  • 3 branches pushed to remote
  • 1 conflict resolved interactively

All dependent branches now have your changes from feature-auth.
```

## Integration Points

### Status Command Integration
- Reuse dependency discovery logic
- Share tree building algorithms
- Use same metadata parsing functions

### Git Integration
- Use `go-git` library for rebase operations
- Handle Git repository state and validation
- Respect Git configuration and hooks

### GitHub Integration
- Coordinate with GitHub API for any PR updates needed
- Maintain consistency between local and remote state

## Technical Implementation Notes

### Key Operations
1. **Tree building** - Parse metadata and construct dependency tree
2. **DFS traversal** - Process branches in correct dependency order
3. **Interactive rebase** - Handle conflicts with user interaction
4. **Metadata updates** - Keep PR references accurate after rebases
5. **Safe pushing** - Use force-with-lease for all pushes

### Data Structures
```go
type CascadeSession struct {
    StartBranch     string
    DependencyTree  *BranchNode
    ProcessedCount  int
    ConflictCount   int
    UpdatedPRs      []int
}
```

### Core Git Operations
```bash
# For each branch in DFS order:
git checkout <branch>
git rebase <parent-branch>
# Handle conflicts interactively...
git push --force-with-lease origin <branch>
```

## Success Criteria
- ✅ Successfully cascade changes through entire dependency tree
- ✅ Handle conflicts with clear, interactive resolution
- ✅ Update PR metadata accurately when parent PRs change
- ✅ Push all successful branches to remote safely
- ✅ Simple, predictable user experience with no complex state management
- ✅ Clear progress indication and final summary