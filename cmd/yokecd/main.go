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
	"github.com/davidmdm/yoke/internal"
	"github.com/davidmdm/yoke/internal/wasi"
)

func main() {
	if err := run(getConfig()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const pluginName = "yokecd"

type ArgoConfig struct {
	Path      string
	Namespace string
	RepoURL   string
	Revision  string
}

type Config struct {
	Argo   ArgoConfig
	Flight string
}

func getConfig() (cfg Config) {
	conf.Var(conf.Environ, &cfg.Argo.Path, "ARGOCD_APP_SOURCE_PATH", conf.Default("."))
	conf.Var(conf.Environ, &cfg.Argo.RepoURL, "ARGOCD_APP_SOURCE_REPO_URL")
	conf.Var(conf.Environ, &cfg.Argo.Revision, "ARGOCD_APP_SOURCE_TARGET_REVISION", conf.Default("main"))
	conf.Var(conf.Environ, &cfg.Argo.Namespace, "ARGOCD_APP_NAMESPACE")
	conf.Var(conf.Environ, &cfg.Flight, "ARGOCD_ENV_FLIGHT")
	conf.Environ.MustParse()
	return
}

func debug(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func run(cfg Config) error {
	enc := json.NewEncoder(os.Stdout)

	if cfg.Flight == "" {
		debug("application is not a flight")
		return HandleAppSource(enc, cfg.Argo)
	}

	debug("application is a flight")

	var flight Flight
	if err := yaml.Unmarshal([]byte(cfg.Flight), &flight); err != nil {
		return fmt.Errorf("failed to parse flight: %w", err)
	}

	debug("flight: %+v\n", flight)

	debug("downloading wasm: %s", flight.Spec.WasmURL)

	resp, err := http.Get(flight.Spec.WasmURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code when fetching %s: %d", flight.Spec.WasmURL, resp.StatusCode)
	}

	wasm, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	debug("executing wasm")

	data, err := wasi.Execute(
		context.Background(),
		wasm,
		flight.Metadata.Name,
		strings.NewReader(flight.Spec.Input),
		flight.Spec.Args...,
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

func HandleAppSource(enc *json.Encoder, argo ArgoConfig) error {
	manifests, err := findManifests(".")
	if err != nil {
		return fmt.Errorf("failed to find manifests: %w", err)
	}

	for _, manifest := range manifests {
		debug("handling manifest: %s", manifest)
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

	debug("%s: %s", manifest, resource.GetKind())

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

	if err := enc.Encode(flight.AsArgoApplication(manifest, argo)); err != nil {
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
