package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/davidmdm/yoke/internal"
	"github.com/davidmdm/yoke/internal/wasi"
	"github.com/davidmdm/yoke/pkg/yoke"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cfg Config) error {
	out := json.NewEncoder(os.Stdout)

	debug("downloading wasm: %s", cfg.Flight.Wasm)

	wasm, err := yoke.LoadWasm(context.Background(), cfg.Flight.Wasm)
	if err != nil {
		return fmt.Errorf("failed to load wasm: %w", err)
	}

	debug("executing wasm")

	data, err := wasi.Execute(context.Background(), wasi.ExecParams{
		Wasm:    wasm,
		Release: cfg.Application.Name,
		Stdin:   strings.NewReader(cfg.Flight.Input),
		Args:    cfg.Flight.Args,
		Env:     map[string]string{"NAMESPACE": cfg.Application.Namespace},
	})
	if err != nil {
		return fmt.Errorf("failed to execute flight wasm: %w", err)
	}

	debug("wasm executed without error")

	var resources internal.List[*unstructured.Unstructured]
	if err := yaml.Unmarshal(data, &resources); err != nil {
		return fmt.Errorf("failed to unmarshal executed flight data: %w", err)
	}

	for _, resource := range resources {
		debug("encoding: %s/%s", resource.GetKind(), resource.GetName())
		if err := out.Encode(resource); err != nil {
			return err
		}
	}

	debug("resources: %d", len(resources))

	return nil
}

func debug(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
