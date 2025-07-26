# UI Components - Charm Integration Patterns

## Component Types

### Spinner Reporter
**Use**: Single operations, unknown duration
**Visual**: Rotating spinner with status message
```
⠋ Creating PR for feature-branch...
⠙ Pushing feature-2 to origin...
⠹ Analyzing stack dependencies...
```

### Tree Progress Reporter  
**Use**: Hierarchical operations, step-by-step progress
**Visual**: ASCII tree with status icons
```
Creating stacked branch
├── ✓ Switch to parent branch (feature-1)
├── ⠋ Create branch feature-2...
├── ⏳ Set up branch metadata
└── ⏳ Update stack relationships
```

### Progress Bar Reporter
**Use**: Counted operations, known total steps
**Visual**: Progress bar with current step indicator
```
Cascading rebase [████████░░] 4/5 branches
⠋ Rebasing feature-3...
```

### Live Tree Updates
**Use**: Dynamic operations where items change state
**Visual**: Tree that updates in real-time
```
Merging ready PRs in order...

main
├── ⚡ feature-1 #123 (merging...)
└── ⏳ feature-2 #124 (waiting)
    └── ⏳ feature-3 #125 (waiting)
```

### Static Rich Display
**Use**: Status views, final results
**Visual**: Beautifully formatted static output
```
📚 Current Stack (3 branches)

main
└── feature-1 ✓ #123 Ready to merge
    └── feature-2 🔄 #124 CI running  
        └── feature-3 📝 #125 Draft
```

## Command → Component Mapping

| Command | Component | Reason |
|---------|-----------|---------|
| `branch <name>` | Tree Progress | Shows hierarchy being built |
| `push` | Spinner + List | Single operation with context |
| `status` | Static Rich Display | Beautiful formatted output |
| `rebase --cascade` | Progress Bar + Tree | Multiple items with progress |
| `merge --auto-cascade` | Live Tree Updates | Dynamic state changes |
| `cleanup` | Progress Bar | Known items to process |

## Charm Library Usage

### Bubble Tea
- Main framework for interactive components
- Handle real-time updates and state management
- Manage component lifecycle

### Lip Gloss  
- All styling and colors
- Tree formatting and alignment
- Status icons and visual hierarchy

### Spinner
- Loading states for operations
- Different spinner styles per operation type

### Progress
- Progress bars for counted operations
- Percentage and step indicators

## Visual Design Patterns

### Icons
- ✓ Success/Complete
- ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏ Spinner states  
- ⏳ Pending/Waiting
- ⚡ In Progress/Active
- ✗ Error/Failed
- 📚 Stack/Collection
- 🔄 Running/Processing
- 📝 Draft/Incomplete

### Colors
- Green: Success states
- Yellow: In progress, warnings
- Red: Errors, failures  
- Blue: Information, links
- Gray: Pending, disabled

### Tree Formatting
- ├── Standard tree branch
- └── Last item in tree
- │   Continuation line
- Proper spacing and alignment