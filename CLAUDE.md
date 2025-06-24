# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gogh is a GitHub repository management tool written in Go, inspired by ghq. It manages GitHub repositories with features like cloning, forking, creating repositories, overlays, hooks, and scripts. A key feature is support for multiple root directories where repositories can be stored.

## Architecture

The project follows clean architecture with four distinct layers:

1. **Core Layer** (`core/`) - Domain entities and interfaces, no external dependencies
2. **Application Layer** (`app/`) - Use cases and business logic, only imports from core
3. **Infrastructure Layer** (`infra/`) - External integrations (filesystem, git, github)
4. **UI Layer** (`ui/`) - CLI implementation using Cobra framework

**Important**: Architecture boundaries are enforced by arch-go. Core cannot depend on other layers, app cannot depend on infra/ui.

## Development Commands

```bash
# Generate code (GraphQL client, mocks)
make gen

# Run linting (golangci-lint + arch-go)
make lint

# Run tests with race detection
make test

# Generate documentation
make man

# Build and install
make install

# Run a single test (always use -p 1 to limit parallel execution)
go test -p 1 -v -run TestName ./path/to/package

# Format Go files with gofumpt (handles formatting and trailing newlines)
go run mvdan.cc/gofumpt@latest -w path/to/file.go
```

## Testing Conventions

- Public function tests: `*_test.go` in separate test packages (import tested package as `testtarget`)
- Private function tests: `*_private_test.go` in same package
- Mock files: Place in separate `*_mock/` directories
- Use `go.uber.org/mock` for mock generation
- **DO NOT use testify package** - use standard library testing only
- **ALWAYS run tests with `-p 1` flag** to limit parallel execution (e.g., `go test -p 1 ./...`)

## Key Components

### Multi-Root Directory Management
- Gogh supports multiple root directories for storing repositories
- Primary root: Default location for new clones
- `roots` command: Manages root directories (add/remove/list/set-primary)
- Root directories are stored in `app/config/workspace.go`

### Repository Operations
- Clone/fork/create/delete: `app/action/` 
- Local repository management: `app/repo/`
- Path resolution: `core/repository/`

### Automation System (Hooks, Scripts, Overlays)

The automation system consists of three interconnected components:

1. **Scripts** - Lua scripts that can be executed in repository contexts
   - Uses `yuin/gopher-lua` and `vadv/gopher-lua-libs`
   - Scripts receive repository context via `gogh` global table
   - Management: `app/script/`, `core/script/`

2. **Overlays** - Predefined files that can be applied to repositories
   - Used for templates, config files, .gitignore, etc.
   - Must be manually applied with `gogh overlay apply` or automatically via hooks
   - Management: `app/overlay/`, `core/overlay/`

3. **Hooks** - Event-driven automation triggers
   - Trigger events: post-clone, post-fork, post-create
   - Can execute either overlays or scripts automatically
   - Pattern-based repository matching
   - Management: `app/hook/`, `core/hook/`

### GitHub Integration
- REST API: `infra/github/`
- GraphQL: `infra/githubv4/` (generated code)

## Important Files

- `cmd/gogh/main.go` - Entry point
- `ui/cli/app.go` - CLI application setup with all commands
- `app/config/config.go` - Configuration management
- `core/workspace/workspace.go` - Workspace interface for multi-root support

## Code Generation

GraphQL client code is generated from `infra/githubv4/*.graphql` files. After modifying these files, run `make gen`.

## Development Notes

- Comments should be in Japanese when beneficial for understanding
- Code and tests should be in English
- Follow existing patterns when adding new commands or features
- Multi-root directory support is a key feature - always consider this in path operations
- The `roots` command manages multiple repository storage locations, not the CLI root command
- Overlays are never applied automatically - they must be applied manually or through hooks
- Hook system can trigger overlay application and script execution based on repository events