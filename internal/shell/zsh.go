package shell

import "fmt"

func ZshActivate() string {
	return fmt.Sprintf(`sw() {
  if [[ $# -eq 0 ]]; then
    eval "$(command workspace unset)"
    return $?
  fi
  case "$1" in
    ls|current|version|--help|-h)
      command workspace "$@"
      ;;
    *)
      eval "$(command workspace set "$@")"
      ;;
  esac
}

%s`, zshCompletion())
}

func zshCompletion() string {
	return `_sw() {
  local curcontext="$curcontext" state line
  local ws_root="${WORKSPACES_ROOT:-$HOME/workspaces}"

  _arguments -C \
    '1:client:->client' \
    '2:env:->env' \
    '3:role:->role' \
    '4:cluster:->cluster' && return

  case $state in
    client)
      local -a clients
      if [[ -d "$ws_root" ]]; then
        clients=("${(@f)$(command workspace ls --clients 2>/dev/null)}")
      fi
      _describe -t clients 'clients' clients && return
      ;;
    env)
      local -a envs
      envs=("${(@f)$(command workspace ls --envs "${words[2]}" 2>/dev/null)}")
      _describe -t envs 'environments' envs && return
      ;;
    role)
      local -a roles
      roles=("${(@f)$(command workspace ls --roles "${words[2]}" "${words[3]}" 2>/dev/null)}")
      _describe -t roles 'roles' roles && return
      ;;
    cluster)
      local -a clusters
      clusters=("${(@f)$(command workspace ls --clusters "${words[2]}" "${words[3]}" 2>/dev/null)}")
      _describe -t clusters 'clusters' clusters && return
      ;;
  esac
}

compdef _sw sw
`
}
