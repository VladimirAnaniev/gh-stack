# Implementation Order

## Command Dependencies & Development Sequence

### 1. `gh stacked branch` - Core Foundation
**Dependencies:** None
**Purpose:** Creates branches with stack metadata
**Key Functions:**
- Branch creation with proper Git operations
- Stack metadata embedding in commits
- Parent relationship establishment

### 2. `gh stacked push` - PR Creation & GitHub Integration  
**Dependencies:** `branch` command
**Purpose:** Amends commits, creates PRs, establishes GitHub linkage
**Key Functions:**
- Commit amendment with stack metadata
- GitHub PR creation via API
- PR description formatting with stack information

### 3. `gh stacked status` - Visibility & Feedback
**Dependencies:** `branch` + `push` commands
**Purpose:** Discovery engine and tree visualization
**Key Functions:**
- Stack discovery from commit metadata
- GitHub PR status resolution
- Tree visualization and status display

### 4. `gh stacked cleanup` - Branch Cleanup & Maintenance
**Dependencies:** `branch` + `push` + `status` commands
**Purpose:** Remove merged/closed branches and abandon unwanted PRs
**Key Functions:**
- Auto-cleanup merged and closed PR branches
- Manual branch abandonment (close PR + delete branch)
- Safety checks to prevent deletion of branches with children

### 5. `gh stacked rebase` - Stack Maintenance
**Dependencies:** `branch` + `push` + `status` + `cleanup` commands  
**Purpose:** Complex cascading rebase logic
**Key Functions:**
- Dependency-aware rebasing
- Cascading updates through stack
- Conflict detection and handling

### 6. `gh stacked merge` - Completion Workflow
**Dependencies:** All previous commands
**Purpose:** Final orchestration of entire workflow
**Key Functions:**
- Dependency-order merge sequencing
- Post-merge stack cleanup
- Automated stack maintenance after merges

## Rationale

Each command builds on the capabilities of the previous ones:

- **Branch** establishes the foundation with metadata tracking
- **Push** creates the GitHub integration layer needed for status
- **Status** provides the visibility needed for safe cleanup and rebase operations
- **Cleanup** provides branch maintenance needed before complex rebase operations  
- **Rebase** provides the maintenance capabilities needed for intelligent merging
- **Merge** orchestrates the complete workflow using all previous capabilities

## Development Strategy

Start with the simplest command (`branch`) and progressively build complexity. Each command should be fully functional and tested before moving to the next, ensuring a solid foundation for dependent functionality.