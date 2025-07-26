# Future Features Plan

## Stack Restructuring

### `gh stacked restack <branch> --onto <new-parent>`
**Purpose**: Move a branch and its dependents to a new parent in the dependency tree

**Use Cases**:
- "I want to rebase feature-B from feature-A onto main"
- "Make this middle branch independent" 
- Restructure stacks as requirements evolve

**Examples**:
```bash
# Move feature-B from feature-A to main
gh stacked restack feature-B --onto main

# Move feature-C from feature-B to feature-A  
gh stacked restack feature-C --onto feature-A

# Make branch independent (special case)
gh stacked restack feature-B --onto main
```

**Behavior**:
1. Validate that target branch and new parent exist
2. Rebase target branch onto new parent
3. Update metadata (change parent-pr reference)
4. Cascade rebase all dependent branches
5. Update GitHub PR base branches

**Implementation Note**: This addresses both "changing stack structure" and "breaking dependencies" use cases identified in the system analysis. Planned for implementation after core MVP is stable.