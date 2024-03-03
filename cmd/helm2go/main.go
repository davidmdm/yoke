package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/davidmdm/ansi"
	"github.com/davidmdm/x/xcontext"
	"github.com/davidmdm/yoke/internal/home"
)

var (
	cache        = filepath.Join(home.Dir, ".cache/yoke")
	schemaGenDir = filepath.Join(cache, "readme-generator-for-helm")
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	packageName := flag.String("package", "values", "package name for generate types")
	values := flag.String("values", "", "path to values file")
	out := flag.String("out", "", "path to outputfile for generated go types")

	flag.Parse()

	if *values == "" {
		return fmt.Errorf("-values is required")
	}
	if *out == "" {
		return fmt.Errorf("-out is required")
	}

	ctx, cancel := xcontext.WithSignalCancelation(context.Background(), syscall.SIGINT)
	defer cancel()

	if err := ensureReadmeGenerator(ctx); err != nil {
		return fmt.Errorf("failed to ensure bitnami/readme-generator installation: %w", err)
	}

	if err := ensureGoJsonSchema(ctx); err != nil {
		return fmt.Errorf("failed to ensure go-jsonschema installation: %w", err)
	}

	*values, _ = filepath.Abs(*values)
	*out, _ = filepath.Abs(*out)

	schemaFile := filepath.Join(os.TempDir(), "values")

	genSchema := exec.CommandContext(ctx, "node", "./bin/index.js", "-v", *values, "-s", schemaFile)

	if err := x(genSchema, WithDir(schemaGenDir)); err != nil {
		return fmt.Errorf("failed to generate jsonschema: %w", err)
	}

	genGoTypes := exec.CommandContext(ctx, "go-jsonschema", schemaFile, "-o", *out, "-p", *packageName, "--only-models")
	if err := x(genGoTypes); err != nil {
		return fmt.Errorf("failed to gen go types: %w", err)
	}

	return nil
}

func ensureReadmeGenerator(ctx context.Context) error {
	if err := os.MkdirAll(cache, 0o755); err != nil {
		return fmt.Errorf("failed to ensure yoke cache: %w", err)
	}

	if _, err := os.Stat(schemaGenDir); err != nil {
		clone := exec.CommandContext(ctx, "git", "clone", "https://github.com/bitnami/readme-generator-for-helm")
		if err := x(clone, WithDir(cache)); err != nil {
			return fmt.Errorf("failed to clone schema generator: %w", err)
		}

		downloadDeps := exec.CommandContext(ctx, "npm", "install")
		if err := x(downloadDeps, WithDir(schemaGenDir)); err != nil {
			return fmt.Errorf("failed to download schema generator dependencies: %w", err)
		}
	}

	return nil
}

func ensureGoJsonSchema(ctx context.Context) error {
	if x(exec.CommandContext(ctx, "command", "-v", "go-jsonschema")) == nil {
		return nil
	}
	if err := x(exec.CommandContext(ctx, "go", "install", "github.com/atombender/go-jsonschema@latest")); err != nil {
		return fmt.Errorf("failed to install go-jsonschema: %w", err)
	}

	return nil
}

var cyan = ansi.MakeStyle(ansi.FgCyan).Sprint

func x(cmd *exec.Cmd, opts ...XOpt) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	for _, apply := range opts {
		apply(cmd)
	}

	fmt.Println()
	fmt.Println("running:", cyan(strings.Join(cmd.Args, " ")))
	fmt.Println()

	return cmd.Run()
}

type XOpt func(*exec.Cmd)

func WithDir(dir string) XOpt {
	return func(c *exec.Cmd) {
		c.Dir = dir
	}
}
