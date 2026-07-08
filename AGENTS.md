# AGENTS.md

## Purpose

This repository contains `kit`, a lightweight Go command-line client for developer workflows.

## Principles

- Keep the CLI lightweight.
- Prefer the Go standard library unless a dependency provides clear value.
- Preserve backward compatibility for the public CLI.
- Favor simple, readable implementations over abstraction.

## Project Layout

- `cmd/kit`: CLI entry point.
- `internal/api`: HTTP client.
- `internal/auth`: authentication.
- `internal/commands`: CLI commands.
- `internal/config`: configuration.
- `internal/git`: Git utilities.
- `internal/exitcode`: exit codes.
- `tools/docsgen`: documentation generator.

## Guidelines

- Follow existing package structure.
- Keep commands focused on one responsibility.
- Keep changes focused on a single concern.
- Add or update tests with behavior changes.
- Run `go vet ./...` and `go test ./...` before submitting changes.
- Update generated documentation when commands change.

## Public Interface

The public CLI is stable. Avoid breaking commands, flags, or output formats except in a planned major release.

## Scope

Prefer improving existing functionality before introducing new abstractions. New dependencies should have a clear, long-term justification.
