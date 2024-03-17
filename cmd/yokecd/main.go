package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/davidmdm/conf"
	"github.com/davidmdm/yoke/internal/wasi"
)

func main() {
	if err := run(getConfig()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type ArgoConfig struct {
	Path     string
	RepoURL  string
	Revision string
}

type Config struct {
	Argo   ArgoConfig
	Flight string
}

func getConfig() (cfg Config) {
	conf.Var(conf.Environ, &cfg.Argo.Path, "ARGOCD_APP_SOURCE_PATH", conf.Default("."))
	conf.Var(conf.Environ, &cfg.Argo.RepoURL, "ARGGOCD_APP_SOURCE_REPO_URL")
	conf.Var(conf.Environ, &cfg.Argo.Revision, "ARGGOCD_APP_SOURCE_REVISION", conf.Default("main"))
	conf.Var(conf.Environ, &cfg.Flight, "ARGOCD_ENV_FLIGHT")
	conf.Environ.MustParse()
	return
}

func run(cfg Config) error {
	if cfg.Flight == "" {
		return HandleAppSource(cfg.Argo)
	}

	var flight Flight
	if err := yaml.Unmarshal([]byte(cfg.Flight), &flight); err != nil {
		return fmt.Errorf("failed to parse flight: %w", err)
	}

	resp, err := http.Get(flight.Spec.WasmURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	wasm, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)

	resources, err := wasi.Execute(
		context.Background(),
		wasm,
		flight.Metadata.Name,
		strings.NewReader(flight.Spec.Input),
		flight.Spec.Args...,
	)
	if err != nil {
		return fmt.Errorf("failed to execute flight wasm: %w", err)
	}

	for _, resource := range resources {
		if err := enc.Encode(resource); err != nil {
			return err
		}
	}

	return nil
}

func HandleAppSource(argo ArgoConfig) error {
	manifests, err := findManifests(argo.Path)
	if err != nil {
		return fmt.Errorf("failed to find manifests: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)

	for _, manifest := range manifests {
		if err := OutputManfiest(enc, argo, manifest); err != nil {
			fmt.Fprintf(os.Stderr, "failed to output manifest %s: %v\n\n", manifest, err)
		}
	}

	return nil
}

func OutputManfiest(enc *json.Encoder, argo ArgoConfig, manifest string) error {
	data, err := os.ReadFile(manifest)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var resource *unstructured.Unstructured
	if err := yaml.Unmarshal(data, &resource); err != nil {
		return fmt.Errorf("cannot unmarshall into unstructured resource: %w", err)
	}

	flight, err := AsFlight(resource)
	if err != nil {
		return fmt.Errorf("failed to convert resource to flight: %w", err)
	}

	if flight == nil {
		if err := enc.Encode(resource); err != nil {
			return fmt.Errorf("failed to encode as json: %w", err)
		}
		return nil
	}

	if err := enc.Encode(flight.AsArgoApplication(argo)); err != nil {
		return fmt.Errorf("failed to encode as json: %w", err)
	}
	return nil
}

func findManifests(root string) (manifests []string, err error) {
	err = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext != ".yml" && ext != ".yaml" {
			return nil
		}
		manifests = append(manifests, path)
		return nil
	})
	return
}
