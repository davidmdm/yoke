package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/davidmdm/yoke/internal"
	"github.com/davidmdm/yoke/internal/wasi"
)

func main() {
	debug(os.Getenv("ARGOCD_APP_PARAMETERS"))

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
	enc := json.NewEncoder(os.Stdout)

	debug("downloading wasm: %s", cfg.Flight.WasmURL)

	resp, err := http.Get(cfg.Flight.WasmURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code when fetching %s: %d", cfg.Flight.WasmURL, resp.StatusCode)
	}

	wasm, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	debug("executing wasm")

	data, err := wasi.Execute(
		context.Background(),
		wasm,
		cfg.Application.Name,
		strings.NewReader(cfg.Flight.Input),
		cfg.Flight.Args...,
	)
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
		if err := enc.Encode(resource); err != nil {
			return err
		}
	}

	debug("resources: %d", len(resources))

	return nil
}

func debug(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
