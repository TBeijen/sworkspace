package shell

import (
	"strings"
	"testing"
)

func TestZshActivateContainsShellFunction(t *testing.T) {
	output := ZshActivate()

	mustContain := []string{
		"sw()",
		"command workspace set",
		"command workspace unset",
		"eval",
	}

	for _, s := range mustContain {
		if !strings.Contains(output, s) {
			t.Errorf("output missing %q", s)
		}
	}
}

func TestZshActivateContainsCompletion(t *testing.T) {
	output := ZshActivate()

	mustContain := []string{
		"_sw()",
		"compdef _sw sw",
		"->client",
		"->env",
		"->role",
		"->cluster",
	}

	for _, s := range mustContain {
		if !strings.Contains(output, s) {
			t.Errorf("completion output missing %q", s)
		}
	}
}
