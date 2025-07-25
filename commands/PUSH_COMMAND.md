# `gh stacked push` Command Specification

## Overview
Annotates all commits in the current branch with stack metadata, pushes to GitHub, and creates or updates a pull request with dependency information.

## Command Signature
```bash
gh stacked push [flags]
```

## Arguments
None - operates on current branch

## Flags
- `--draft`: Create PR as draft
- `--force`: Force push even if conflicts exist

## Behavior

### Core Functionality

#### 1. Parent Detection via Git History
- Use `git merge-base` to find where current branch diverged
- Analyze Git history to identify parent branch
- Check parent branch for existing gh-stacked metadata
- Determine parent PR number (if parent has PR)

#### 2. PR Creation/Update
- Create new PR or update existing PR for current branch
- Set PR base branch correctly based on dependency
- Get PR number for current branch

#### 3. Commit Annotation
- Amend ALL commits in current branch with complete stack metadata
- Add `gh-stacked: pr=<current-pr> parent-pr=<parent-pr>` to commit messages
- Include both own PR number and parent PR number for resilient tracking
- Preserve original commit messages and structure

#### 4. Force Push
- Push annotated commits to remote (force push required due to amended commits)
- Update PR with dependency information

## Examples

### Basic Usage
```bash
# First branch in tree (no parent)
git checkout feature-auth
git commit -m "feat: add login system"
gh stacked push
# → Creates PR #123, amends commit with "pr=123", pushes

# Dependent branch
git checkout feature-auth-tests  
git commit -m "test: add auth tests"
gh stacked push
# → Creates PR #124, amends with "pr=124 parent-pr=123", pushes
```

### Advanced Usage
```bash
# Create draft PR
gh stacked push --draft

# Force push (when rebasing has occurred)
gh stacked push --force
```

## Dual Annotation Strategy

### Complete Metadata Format
```bash
# Root branch commits (no parent)
feat: add user authentication

gh-stacked: pr=123

# Dependent branch commits
test: add auth unit tests

gh-stacked: pr=124 parent-pr=123
```

### Why Both PR Numbers Matter

**Own PR Number (`pr=124`):**
- **Commit ownership** - Clear which commits belong to which PR
- **Resilient tracking** - Works even after parent PRs are merged/deleted
- **Discovery efficiency** - Group commits by PR number

**Parent PR Number (`parent-pr=123`):**
- **Dependency tracking** - Build dependency trees
- **Rebase logic** - Know what to rebase when parent changes
- **Merge ordering** - Understand merge dependencies

## Operation Sequence

### 1. Parent Detection Logic
```bash
# Find divergence point
git merge-base current-branch main
# → Returns commit hash where branch diverged

# Determine parent branch
git branch --contains <divergence-commit>
# → Find which branch contains the divergence point

# Extract parent PR from parent branch commits
git log parent-branch --grep="gh-stacked:" -1
# → Parse "pr=123" to get parent PR number
```

### 2. PR Operations
- **Check existing PR**: Query GitHub API for PR on current branch
- **Create new PR**: If none exists, create with correct base branch
- **Update existing PR**: If exists, update as needed
- **Get PR number**: Extract PR number for annotation

### 3. Commit Annotation Process
```bash
# Get all commits in current branch (not in main)
git log main..HEAD --reverse --format="%H"

# For each commit:
# 1. Get original message: git log -1 --format="%B" <commit>
# 2. Append complete metadata
# 3. Amend commit: git commit --amend -m "<new-message>"
```

### 4. Force Push
- Push all amended commits to remote
- Force push is required due to commit message changes

## Parent Resolution Rules

1. **Branch from main/master**: No parent-pr (root of dependency tree)
   ```bash
   gh-stacked: pr=123
   ```

2. **Branch from branch with metadata**: Use that branch's PR as parent
   ```bash
   gh-stacked: pr=124 parent-pr=123
   ```

3. **Branch from branch without metadata**: Error - parent must be pushed first
   ```bash
   Error: Parent branch 'feature-auth' has no PR. Push parent branch first.
   ```

## GitHub Integration

### PR Configuration
```bash
# Root branch PR
base: main
head: feature-auth

# Dependent branch PR  
base: feature-auth  # PR base matches parent branch
head: feature-auth-tests
```

### PR Description Template
```markdown
**Depends on:** #123 (feature-auth)

---
*Created with gh-stacked*
```

## Validation Rules

### Pre-conditions
- ✅ Must be in a Git repository with GitHub remote
- ✅ Must be on a branch (not detached HEAD)
- ✅ Branch must have commits (not empty)
- ✅ Must have GitHub authentication configured
- ✅ Parent branch must have PR if this is a dependent branch

### Post-conditions
- ✅ All commits annotated with complete stack metadata (pr + parent-pr)
- ✅ Branch pushed to GitHub
- ✅ PR created or updated with correct base branch
- ✅ Dependencies clearly documented in PR

## Error Handling

### Dependency Validation
- **Parent has no PR**: Require parent to be pushed first
- **Parent PR is closed/merged**: Handle gracefully, suggest rebase
- **Circular dependencies**: Detect and prevent

### Git Operations
- **No commits**: Error message about empty branch
- **Merge conflicts**: Show conflict resolution instructions
- **Force push required**: Automatic (always needed due to amended commits)
- **Remote issues**: Network error handling and retry logic

### GitHub API
- **Authentication failed**: Clear instructions for gh auth setup
- **PR creation failed**: Show GitHub error with troubleshooting
- **Repository not found**: Verify remote configuration
- **Rate limiting**: Handle with exponential backoff

## Output

### Success Output
```bash
✓ Detected parent: feature-auth (PR #123)
✓ Created PR #124
✓ Annotated 3 commits with metadata (pr=124 parent-pr=123)
✓ Force pushed to origin/feature-auth-tests
ℹ PR URL: https://github.com/owner/repo/pull/124
ℹ Dependencies: #123 → #124
```

### Root Branch Output
```bash
✓ Starting new dependency tree from main
✓ Created PR #123
✓ Annotated 2 commits with metadata (pr=123)
✓ Pushed to origin/feature-auth
ℹ PR URL: https://github.com/owner/repo/pull/123
```

### Update Existing PR
```bash
✓ Found existing PR #124
✓ Updated 2 commits with latest metadata (pr=124 parent-pr=123)
✓ Force pushed to origin/feature-auth-tests
ℹ PR URL: https://github.com/owner/repo/pull/124
```

## Technical Implementation Notes

### Critical Operation Order
1. **Parent detection** - Must happen first
2. **PR creation/update** - Get PR number before annotation
3. **Commit annotation** - Add complete metadata to all commits
4. **Force push** - Required due to amended commits

### Metadata Parsing
```go
// Parse metadata from commit messages
type StackMetadata struct {
    PR       int  // Own PR number
    ParentPR *int // Parent PR number (nil for root)
}

// Example parsing
"gh-stacked: pr=124 parent-pr=123" → {PR: 124, ParentPR: &123}
"gh-stacked: pr=123"               → {PR: 123, ParentPR: nil}
```

## Success Criteria
- ✅ All commits annotated with complete metadata (own PR + parent PR)
- ✅ Robust dependency tracking that survives PR lifecycle changes
- ✅ Branch successfully pushed to GitHub with force push
- ✅ PR created/updated with correct base and dependency information
- ✅ Clear feedback about dependency relationships
- ✅ Minimal interface focused on core functionality