package cmd

import (
	"fmt"
	"os"
	"strings"
)

func Unset() error {
	vars := os.Getenv("__SWORKSPACE_VARS")
	if vars == "" {
		return nil
	}

	for _, v := range strings.Split(vars, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			fmt.Printf("unset %s\n", v)
		}
	}

	fmt.Println("unset __SWORKSPACE_VARS")
	fmt.Println("unset __SWORKSPACE_ACTIVE")

	return nil
}
