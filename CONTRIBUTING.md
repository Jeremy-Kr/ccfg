# Contributing to ccfg

Thank you for your interest in contributing to ccfg. This document covers the development setup, coding conventions, and PR process.

## Development Setup

### Prerequisites

- **Go 1.26+** -- [Install Go](https://go.dev/doc/install)
- **golangci-lint** (optional, for linting) -- [Install](https://golangci-lint.run/welcome/install/)

### Build and Run

```bash
# Build
go build -o ccfg ./cmd/ccfg

# Run
./ccfg

# Or use the Makefile
make build
make run
```

### Testing

```bash
# Run all tests
go test ./...

# Test a specific package
go test ./internal/scanner/...

# Run with verbose output
go test -v ./...
```

### Linting

```bash
# Run linter
golangci-lint run

# Or use the Makefile
make lint
```

## Project Structure

```
cmd/ccfg/          Entry point
internal/
  scanner/         Config file discovery (paths, scan logic)
  parser/          Config file parsing (JSON, JSONC, Markdown)
  model/           Shared type definitions
  tui/             Bubbletea model, view, components
docs/              PRD, technical design, roadmap
```

## Coding Conventions

- **Formatting:** Go standard style (`gofmt`). Run `gofmt` before committing.
- **Encapsulation:** Use `internal/` packages. Do not expose implementation details.
- **Error handling:** Wrap errors with context using `fmt.Errorf("context: %w", err)`.
- **Interfaces:** Define interfaces at the call site, not in the implementing package.
- **Tests:** Place test files (`*_test.go`) in the same package as the code being tested.
- **Naming:** Follow Go conventions -- `camelCase` for unexported, `PascalCase` for exported identifiers.
- **Dependencies:** Keep dependencies minimal. Only add a library when truly necessary.

## Commit Message Convention

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
type(scope): description
```

### Types

| Type | Purpose |
|---|---|
| `feat` | New feature |
| `fix` | Bug fix |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `docs` | Documentation only |
| `test` | Adding or updating tests |
| `chore` | Maintenance tasks (dependencies, CI, etc.) |
| `style` | Code style changes (formatting, no logic change) |
| `perf` | Performance improvement |
| `ci` | CI/CD changes |

### Examples

```
feat(scanner): add Linux path support
fix(parser): handle trailing commas in JSONC
docs(readme): add installation instructions
refactor(tui): extract tree rendering into separate component
```

## Pull Request Process

1. **Fork** the repository and create a feature branch from `main`.
2. **Make your changes** following the coding conventions above.
3. **Write tests** for new functionality.
4. **Run tests and lint** before submitting:
   ```bash
   go test ./...
   golangci-lint run
   ```
5. **Open a PR** against `main` with a clear title and description.
6. **Respond to review feedback** promptly.

### PR Guidelines

- Keep PRs focused -- one logical change per PR.
- Include a brief description of what the PR does and why.
- Reference any related issues (e.g., `Fixes #42`).
- Ensure all CI checks pass.

## Key Constraints

- **Read-only:** ccfg must never modify any configuration file. This is a core design principle.
- **macOS + Linux:** Support both platforms. Path resolution differs by OS.
- **Single binary:** No external runtime dependencies.

## Questions?

Open an issue on GitHub if you have questions or suggestions.
