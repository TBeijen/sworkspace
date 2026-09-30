package cmd

import (
	"fmt"
	"os"
)

func Current() error {
	active := os.Getenv("__WS_ACTIVE")
	if active == "" {
		fmt.Println("no workspace active")
		return nil
	}
	fmt.Println(active)
	return nil
}
