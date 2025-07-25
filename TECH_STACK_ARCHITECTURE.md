# Technology Stack & Code Architecture

## Overview
`gh-stacked` is a GitHub CLI extension built in Go that provides stacked PR workflow management with beautiful terminal UI.

## Technology Stack

### Core Dependencies

```go
// CLI Framework
github.com/spf13/cobra              // Command structure and argument parsing

// GitHub Integration
github.com/cli/go-gh                // GitHub CLI extension framework and API client

// Git Operations
github.com/go-git/go-git/v5         // Pure Go Git implementation for local operations

// Terminal UI
github.com/charmbracelet/bubbletea  // TUI framework for interactive components
github.com/charmbracelet/bubbles    // Pre-built UI components (tables, trees, spinners)
github.com/charmbracelet/lipgloss   // Styling and layout for terminal output

```

### Why These Choices

**Cobra Framework:**
- ✅ Used by GitHub CLI itself - consistent patterns
- ✅ Rich flag/argument handling
- ✅ Automatic help generation
- ✅ Subcommand structure perfect for `gh stacked <command>`

**go-gh Library:**
- ✅ Official GitHub CLI extension framework
- ✅ Automatic authentication (inherits from `gh auth`)
- ✅ Repository detection and GitHub API client
- ✅ Consistent with gh conventions

**go-git Library:**
- ✅ Pure Go implementation - no external Git dependency
- ✅ Full Git operations (branches, commits, refs)
- ✅ Cross-platform compatibility
- ✅ Thread-safe for concurrent operations

**Charm Bracelet Suite:**
- ✅ Free and open source (MIT licensed)
- ✅ Production-ready (used by major tools)
- ✅ Rich terminal UI components
- ✅ Perfect for tree visualization and interactive displays

## Project Structure

```
gh-stacked/
├── gh-stacked                 # Main executable (required naming)
├── main.go                    # Entry point
├── go.mod                     # Go module definition
├── go.sum                     # Dependency lock file
├── 
├── cmd/                       # Cobra command definitions
│   ├── root.go               # Root command setup
│   ├── branch.go             # gh stacked branch
│   ├── push.go               # gh stacked push
│   ├── status.go             # gh stacked status
│   ├── cleanup.go            # gh stacked cleanup
│   ├── rebase.go             # gh stacked rebase
│   └── merge.go              # gh stacked merge
├── 
├── pkg/                       # Core business logic
│   ├── discovery/            # Stack discovery from commits/PRs
│   │   ├── commit_parser.go  # Parse stack metadata from commits
│   │   ├── pr_resolver.go    # Resolve PR information
│   │   └── stack_builder.go  # Build dependency trees
│   ├── 
│   ├── git/                  # Git operations wrapper
│   │   ├── repository.go     # Repository operations
│   │   ├── branch.go         # Branch management
│   │   └── commit.go         # Commit operations and metadata
│   ├── 
│   ├── github/               # GitHub API operations
│   │   ├── client.go         # API client wrapper
│   │   ├── pr.go             # Pull request operations
│   │   └── repository.go     # Repository information
│   ├── 
│   ├── stack/                # Stack management logic
│   │   ├── models.go         # Data structures
│   │   ├── operations.go     # Stack operations (rebase, merge)
│   │   └── validation.go     # Stack integrity checks
│   └── 
│   └── ui/                   # Terminal UI components
│       ├── tree.go           # Tree visualization using Bubbletea
│       ├── status.go         # Status display components
│       └── styles.go         # Lipgloss styling definitions
├── 
├── internal/                  # Private utilities
│   └── utils/                # Helper functions
└── 
└── docs/                     # Documentation
    ├── STACK_DEPENDENCY_DESIGN.md
    └── TECH_STACK_ARCHITECTURE.md
```

## Command Architecture

### Cobra Integration Pattern

```go
// main.go - Entry point
func main() {
    cmd.Execute()
}

// cmd/root.go - Root command setup
var rootCmd = &cobra.Command{
    Use:   "stacked",
    Short: "Manage stacked pull requests",
}

func Execute() {
    rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(branchCmd)
    rootCmd.AddCommand(pushCmd)
    rootCmd.AddCommand(statusCmd)
    rootCmd.AddCommand(cleanupCmd)
    rootCmd.AddCommand(rebaseCmd)
    rootCmd.AddCommand(mergeCmd)
}
```

### GitHub CLI Extension Requirements

- **Executable name:** Must be `gh-stacked` (matching repository name)
- **Location:** Executable at repository root
- **Installation:** `gh extension install vladimir-ananiev/gh-stacked`
- **Usage:** `gh stacked <command>` (GitHub CLI strips `gh-` prefix)

## Data Flow Architecture

### Stack Discovery Process

1. **Local Git Scan** (`pkg/git/`) - Scan commits for stack metadata
2. **GitHub API Query** (`pkg/github/`) - Resolve PR information  
3. **Stack Building** (`pkg/discovery/`) - Build dependency tree
4. **UI Rendering** (`pkg/ui/`) - Display results with Bubbletea

### Command Execution Flow

```
User Input → Cobra Command → Business Logic → Git/GitHub APIs → UI Output
    ↓            ↓              ↓               ↓              ↓
gh stacked → cmd/status.go → pkg/stack/ → pkg/git/ → pkg/ui/tree.go
  status                                    pkg/github/
```

## Extension Distribution

### Development
```bash
go mod init github.com/vladimir-ananiev/gh-stacked
go build -o gh-stacked
```

### Installation
```bash
gh extension install vladimir-ananiev/gh-stacked
```

### Usage
```bash
gh stacked status
gh stacked branch feature-name
gh stacked push --draft
gh stacked cleanup --merged
```

## Benefits of This Architecture

- **Modular Design** - Clear separation of concerns
- **GitHub Native** - Leverages official GitHub CLI patterns
- **Rich UI** - Beautiful terminal interface with Charm tools
- **Cross-Platform** - Pure Go, works everywhere GitHub CLI works
- **Extension Ecosystem** - Easy discovery and installation via gh
- **Developer Friendly** - Familiar patterns for Go/CLI developers