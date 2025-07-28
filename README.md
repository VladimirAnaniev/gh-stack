# gh-stack

A GitHub CLI extension for managing stacked pull requests.

## Overview

`gh-stack` simplifies working with stacked (dependent) pull requests by providing tools to:

- **Visualize PR dependencies** - See the dependency tree of your open PRs
- **Cascade rebases** - Automatically rebase dependent branches when base branches change
- **Streamlined workflow** - Manage complex PR stacks with simple commands

## Installation

Install as a GitHub CLI extension:

```bash
gh extension install VladimirAnaniev/gh-stack
```

Or build from source:

```bash
git clone https://github.com/VladimirAnaniev/gh-stack.git
cd gh-stack
go build -o gh-stack
gh extension install .
```

## Usage

### View Stack Status

See the dependency tree of your open PRs:

```bash
gh stack
```

This shows:
- Current branch (highlighted)
- PR status indicators (draft, approved, conflicts, etc.)
- Dependency relationships between PRs

### Cascade Rebase

When a base branch changes, cascade the rebase through all dependent branches:

```bash
gh stack cascade
```

This will:
1. Checkout and pull the main branch
2. For each PR targeting main: checkout, rebase, and push
3. For each dependent PR: checkout, rebase on its parent, and push
4. Handle merge conflicts with clear instructions

## How It Works

The tool builds a dependency tree by analyzing the base and head branches of your open PRs. It uses the GitHub CLI for authentication and API access, and go-git for local Git operations.

## Status Indicators

- 🔄 - Ready for review
- ✅ - Approved and ready to merge
- ❌ - Changes requested
- ⚠️ - Merge conflicts detected
- 📝 - Draft PR

## Requirements

- [GitHub CLI](https://cli.github.com/) installed and authenticated
- Go 1.19+ (for building from source)
- Git repository with GitHub remote

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## License

MIT License - see LICENSE file for details.