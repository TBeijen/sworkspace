package cmd

import (
	"fmt"

	"github.com/tbeijen/sworkspace/internal/shell"
)

func Activate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: sworkspace activate <shell>")
	}
	switch args[0] {
	case "zsh":
		fmt.Print(shell.ZshActivate())
		return nil
	default:
		return fmt.Errorf("unsupported shell: %s (supported: zsh)", args[0])
	}
}
