package main

import (
	"flag"
	"os"
	"path/filepath"
)

type GlobalSettings struct {
	KubeConfigPath string
}

var home string

func init() {
	home, _ = os.UserHomeDir()
}

func RegisterGlobalFlags(flagset *flag.FlagSet, settings *GlobalSettings) {
	flagset.StringVar(&settings.KubeConfigPath, "kubeconfig", filepath.Join(home, "/.kube/config"), "path to kube config")
}
