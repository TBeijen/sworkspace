package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func testdataRoot(t *testing.T) string {
	t.Helper()
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "..", "testdata", "workspaces")
}

func withRoot(t *testing.T) {
	t.Helper()
	t.Setenv("SWORKSPACE_ROOT", testdataRoot(t))
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestSetEmitsExports(t *testing.T) {
	withRoot(t)

	output := captureStdout(t, func() {
		err := Set([]string{"acme", "staging", "ad"})
		if err != nil {
			t.Fatalf("Set: %v", err)
		}
	})

	mustContain := []string{
		`export AWS_PROFILE="acme-staging-ad"`,
		`export WS_CLIENT="acme"`,
		`export WS_ENV="staging"`,
		`export WS_ROLE="ad"`,
		`export __SWORKSPACE_VARS=`,
		`export __SWORKSPACE_ACTIVE="acme/staging/ad"`,
	}

	for _, s := range mustContain {
		if !strings.Contains(output, s) {
			t.Errorf("output missing %q\nfull output:\n%s", s, output)
		}
	}
}

func TestSetWithCluster(t *testing.T) {
	withRoot(t)

	output := captureStdout(t, func() {
		err := Set([]string{"acme", "staging", "ad", "cluster1"})
		if err != nil {
			t.Fatalf("Set: %v", err)
		}
	})

	if !strings.Contains(output, `export KUBECONFIG=`) {
		t.Errorf("output missing KUBECONFIG\nfull output:\n%s", output)
	}
	if !strings.Contains(output, "cluster1") {
		t.Errorf("output missing cluster path\nfull output:\n%s", output)
	}
}

func TestSetMissingCluster(t *testing.T) {
	withRoot(t)

	err := Set([]string{"acme", "staging", "ad", "nonexistent"})
	if err == nil {
		t.Fatal("expected error for missing cluster")
	}
}

func TestSetTooFewArgs(t *testing.T) {
	err := Set([]string{"acme", "staging"})
	if err == nil {
		t.Fatal("expected error for too few args")
	}
}

func TestUnsetEmitsUnsets(t *testing.T) {
	t.Setenv("__SWORKSPACE_VARS", "WS_CLIENT,WS_ENV,WS_ROLE,AWS_PROFILE")

	output := captureStdout(t, func() {
		err := Unset()
		if err != nil {
			t.Fatalf("Unset: %v", err)
		}
	})

	mustContain := []string{
		"unset WS_CLIENT",
		"unset WS_ENV",
		"unset WS_ROLE",
		"unset AWS_PROFILE",
		"unset __SWORKSPACE_VARS",
		"unset __SWORKSPACE_ACTIVE",
	}

	for _, s := range mustContain {
		if !strings.Contains(output, s) {
			t.Errorf("output missing %q\nfull output:\n%s", s, output)
		}
	}
}

func TestUnsetNoVars(t *testing.T) {
	t.Setenv("__SWORKSPACE_VARS", "")

	output := captureStdout(t, func() {
		err := Unset()
		if err != nil {
			t.Fatalf("Unset: %v", err)
		}
	})

	if output != "" {
		t.Errorf("expected empty output, got: %s", output)
	}
}

func TestCurrentActive(t *testing.T) {
	t.Setenv("__SWORKSPACE_ACTIVE", "acme/staging/ad")

	output := captureStdout(t, func() {
		err := Current()
		if err != nil {
			t.Fatalf("Current: %v", err)
		}
	})

	if strings.TrimSpace(output) != "acme/staging/ad" {
		t.Errorf("got %q, want %q", strings.TrimSpace(output), "acme/staging/ad")
	}
}

func TestCurrentInactive(t *testing.T) {
	t.Setenv("__SWORKSPACE_ACTIVE", "")

	output := captureStdout(t, func() {
		err := Current()
		if err != nil {
			t.Fatalf("Current: %v", err)
		}
	})

	if strings.TrimSpace(output) != "no workspace active" {
		t.Errorf("got %q, want %q", strings.TrimSpace(output), "no workspace active")
	}
}

func TestActivateZsh(t *testing.T) {
	output := captureStdout(t, func() {
		err := Activate([]string{"zsh"})
		if err != nil {
			t.Fatalf("Activate: %v", err)
		}
	})

	mustContain := []string{
		"sw()",
		"command sworkspace",
		"compdef _sw sw",
		"_sw()",
	}

	for _, s := range mustContain {
		if !strings.Contains(output, s) {
			t.Errorf("output missing %q", s)
		}
	}
}

func TestActivateUnsupportedShell(t *testing.T) {
	err := Activate([]string{"fish"})
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
}

func TestActivateNoArgs(t *testing.T) {
	err := Activate([]string{})
	if err == nil {
		t.Fatal("expected error for missing shell arg")
	}
}

func TestLsAll(t *testing.T) {
	withRoot(t)
	os.Args = []string{"sworkspace", "ls"}

	output := captureStdout(t, func() {
		err := Ls()
		if err != nil {
			t.Fatalf("Ls: %v", err)
		}
	})

	if !strings.Contains(output, "acme/staging") {
		t.Errorf("output missing acme/staging\nfull output:\n%s", output)
	}
	if !strings.Contains(output, "acme/production") {
		t.Errorf("output missing acme/production\nfull output:\n%s", output)
	}
}
