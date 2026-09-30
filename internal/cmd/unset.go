package cmd

import (
	"fmt"
	"os"
	"strings"
)

func Unset() error {
	vars := os.Getenv("__WS_VARS")
	if vars == "" {
		return nil
	}

	for _, v := range strings.Split(vars, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			fmt.Printf("unset %s\n", v)
		}
	}

	fmt.Println("unset __WS_VARS")
	fmt.Println("unset __WS_ACTIVE")

	return nil
}
