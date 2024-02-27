package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var kubeconfig string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	kubeconfig = filepath.Join(home, ".kube/config")

	fmt.Println(kubeconfig)
}

func TestMain(m *testing.M) {
	must(x(exec.Command("kind", "delete", "cluster")))
	must(x(exec.Command("kind", "create", "cluster")))
	must(x(exec.Command("kubectl", "config", "use-context", "kind-kind")))

	fmt.Println()
	fmt.Println("cluster ready for tests")
	fmt.Println()

	os.Exit(m.Run())
}

func x(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
