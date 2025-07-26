# Core Metadata Management Service

## Purpose
Centralized service for managing `gh-stacked` metadata annotations in commit messages. Provides consistent metadata operations across all commands and ensures single source of truth for metadata handling.

## Responsibilities

### 1. Metadata Format Management
- Define and enforce canonical format: `gh-stacked: pr=X parent-pr=Y`
- Parse metadata from commit messages
- Serialize metadata back to commit messages
- Validate metadata structure and consistency

### 2. Core Operations
- **Annotate**: Add metadata to commits (used by PUSH)
- **Parse**: Extract metadata from commits (used by STATUS, CASCADE, CLEANUP, MERGE)
- **Update**: Modify existing metadata (used by CASCADE for orphaned branches)
- **Validate**: Check metadata integrity (used by STATUS as validation layer)

### 3. Shared Across Commands
- **PUSH**: Creates initial metadata annotations
- **STATUS**: Parses metadata to build dependency trees and validates consistency
- **CASCADE**: Updates metadata when rebasing orphaned branches (removes parent-pr)
- **CLEANUP**: Delegates metadata updates to CASCADE
- **MERGE**: Delegates metadata updates to CASCADE
- **BRANCH**: No metadata operations (pure Git branch creation)

## Key Benefits
- **Consistency**: All commands use same metadata format and parsing logic
- **Single Source of Truth**: One place to define metadata rules and operations
- **Validation**: Centralized validation ensures metadata integrity before operations
- **Maintainability**: Changes to metadata format only need to be made in one place

## Integration Pattern
Commands import and use the metadata service rather than implementing their own metadata logic:

```go
// Example usage pattern
metadataService.AnnotateCommits(branch, prNumber, parentPR)  // PUSH
tree := metadataService.BuildDependencyTree()               // STATUS  
metadataService.RemoveParentPR(branch)                      // CASCADE (for orphans)
validation := metadataService.ValidateMetadata(tree)        // STATUS (validation layer)
```

This ensures consistent behavior and eliminates duplication of metadata handling logic across commands.