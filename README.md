# workspace

Shell workspace environment switcher. Sets and unsets environment variables
per shell session, independent of the current directory.

## Why

Tools like [direnv](https://direnv.net/) and
[mise](https://mise.jdx.dev/) tie environment to the current directory.
That works for project-level config but not for cross-cutting concerns like
AWS profiles, Kubernetes contexts, or team/environment identifiers that should
persist as you move between directories.

`workspace` fills that gap: activate a workspace once, work across directories,
deactivate when done. Multiple shells can have different workspaces active
independently.

## How it works

`workspace` is a Go binary that outputs shell commands. A thin shell function
(`sw`) evals the output, so exports and unsets happen in the calling shell.

```zsh
# In ~/.zshrc:
eval "$(workspace activate zsh)"

# Then:
sw acme staging ad           # activate workspace (client / env / role)
sw acme staging ad cluster1  # activate with kubeconfig
sw                           # deactivate
sw ls                        # list available workspaces
sw current                   # show active workspace
```

## Config

Workspaces live under `$WORKSPACES_ROOT` (default `~/workspaces`), organized as
`<client>/<env>/`:

```
~/workspaces/
  acme/
    staging/
      workspace.toml           # base env vars
      workspace.ad.toml        # role overlay (ad = admin)
      workspace.ro.toml        # role overlay (ro = read-only)
      .kube/cluster1           # kubeconfig files
    production/
      workspace.toml
      workspace.ad.toml
      workspace.ro.toml
  bigcorp/
    dev/
      workspace.toml
      workspace.dev.toml
```

The `<role>` comes from the file basename suffix: `workspace.<role>.toml`.

`workspace.toml` (base, loaded for every role):

```toml
[env]
WS_CLIENT = "acme"
WS_ENV    = "staging"
```

`workspace.ad.toml` (role overlay, merged on top of base):

```toml
[env]
WS_ROLE     = "ad"
AWS_PROFILE = "acme-staging-ad"
```

On activate, base and role configs are merged. All exported vars are tracked in
`__WS_VARS` so deactivate can cleanly unset them.

## Install

### From GitHub releases (via mise)

```toml
[tools]
"github:tbeijen/workspace" = "0.1"
```

### From source

```sh
go install github.com/tbeijen/workspace@latest
```

### Build locally

```sh
go build -o workspace .
cp workspace ~/.local/bin/
```

## Shell support

Currently: zsh. Bash support can be added later.

## License

MIT
