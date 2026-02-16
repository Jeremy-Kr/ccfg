# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Translated all Korean comments, strings, and test messages to English
- Translated all documentation to English for open-source release

### Fixed

- Homebrew tap skip_upload handling when token is not available

## [0.1.0] - 2026-02-16

### Added

- TUI dashboard with tree view, preview panel, and search
- Config file scanner for Claude Code settings, instructions, MCP, and policy files
- JSON/JSONC parser with syntax highlighting (Glamour + Chroma)
- Merged config view with search functionality
- Agent character card and skill ability card UI
- Usage ranking view with tool/agent/skill statistics
- Virtual node tree for hooks and MCP server entries
- File watcher with fsnotify for auto-refresh
- Scrollbar support for tree, preview, and ranking panels
- `--version` flag with ldflags injection
- MIT license
- GoReleaser config with Homebrew tap support
- GitHub Actions release automation workflow
- golangci-lint, goreleaser, and Makefile configuration

### Fixed

- Directory preview errors and symlink handling
- Subdirectory recursive scan and tree depth limit removal
- Panel height causing left tree panel misalignment
- Line wrapping within panels (MaxWidth fix)
- Ranking view scroll and cursor visibility issues
- Scrollbar rendering and ranking width calculation
- Tree panel width calculation
- Project agent path detection

[Unreleased]: https://github.com/jeremy-kr/ccfg/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/jeremy-kr/ccfg/releases/tag/v0.1.0
