# Claude Code Project Configuration

## Project Overview
`gh-stacked` is a CLI tool for managing stacked Pull Request workflows on GitHub. It aims to simplify:
- Creating dependent PRs that build on each other
- Cascading rebases when base branches change  
- Visualizing PR dependency trees
- Streamlined workflow for stacked development

## Competitive Analysis
Existing tools in this space include:
- **git-spr** (Go): Each commit becomes a PR, simple workflow
- **Graphite CLI** (TypeScript): Full-featured stacking with team collaboration
- **stack-pr** (Modular): Commit-to-PR mapping with dependency management
- **ghstack** (Python): Facebook's approach to stacked diffs

## Development Standards

### Test-Driven Development (TDD) Approach
**MANDATORY**: All code development follows strict TDD workflow:

1. **Interface-First Design**: Always define interfaces before implementation
2. **Test-First**: Write tests that validate the interface behavior before any implementation
3. **Red-Green-Refactor**: Verify tests fail → implement minimal code → refactor for quality
4. **Validation**: All changes must be verified by automated tests

#### TDD Workflow Steps:
```bash
# 1. Define interface in Go
type MyService interface {
    DoSomething(input string) (result string, error)
}

# 2. Write tests that validate the interface
func TestMyService_DoSomething(t *testing.T) {
    // Test cases that validate expected behavior
}

# 3. Verify tests fail (Red)
go test ./...

# 4. Implement minimal functionality (Green)
# 5. Refactor and improve (Refactor)
# 6. Verify all tests pass
```

### Testing Requirements

#### Test Types & Coverage
- **Unit Tests**: Test individual functions and methods in isolation
- **Integration Tests**: Test component interactions with minimal mocking
- **Black Box Tests**: End-to-end tests running actual commands on real GitHub repositories
- **Use minimal mocking**: Prefer real implementations and test data over mocks when possible

#### Test Standards
- Tests must be deterministic and repeatable
- Use descriptive test names that explain the scenario: `TestBranchCommand_WhenParentHasMetadata_SetsCorrectParentPR`
- Test both success and error cases
- Include edge cases and boundary conditions
- All tests must pass before any commit

### Go Code Standards

#### Style Guidelines
- Follow **Google Go Style Guide**: https://google.github.io/styleguide/go/
- Follow **Uber Go Style Guide**: https://github.com/uber-go/guide/blob/master/style.md
- Use `gofmt`, `goimports`, and `golint` consistently
- Write idiomatic Go code with proper error handling

#### Library Constraints
- **ONLY use libraries specified in design documents**
- **Pre-approved libraries**:
  - `github.com/cli/go-gh` (GitHub CLI integration)
  - `github.com/go-git/go-git/v5` (Git operations)
  - `github.com/spf13/cobra` (CLI framework, inherited from go-gh)
  - Standard library packages
- **NO random library imports** - discuss and document any new dependencies first

### Task Planning & Documentation

#### Before Any Implementation
1. **Create `tasks.md`**: Break down work into fine-grained, reviewable tasks
2. **Review tasks**: Get approval on approach before starting implementation
3. **Document decisions**: Write `.md` files for any architectural or design decisions

#### Task Breakdown Example
```markdown
# Feature: Branch Command Implementation

## Tasks
1. [ ] Define `BranchService` interface with `CreateBranch(name, from string) error`
2. [ ] Write unit tests for `BranchService` covering:
   - Valid branch creation
   - Invalid branch names
   - Already existing branches
   - Git repository errors
3. [ ] Implement `BranchService` struct with go-git integration
4. [ ] Write integration tests with real git repository
5. [ ] Create CLI command wrapper using Cobra
6. [ ] Add end-to-end black box tests
```

### Code Review & Quality

#### Before Any Commit
- All tests pass: `go test ./...`
- Code is properly formatted: `gofmt -s -w .`
- No linting errors: `golint ./...`
- Dependencies are tidy: `go mod tidy`

#### When Implementation is Harder Than Expected
- **STOP** and discuss approach
- Document the complexity and alternative solutions
- Update tasks.md with revised plan
- Get alignment before proceeding

### Git Workflow
- Use conventional commit messages: `feat:`, `fix:`, `docs:`, `refactor:`, etc.
- Create feature branches for new work
- All changes must pass automated tests
- Never commit secrets, API keys, or sensitive data
- Each commit should represent a complete, tested unit of work

