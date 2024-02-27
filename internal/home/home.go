package home

import (
	"os"
	"path/filepath"
)

var (
	Dir        string
	Kubeconfig string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	Dir = home
	Kubeconfig = filepath.Join(home, ".kube/config")
}
