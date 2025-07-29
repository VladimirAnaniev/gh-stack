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

Example output:
```
Stack Status (current: feature/auth-improvements)

🔄 feature/auth-improvements ← current #123 Add OAuth2 support
├── 📝 feature/user-profiles #124 User profile management
│   └── ✅ feature/dashboard #125 Enhanced dashboard UI
└── ⚠️ feature/notifications #126 Real-time notifications
```

This shows:
- Current branch (highlighted with "← current")
- PR status indicators (🔄 ready, 📝 draft, ✅ approved, ❌ changes requested, ⚠️ conflicts)
- Dependency relationships between PRs in a tree structure

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

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## Support

If you encounter any issues or have questions:
- Check existing [issues](https://github.com/VladimirAnaniev/gh-stack/issues)
- Create a new issue with detailed information
- Include your OS, Go version, and GitHub CLI version

## Changelog

See [releases](https://github.com/VladimirAnaniev/gh-stack/releases) for version history and changes.

## License

MIT License - see [LICENSE](LICENSE) file for details.