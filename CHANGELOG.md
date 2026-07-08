# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0]

### Added

- First public release
- Authentication with GitHub (PAT and OAuth device flow)
- Secure credential storage (OS keychain with file fallback)
- GitHub repository, issue, and pull request workflows
- Git repository context detection
- Shell completion (bash, zsh, fish)
- Configuration management
- Structured exit codes for scripting
- `kit doctor` diagnostics command
- Generated command documentation
- Live integration tests (opt-in)
- GitHub Actions CI and release automation
- Cross-platform binaries (macOS, Linux, Windows)

### Security

- Credentials stored in the operating system's secure credential store when available
- Fallback to file-based storage only when necessary
