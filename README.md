# sworkspace

Shell workspace environment switcher. Sets and unsets environment variables
per shell session, independent of the current directory.

## Why

Tools like [direnv](https://direnv.net/) and
[mise](https://mise.jdx.dev/) tie environment to the current directory.
That works for project-level config but not for cross-cutting concerns like
AWS profiles, Kubernetes contexts, or team/environment identifiers that should
persist as you move between directories.

`sworkspace` fills that gap: activate a workspace once, work across directories,
deactivate when done. Multiple shells can have different workspaces active
independently.

## How it works

`sworkspace` is a Go binary that outputs shell commands. A thin shell function
(`sw`) evals the output, so exports and unsets happen in the calling shell.

```zsh
# In ~/.zshrc:
eval "$(sworkspace activate zsh)"

# Then:
sw acme staging ad           # activate workspace (client / env / role)
sw acme staging ad cluster1  # activate with kubeconfig
sw                           # deactivate
sw ls                        # list available workspaces
sw current                   # show active workspace
```

## Config

Workspaces live under `$SWORKSPACE_ROOT` (default `~/workspaces`), organized as
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
`__SWORKSPACE_VARS` so deactivate can cleanly unset them. Internal state uses
the `__SWORKSPACE_` prefix; user-configurable settings use `SWORKSPACE_`.

## Install

### Using mise

Install the latest eligible release globally:

```sh
mise use -g github:TBeijen/sworkspace
```

Mise defaults to a seven-day minimum release age as a supply-chain safeguard.
A newly published release may therefore not appear immediately. Only bypass
this protection when you trust the latest release:

```sh
MISE_MINIMUM_RELEASE_AGE=0s mise use -g github:TBeijen/sworkspace@latest
```

### From source

```sh
go install github.com/tbeijen/sworkspace@latest
```

### Build locally

```sh
git clone https://github.com/tbeijen/sworkspace.git
cd sworkspace
mkdir -p bin
go build -o bin/sworkspace .
cp bin/sworkspace ~/.local/bin/
```

## Shell support

Currently: zsh. Bash support can be added later.

## License

MIT
