package main

import (
	"fmt"
	"os"

	"github.com/tbeijen/workspace/internal/cmd"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "activate":
		err = cmd.Activate(os.Args[2:])
	case "set":
		err = cmd.Set(os.Args[2:])
	case "unset":
		err = cmd.Unset()
	case "ls":
		err = cmd.Ls()
	case "current":
		err = cmd.Current()
	case "version":
		fmt.Println(version)
	case "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "workspace: unknown command %q\n", os.Args[1])
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "workspace: %s\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `Usage: workspace <command> [args]

Commands:
  activate <shell>                        Emit shell integration (zsh)
  set <client> <env> <role> [cluster]     Emit export statements
  unset                                   Emit unset statements
  ls                                      List available workspaces
  current                                 Show active workspace
  version                                 Print version

Environment:
  WORKSPACES_ROOT    Workspace directory (default: ~/workspaces)
`)
}
