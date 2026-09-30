package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func testdataRoot(t *testing.T) string {
	t.Helper()
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "..", "testdata", "workspaces")
}

func withRoot(t *testing.T) {
	t.Helper()
	t.Setenv("WORKSPACES_ROOT", testdataRoot(t))
}

func TestLoadEnv(t *testing.T) {
	withRoot(t)

	env, err := LoadEnv("acme", "staging", "ad")
	if err != nil {
		t.Fatalf("LoadEnv: %v", err)
	}

	want := map[string]string{
		"WS_CLIENT":   "acme",
		"WS_ENV":      "staging",
		"WS_ROLE":     "ad",
		"AWS_PROFILE": "acme-staging-ad",
	}

	for k, v := range want {
		if got := env[k]; got != v {
			t.Errorf("env[%q] = %q, want %q", k, got, v)
		}
	}

	if len(env) != len(want) {
		t.Errorf("got %d env vars, want %d", len(env), len(want))
	}
}

func TestLoadEnvOverlay(t *testing.T) {
	withRoot(t)

	env, err := LoadEnv("acme", "staging", "ro")
	if err != nil {
		t.Fatalf("LoadEnv: %v", err)
	}

	if got := env["AWS_PROFILE"]; got != "acme-staging-ro" {
		t.Errorf("AWS_PROFILE = %q, want %q", got, "acme-staging-ro")
	}
	if got := env["WS_ROLE"]; got != "ro" {
		t.Errorf("WS_ROLE = %q, want %q", got, "ro")
	}
}

func TestLoadEnvMissingBase(t *testing.T) {
	withRoot(t)

	_, err := LoadEnv("nonexistent", "env", "role")
	if err == nil {
		t.Fatal("expected error for missing workspace")
	}
}

func TestLoadEnvMissingRole(t *testing.T) {
	withRoot(t)

	_, err := LoadEnv("acme", "staging", "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing role")
	}
}

func TestKubeconfigPath(t *testing.T) {
	withRoot(t)

	p, err := KubeconfigPath("acme", "staging", "cluster1")
	if err != nil {
		t.Fatalf("KubeconfigPath: %v", err)
	}

	expected := filepath.Join(testdataRoot(t), "acme", "staging", ".kube", "cluster1")
	if p != expected {
		t.Errorf("got %q, want %q", p, expected)
	}
}

func TestKubeconfigPathMissing(t *testing.T) {
	withRoot(t)

	_, err := KubeconfigPath("acme", "staging", "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing cluster")
	}
}

func TestListWorkspaces(t *testing.T) {
	withRoot(t)

	workspaces, err := ListWorkspaces()
	if err != nil {
		t.Fatalf("ListWorkspaces: %v", err)
	}

	if len(workspaces) < 2 {
		t.Fatalf("got %d workspaces, want at least 2", len(workspaces))
	}

	found := false
	for _, ws := range workspaces {
		if ws.Client == "acme" && ws.Env == "staging" {
			found = true
			if len(ws.Roles) < 2 {
				t.Errorf("staging has %d roles, want at least 2", len(ws.Roles))
			}
		}
	}
	if !found {
		t.Error("acme/staging not found in workspace list")
	}
}

func TestListClients(t *testing.T) {
	withRoot(t)

	clients, err := ListClients()
	if err != nil {
		t.Fatalf("ListClients: %v", err)
	}

	if len(clients) == 0 {
		t.Fatal("no clients found")
	}
	if clients[0] != "acme" {
		t.Errorf("first client = %q, want %q", clients[0], "acme")
	}
}

func TestListEnvs(t *testing.T) {
	withRoot(t)

	envs, err := ListEnvs("acme")
	if err != nil {
		t.Fatalf("ListEnvs: %v", err)
	}

	if len(envs) < 2 {
		t.Fatalf("got %d envs, want at least 2", len(envs))
	}
}

func TestListRolesFor(t *testing.T) {
	withRoot(t)

	roles := ListRolesFor("acme", "staging")
	if len(roles) < 2 {
		t.Fatalf("got %d roles, want at least 2", len(roles))
	}

	hasAd := false
	hasRo := false
	for _, r := range roles {
		if r == "ad" {
			hasAd = true
		}
		if r == "ro" {
			hasRo = true
		}
	}
	if !hasAd || !hasRo {
		t.Errorf("expected roles ad and ro, got %v", roles)
	}
}

func TestListClusters(t *testing.T) {
	withRoot(t)

	clusters, err := ListClusters("acme", "staging")
	if err != nil {
		t.Fatalf("ListClusters: %v", err)
	}

	if len(clusters) != 1 || clusters[0] != "cluster1" {
		t.Errorf("got clusters %v, want [cluster1]", clusters)
	}
}

func TestRootDefault(t *testing.T) {
	os.Unsetenv("WORKSPACES_ROOT")
	root := Root()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, "workspaces")
	if root != expected {
		t.Errorf("Root() = %q, want %q", root, expected)
	}
}

func TestRootEnvOverride(t *testing.T) {
	t.Setenv("WORKSPACES_ROOT", "/custom/path")
	if got := Root(); got != "/custom/path" {
		t.Errorf("Root() = %q, want %q", got, "/custom/path")
	}
}
