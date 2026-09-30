package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/tbeijen/workspace/internal/config"
)

func Ls() error {
	args := os.Args[2:]

	if len(args) >= 2 {
		switch args[0] {
		case "--clients":
			return lsClients()
		case "--envs":
			return lsEnvs(args[1])
		case "--roles":
			if len(args) >= 3 {
				return lsRoles(args[1], args[2])
			}
			return fmt.Errorf("usage: workspace ls --roles <client> <env>")
		case "--clusters":
			if len(args) >= 3 {
				return lsClusters(args[1], args[2])
			}
			return fmt.Errorf("usage: workspace ls --clusters <client> <env>")
		}
	}

	if len(args) == 1 && args[0] == "--clients" {
		return lsClients()
	}

	return lsAll()
}

func lsAll() error {
	workspaces, err := config.ListWorkspaces()
	if err != nil {
		return err
	}

	for _, ws := range workspaces {
		roles := strings.Join(ws.Roles, ", ")
		fmt.Printf("%s/%s  [%s]\n", ws.Client, ws.Env, roles)
	}

	return nil
}

func lsClients() error {
	clients, err := config.ListClients()
	if err != nil {
		return err
	}
	for _, c := range clients {
		fmt.Println(c)
	}
	return nil
}

func lsEnvs(client string) error {
	envs, err := config.ListEnvs(client)
	if err != nil {
		return err
	}
	for _, e := range envs {
		fmt.Println(e)
	}
	return nil
}

func lsRoles(client, env string) error {
	roles := config.ListRolesFor(client, env)
	for _, r := range roles {
		fmt.Println(r)
	}
	return nil
}

func lsClusters(client, env string) error {
	clusters, err := config.ListClusters(client, env)
	if err != nil {
		return err
	}
	for _, c := range clusters {
		fmt.Println(c)
	}
	return nil
}
