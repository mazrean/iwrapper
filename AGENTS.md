# Repository Guidelines — iwrapper

Interface wrapper code generator for Go.

> Agent configuration is managed via [apm](https://github.com/microsoft/apm).
> Common conventions live in `mazrean/apm-plackage/common`; Go-specific rules
> come from `mazrean/apm-plackage/go`. Run `apm install` to materialise locally.

## Build & Test

- `go test -v ./...`
- `go build ./...`
- `golangci-lint run`

## Conventions

- Specs go under `specs/`; use `mazrean/agent-skills/skills/writing-*`.
- Commit using Conventional Commits (`committing-code` skill).
- Use the Go 1.24+ `tool` directive for build tools; see `using-go-tool-directive` skill.
