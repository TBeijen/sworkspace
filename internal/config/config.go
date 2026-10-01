package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

type WorkspaceConfig struct {
	Env map[string]string `toml:"env"`
}

type Workspace struct {
	Client string
	Env    string
	Roles  []string
}

func Root() string {
	if r := os.Getenv("SWORKSPACE_ROOT"); r != "" {
		return r
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("~", "workspaces")
	}
	return filepath.Join(home, "workspaces")
}

func LoadEnv(client, env, role string) (map[string]string, error) {
	dir := filepath.Join(Root(), client, env)

	base, err := loadTOML(filepath.Join(dir, "workspace.toml"))
	if err != nil {
		return nil, fmt.Errorf("loading base config: %w", err)
	}

	overlay, err := loadTOML(filepath.Join(dir, fmt.Sprintf("workspace.%s.toml", role)))
	if err != nil {
		return nil, fmt.Errorf("loading role config: %w", err)
	}

	merged := make(map[string]string, len(base.Env)+len(overlay.Env))
	for k, v := range base.Env {
		merged[k] = v
	}
	for k, v := range overlay.Env {
		merged[k] = v
	}

	return merged, nil
}

func KubeconfigPath(client, env, cluster string) (string, error) {
	p := filepath.Join(Root(), client, env, ".kube", cluster)
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("kubeconfig not found: %s", p)
	}
	return p, nil
}

func ListWorkspaces() ([]Workspace, error) {
	root := Root()
	var workspaces []Workspace

	clients, err := listDirs(root)
	if err != nil {
		return nil, fmt.Errorf("reading workspace root %s: %w", root, err)
	}

	for _, client := range clients {
		envs, err := listDirs(filepath.Join(root, client))
		if err != nil {
			continue
		}
		for _, env := range envs {
			envDir := filepath.Join(root, client, env)
			if _, err := os.Stat(filepath.Join(envDir, "workspace.toml")); err != nil {
				continue
			}
			roles := listRoles(envDir)
			workspaces = append(workspaces, Workspace{
				Client: client,
				Env:    env,
				Roles:  roles,
			})
		}
	}

	return workspaces, nil
}

func ListClients() ([]string, error) {
	return listDirs(Root())
}

func ListEnvs(client string) ([]string, error) {
	var envs []string
	dirs, err := listDirs(filepath.Join(Root(), client))
	if err != nil {
		return nil, err
	}
	for _, d := range dirs {
		if _, err := os.Stat(filepath.Join(Root(), client, d, "workspace.toml")); err == nil {
			envs = append(envs, d)
		}
	}
	return envs, nil
}

func ListRolesFor(client, env string) []string {
	return listRoles(filepath.Join(Root(), client, env))
}

func ListClusters(client, env string) ([]string, error) {
	kubeDir := filepath.Join(Root(), client, env, ".kube")
	entries, err := os.ReadDir(kubeDir)
	if err != nil {
		return nil, nil
	}
	var clusters []string
	for _, e := range entries {
		if !e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			clusters = append(clusters, e.Name())
		}
	}
	sort.Strings(clusters)
	return clusters, nil
}

func loadTOML(path string) (*WorkspaceConfig, error) {
	var cfg WorkspaceConfig
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	if cfg.Env == nil {
		cfg.Env = make(map[string]string)
	}
	return &cfg, nil
}

func listDirs(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Strings(dirs)
	return dirs, nil
}

func listRoles(envDir string) []string {
	entries, err := os.ReadDir(envDir)
	if err != nil {
		return nil
	}
	var roles []string
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, "workspace.") && strings.HasSuffix(name, ".toml") && name != "workspace.toml" {
			role := strings.TrimPrefix(name, "workspace.")
			role = strings.TrimSuffix(role, ".toml")
			roles = append(roles, role)
		}
	}
	sort.Strings(roles)
	return roles
}
