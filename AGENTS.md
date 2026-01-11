# AGENTS.md

## Purpose
Gogh is a GitHub repository manager (ghq-like) with automation features (overlay/script/hook/extra). The CLI is the primary entry point.

## Architecture (enforced by arch-go)
- core/: domain entities and interfaces, no dependencies on other layers
- app/: use cases and business logic, depends only on core
- infra/: external integrations (filesystem, git, GitHub REST/GraphQL)
- ui/: CLI adapter using Cobra

## Entry Point / DI
- cmd/gogh/main.go builds stores/services and passes them to ui/cli.
- ui/cli/app.go wires all commands and persists stores in PersistentPostRunE.

## Config & Persistence
- Stores are TOML-based (app/config/*) and implement core/store interfaces.
- Paths are XDG-based and can be overridden by GOGH_*_PATH env vars.
- config.LoadAlternative loads v4 then falls back to v0 stores.
- "migrate" forces re-save to the current v4 format.

## Key Domain Concepts
- Reference = host/owner/name with optional alias (core/repository/parser.go).
- Layout = root/host/owner/name (infra/filesystem/layout_service.go).
- Finder scans roots to resolve refs/paths (infra/filesystem/finder_service.go).
- Multi-root is first-class (workspace roots + primary root).

## Automation System
- Overlay: file content stored separately; apply writes file to repo path.
- Script: Lua code stored separately; run uses gopher-lua.
- Hook: triggers on post-clone/post-fork/post-create; runs overlay or script.
- Extra: bundle of overlays+hooks for reuse.
- Hook invocation: app/hook/invoke/usecase.go.

## Script Execution Pipeline
- script invoke spawns "gogh script run" as a subprocess and sends a gob payload.
- Hidden command: ui/cli/commands/script_run.go; skips execution when GOGH_TEST_MODE=1.
- Lua globals table is named "gogh" (see lua/gogh.lua).

## Git/GitHub Integration
- Git ops are via core/git interface + infra/git (go-git).
- GitHub hosting uses REST + GraphQL (infra/github + infra/githubv4).
- Device auth is implemented in infra/github/authenticate_service.go.

## Dev Commands (from CLAUDE.md)
- make gen: generate GraphQL client + mocks
- make lint: golangci-lint + arch-go
- make test: go test with race
- make man: docs
- make install: build/install
- gofmt: use gofumpt

## Testing Conventions (from CLAUDE.md)
- Use standard library testing (no testify).
- Public funcs: *_test.go in separate package.
- Private funcs: *_private_test.go in same package.
- Always run tests with -p 1.
