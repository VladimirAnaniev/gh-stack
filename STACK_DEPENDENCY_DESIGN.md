# Stacked PR Dependency Tracking Design

## Overview
Design for tracking branch dependencies in `gh stacked` extension using PR-based references embedded in commit messages.

## Core Approach: PR-Based References

### Metadata in Commit Messages
Embed stack metadata directly in Git commit messages:
```bash
feat: add user authentication

gh-stacked: pr=124 parent-pr=123
```

### Annotation Strategy
- Users commit freely with normal messages
- When running `gh stacked push`, automatically amend all commits in branch with stack metadata
- Force push changes (automated, transparent to user)
- Create/update GitHub PR

## Why PR IDs Over Branch Names

### Problems with Branch-Based References
1. **Branch deletion after merge** - GitHub auto-deletes merged branches
2. **Branch name collisions** - Same name reused for different features
3. **Branch renaming** - User renames branch, breaks child references
4. **Multiple stacks** - Same branch names across different stacks

### Benefits of PR-Based References
- ✅ **Stable identifiers** - PR numbers never change or disappear
- ✅ **GitHub integration** - Easy to query PR status via API
- ✅ **Merge detection** - Can detect when parent PR was merged via UI
- ✅ **Globally unique** - No collisions across stacks

## Stack Metadata Format

```bash
# First branch in dependency tree (parent is main/master)
feat: implement user login
gh-stacked: pr=123

# Subsequent branches
feat: add login tests  
gh-stacked: pr=124 parent-pr=123

feat: refactor auth flow
gh-stacked: pr=125 parent-pr=124
```

## Discovery Process

1. **Scan local commits** for `gh-stacked:` metadata in all branches
2. **Query GitHub API** for PR details using PR numbers
3. **Build dependency tree** from parent-pr relationships
4. **Handle merged PRs** by detecting merge status and rebasing children onto target branch
5. **Cross-reference** with current branch states and PR statuses

## Handling Edge Cases

### UI Merges
- Detect when `parent-pr=123` is merged via GitHub API
- Automatically suggest/perform rebase of child branches onto target (usually `main`)

### Manual Branches  
- Branches created outside `gh stacked` have no metadata
- Ignored in stack operations until explicitly incorporated
- Future enhancement: allow manual addition to existing stacks

### Rebases/Squashes
- Branch-based references survive rebases (PR number unchanged)
- Stack metadata preserved through Git history rewriting
- Child branches maintain correct parent relationships

## User Workflow

```bash
# Create dependency tree
gh stacked branch feature-auth      # Creates branch, no parent
git commit -m "feat: add login"     # Normal commit
gh stacked push                     # Amends commit with metadata, creates PR #123

# Continue dependency tree  
gh stacked branch feature-auth-tests # Creates branch from feature-auth
git commit -m "test: add auth tests" # Normal commit
gh stacked push                      # Amends commit, creates PR #124

# Status shows dependency tree using PR relationships
gh stacked status --tree
```

## Benefits

- **Transparent to users** - Normal Git workflow preserved
- **Collaboration friendly** - Works across devices and team members
- **Resilient** - Survives rebases, force pushes, UI operations  
- **No external files** - All metadata in Git history
- **GitHub native** - Leverages PR system as source of truth