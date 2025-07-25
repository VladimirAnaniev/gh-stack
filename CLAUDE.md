# Claude Code Project Configuration

## Project Overview
`stacked-gh` is a CLI tool for managing stacked Pull Request workflows on GitHub. It aims to simplify:
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

### Code Style & Formatting
- Use consistent indentation (2 spaces for JS/TS, 4 for Python)
- Follow language-specific linting rules
- Always run linters and type checkers before committing
- Use meaningful variable and function names

### Git Workflow
- Use conventional commit messages: `feat:`, `fix:`, `docs:`, `refactor:`, etc.
- Create feature branches for new work
- Always test before committing
- Never commit secrets, API keys, or sensitive data

### Testing Philosophy
- Write tests for new features
- Run tests before committing changes
- Aim for good test coverage on critical paths
- Use descriptive test names

### File Organization
- Keep related files together
- Use clear directory structure
- Document architectural decisions in this file
- Update this file as the project evolves

## Claude Code Workflows

### When Adding Features
1. Understand requirements clearly
2. Plan implementation approach
3. Write code following project conventions
4. Add tests for new functionality
5. Run linters and tests
6. Commit with clear message

### When Debugging
1. Reproduce the issue
2. Identify root cause
3. Fix in small, testable increments
4. Verify fix works
5. Add tests to prevent regression

### When Refactoring
1. Understand existing code first
2. Make small, safe changes
3. Test after each change
4. Preserve existing functionality
5. Update documentation if needed

## Project Commands
(Update these as the project grows)
- Build: `npm run build` (or equivalent)
- Test: `npm test` (or equivalent)
- Lint: `npm run lint` (or equivalent)
- Type check: `npm run typecheck` (or equivalent)

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
- **Natural workflow** - `gh stacked push` feels more intuitive than `stacked-gh push`
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