## Claude Code Workflows

### When Adding Features (TDD Workflow)
1. **Plan**: Create `tasks.md` with fine-grained breakdown
2. **Design**: Define interfaces and data structures first
3. **Test**: Write failing tests that validate the interface
4. **Implement**: Write minimal code to make tests pass
5. **Refactor**: Improve code quality while keeping tests green
6. **Verify**: Run all tests, linting, and formatting checks
7. **Document**: Update relevant `.md` files with decisions made
8. **Commit**: Commit complete, tested functionality

### When Debugging
1. **Reproduce**: Create a failing test that demonstrates the bug
2. **Investigate**: Identify root cause through existing tests
3. **Fix**: Make minimal changes to fix the issue
4. **Test**: Ensure fix works and doesn't break existing functionality
5. **Prevent**: Add regression tests if not already covered
6. **Commit**: Commit fix with clear explanation

### When Refactoring
1. **Safety**: Ensure comprehensive test coverage exists first
2. **Small Steps**: Make incremental changes with tests passing
3. **Interface Preservation**: Keep public interfaces stable
4. **Documentation**: Update docs to reflect any architectural changes
5. **Validation**: Run full test suite after each change

## Project Commands

### Development Commands
- **Test**: `go test ./...` - Run all tests
- **Test with coverage**: `go test -cover ./...` - Run tests with coverage report
- **Build**: `go build ./...` - Build all packages
- **Format**: `gofmt -s -w .` - Format all Go files
- **Imports**: `goimports -w .` - Fix imports in all Go files  
- **Lint**: `golint ./...` - Lint all packages
- **Vet**: `go vet ./...` - Static analysis
- **Mod tidy**: `go mod tidy` - Clean up dependencies

### Quality Checks (run before commit)
```bash
go test ./...           # All tests pass
gofmt -s -w .          # Code formatted
goimports -w .         # Imports fixed
golint ./...           # No lint errors
go vet ./...           # No vet warnings
go mod tidy            # Dependencies clean
```

### GitHub CLI Extension Commands
- **Install locally**: `gh extension install .` - Install current version for testing
- **Build extension**: `go build -o gh-stacked` - Build the extension binary

## CLI Commands Design

### Core Commands:
```bash
# Create new branch (works for first branch or stacking on current)
gh stacked branch <branch-name>

# Smart push: opens PR first time, updates existing PR after
gh stacked push [--draft]

# Status of current stack and PRs with tree visualization
gh stacked status

# Cleanup merged/closed branches and abandon unwanted PRs
gh stacked cleanup [--merged | --closed | <branch-name>]

# Cascading rebase when base branches change
gh stacked rebase [--cascade]

# Merge ready PRs in dependency order
gh stacked merge [--auto-cascade]
```

### Key Features:
- **Smart PR linking**: Automatically adds dependency info to PR descriptions
- **Tree visualization**: ASCII tree showing PR relationships  
- **Cascading operations**: Rebase/merge propagates through dependent PRs
- **GitHub integration**: Uses GitHub API for PR management

## Architecture Decisions

### ✅ **GitHub CLI Extension Approach**  
**Decision**: Build as `gh stacked` extension using `github.com/cli/go-gh` library

**Rationale**:
- **Target audience alignment** - Developers using stacked PRs likely already use `gh` CLI
- **Authentication solved** - Inherits `gh auth` automatically, biggest friction eliminated  
- **Natural workflow** - `gh stacked push` feels more intuitive than `gh-stacked push`
- **GitHub-first design** - Tool is GitHub-specific anyway
- **Faster MVP** - Focus on core logic instead of auth plumbing
- **Built-in GitHub integration** - Repository detection, environment setup included

### Technical Architecture
- **Extension framework**: GitHub CLI extension using `go-gh` library
- **Language**: Go with Cobra CLI framework (inherited from `go-gh`) 
- **Git operations**: go-git library for local Git operations
- **GitHub API**: Integrated via `go-gh` REST client (wraps `google/go-github`)
- **Authentication**: Automatic via GitHub CLI's auth system
- **Dependency tracking**: PR-based references in commit messages with dual annotation
- **PR linking**: Automatic dependency information in PR descriptions

## Dependencies & Tools
(List key dependencies and their purposes as they're added)

---
*This file should be updated as the project evolves. Use `#` in Claude sessions to quickly add memories here.*