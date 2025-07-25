# `gh stacked branch` Command Specification

## Overview
Creates a new Git branch for stacked development. Parent relationships are derived during push operations using Git history analysis.

## Command Signature
```bash
gh stacked branch <branch-name> [flags]
```

## Arguments
- `<branch-name>` (required): Name of the new branch to create

## Flags
- `--from <branch>`: Create branch from specific branch instead of current branch

## Behavior

### Core Functionality

#### 1. Branch Creation
- Create new Git branch from current HEAD (or `--from` branch)
- Switch to the newly created branch
- Ensure branch name doesn't conflict with existing branches

#### 2. Pure Git Operation
- No metadata storage or configuration files needed
- Parent relationships derived later during `push` using Git history analysis
- Clean branch creation with no extra commits or files

## Examples

### Basic Usage
```bash
# Create first branch in dependency tree (from main)
git checkout main
gh stacked branch feature-auth

# Create dependent branch (from feature-auth)  
git checkout feature-auth
gh stacked branch feature-auth-tests
```

### Advanced Usage
```bash
# Create branch from specific parent
gh stacked branch feature-cleanup --from feature-auth
```

## Validation Rules

### Pre-conditions
- ✅ Must be in a Git repository
- ✅ Must have GitHub remote configured
- ✅ Working directory must be clean (no uncommitted changes)
- ✅ Branch name must be valid Git branch name
- ✅ Branch name must not already exist
- ✅ If using `--from`, that branch must exist

### Post-conditions
- ✅ New branch created and checked out
- ✅ Ready for normal Git workflow (`git commit`, etc.)
- ✅ Parent relationship discoverable via Git history

## Error Handling

### Input Validation
- **Invalid branch name**: Show Git naming conventions
- **Branch already exists**: Suggest alternative names or checkout existing
- **Dirty working directory**: Prompt to commit or stash changes
- **No Git repository**: Clear error message with setup instructions
- **--from branch doesn't exist**: List available branches

### Git Operations
- **Branch creation fails**: Show underlying Git error
- **Remote issues**: Warn but don't block (can work offline)
- **Permission issues**: Clear error with troubleshooting steps

## Output

### Success Output
```bash
✓ Created branch 'feature-auth-tests'
ℹ Next: Make your changes and run 'gh stacked push' to create PR
```

### From main/master
```bash
✓ Created branch 'feature-auth'
ℹ Next: Make your changes and run 'gh stacked push' to create PR
```

### Verbose Output (`--verbose`)
```bash
✓ Current branch: feature-auth
✓ Creating branch 'feature-auth-tests' from 'feature-auth'
✓ Switched to branch 'feature-auth-tests'
ℹ Branch is ready for development
ℹ Parent relationship will be determined during push
```

## Parent Detection Strategy (for Push Command)

When `gh stacked push` runs, it will determine parent relationships using:

### Git History Analysis
1. **Find branch point**: Use `git merge-base` to find where current branch diverged
2. **Identify parent branch**: Determine which branch the current branch was created from
3. **Check for existing metadata**: Look for stacked-gh metadata in parent branch commits
4. **Build dependency chain**: Construct parent-pr relationships

### Detection Logic
```bash
# Example: On branch feature-auth-tests
git merge-base feature-auth-tests main
# → Returns commit hash where feature-auth-tests diverged

# Check if feature-auth has stacked metadata
git log feature-auth --grep="stacked-gh:" -1
# → If found, feature-auth is the parent
# → If not found, feature-auth might be root of tree
```

## Integration Points

### Git Integration
- Use `go-git` library for branch operations
- Respect Git configuration (user.name, user.email)
- Honor .gitignore and Git hooks
- No additional Git storage mechanisms needed

### GitHub Integration  
- Validate GitHub remote exists
- No immediate GitHub operations (defer to `push` command)

### Future Commands
- **Push command**: Will analyze Git history to determine parent relationships
- **Status command**: Will discover dependency trees from existing commit metadata
- **Other commands**: Will use Git history analysis for branch relationships

## Technical Implementation Notes

### Key Operations
1. Validate inputs and pre-conditions
2. Create Git branch using go-git (from current or --from branch)
3. Switch to new branch
4. Display success message with next steps
5. No additional storage or configuration needed

### Git History Integration
- Branch creation point is preserved in Git history
- Parent relationships discoverable via `git merge-base` and commit analysis
- Clean separation: branch creation is pure Git, metadata added during push

## Success Criteria
- ✅ New branch created successfully using standard Git operations
- ✅ No temporary files or configuration storage needed
- ✅ Parent relationships discoverable via Git history analysis
- ✅ User can continue with normal Git workflow
- ✅ Clean, minimal interface focused on branch creation only
- ✅ Proper error handling for all failure scenarios