package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"slices"
)

type Config struct {
	KubeConfigPath string
	ReleaseName    string
	PlatterPath    string
	PlatterArgs    []string
}

var home string

func init() {
	home, _ = os.UserHomeDir()
}

func GetConfig() (Config, error) {
	var cfg Config

	flag.StringVar(&cfg.KubeConfigPath, "kubeconfig", filepath.Join(home, "/.kube/config"), "path to kube config")

	args := os.Args[1:]
	if idx := slices.Index(args, "--"); idx > -1 {
		if len(args) >= idx {
			cfg.PlatterArgs = args[idx+1:]
		}
		args = args[:idx]
	}

	if err := flag.CommandLine.Parse(args); err != nil {
		return cfg, err
	}

	cfg.ReleaseName = flag.Arg(0)
	cfg.PlatterPath = flag.Arg(1)

	if cfg.ReleaseName == "" || cfg.PlatterPath == "" {
		return cfg, errors.New("two positional arguments required:\n\n   usage: halloumi [release-name] [path]")
	}

	return cfg, nil
}
