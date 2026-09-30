# workspace

## What this is

A Go CLI that outputs shell commands for workspace environment switching.
Binary name: `workspace`. User-facing shell function: `sw`.

The binary never modifies the shell directly. It prints export/unset statements
to stdout; the shell function evals them. This is the same pattern used by
mise (`mise activate zsh`) and fnox (`fnox activate zsh`).

## Building

```sh
go build -o workspace .
```

## Testing

```sh
go test ./...
```

## Project layout

```
main.go              # CLI entry, os.Args dispatch to cmd/*
internal/
  cmd/
    activate.go      # `workspace activate zsh` — emit shell function + completion
    set.go           # `workspace set <client> <env> <role> [cluster]` — emit exports
    unset.go         # `workspace unset` — emit unsets based on __WS_VARS
    ls.go            # `workspace ls` — list workspaces from filesystem
    current.go       # `workspace current` — print __WS_ACTIVE
  config/
    config.go        # TOML parsing, workspace discovery, path resolution
  shell/
    zsh.go           # zsh shell function + completion templates
testdata/
  workspaces/        # test fixtures mirroring ~/workspaces structure
```

## Design principles

- The binary is stateless. All state lives in env vars (`__WS_VARS`, `__WS_ACTIVE`).
- Output is plain shell commands, one per line. No interactive prompts, no TUI.
- Config files are TOML with a single `[env]` section. Keep the format minimal.
- Workspace root is `$WORKSPACES_ROOT` or `~/workspaces`.
- No CLI framework (cobra, etc). Subcommand dispatch via os.Args is sufficient.

## Workspace directory structure

```
~/workspaces/
  <client>/
    <env>/
      workspace.toml             # base env vars (loaded for every role)
      workspace.<role>.toml      # role overlay (merged on top of base)
      .kube/<cluster>            # optional kubeconfig files
```

Roles are discovered from file basenames: `workspace.ad.toml` defines role `ad`.

## Code style

- Standard Go: `gofmt`, no linter beyond `go vet`.
- No comments unless the why is non-obvious.
- Minimal dependencies. Currently only `github.com/BurntSushi/toml`.
- Error messages go to stderr. Shell commands go to stdout. Never mix them.
- Exit 0 on success (even if output is empty). Exit 1 on errors.

## Config format

`workspace.toml` (base):
```toml
[env]
WS_CLIENT = "acme"
WS_ENV    = "staging"
```

`workspace.<role>.toml` (overlay):
```toml
[env]
WS_ROLE     = "ad"
AWS_PROFILE = "acme-staging-ad"
```

## Not yet implemented

- Hooks (`[hooks]` section in config): activate/deactivate commands.
- Bash shell support.
- `workspace completions zsh` as a standalone subcommand.
