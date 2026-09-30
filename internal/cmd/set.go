package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tbeijen/workspace/internal/config"
)

func Set(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: workspace set <client> <env> <role> [cluster]")
	}

	client, env, role := args[0], args[1], args[2]

	envVars, err := config.LoadEnv(client, env, role)
	if err != nil {
		return err
	}

	var cluster string
	if len(args) >= 4 {
		cluster = args[3]
		kubePath, err := config.KubeconfigPath(client, env, cluster)
		if err != nil {
			return err
		}
		envVars["KUBECONFIG"] = kubePath
	}

	keys := sortedKeys(envVars)
	for _, k := range keys {
		fmt.Printf("export %s=%q\n", k, envVars[k])
	}

	fmt.Printf("export __WS_VARS=%q\n", strings.Join(keys, ","))
	fmt.Printf("export __WS_ACTIVE=%q\n", formatActive(client, env, role))

	return nil
}

func formatActive(client, env, role string) string {
	return fmt.Sprintf("%s/%s/%s", client, env, role)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
