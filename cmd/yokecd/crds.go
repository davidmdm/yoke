package main

import (
	"cmp"
	"encoding/json"
	"path/filepath"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	FlightAPIVersion = "yoke.cloud/v1alpha1"
	FlightKind       = "Flight"
)

type Resource[T any] struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata,omitempty"`
	Spec       T        `json:"spec"`
}

type Metadata struct {
	Labels      map[string]any `json:"labels,omitempty"`
	Annotations map[string]any `json:"annotations,omitempty"`
	Name        string         `json:"name"`
	Namespace   string         `json:"namespace,omitempty"`
}

type ApplicationSource struct {
	RepoURL        string       `json:"repoURL"`
	Path           string       `json:"path"`
	TargetRevision string       `json:"targetRevision"`
	Plugin         SourcePlugin `json:"plugin"`
}

type PluginEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PluginParameter struct {
	Name   string `json:"name"`
	String string `json:"string,omitempty"`
}
type SourcePlugin struct {
	Name       string            `json:"name"`
	Env        []PluginEnv       `json:"env,omitempty"`
	Parameters []PluginParameter `json:"parameters,omitempty"`
}

type ApplicationSpec struct {
	Source      ApplicationSource `json:"source,omitempty"`
	Project     string            `json:"project"`
	Destination struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"destination"`
}

type Application Resource[ApplicationSpec]

type FlightSpec struct {
	ApplicationSpec `json:",omitempty"`
	WasmURL         string   `yaml:"wasmURL,omitempty" json:"wasmURL"`
	Args            []string `yaml:"args,omitempty" json:"args,omitempty"`
	Input           string   `yaml:"input,omitempty" json:"input,omitempty"`
}

type Flight Resource[FlightSpec]

func (flight Flight) AsArgoApplication(manifest string, argo ArgoConfig) Application {
	data, _ := json.MarshalIndent(flight, "", "  ")

	appSpec := flight.Spec.ApplicationSpec

	manifestDir, _ := filepath.Split(manifest)

	appSpec.Source.Path = filepath.Join(argo.Path, manifestDir)
	appSpec.Source.RepoURL = argo.RepoURL
	appSpec.Source.TargetRevision = argo.Revision

	appSpec.Destination.Namespace = cmp.Or(appSpec.Destination.Namespace, argo.Namespace)

	appSpec.Source.Plugin.Name = cmp.Or(appSpec.Source.Plugin.Name, argo.PluginName)

	appSpec.Source.Plugin.Env = append(appSpec.Source.Plugin.Env, PluginEnv{
		Name:  "FLIGHT",
		Value: string(data),
	})

	appSpec.Source.Plugin.Parameters = append(appSpec.Source.Plugin.Parameters, PluginParameter{
		Name:   "manifest-file",
		String: filepath.Join(argo.Path, manifest),
	})

	return Application{
		APIVersion: "argoproj.io/v1alpha1",
		Kind:       "Application",
		Metadata:   flight.Metadata,
		Spec:       appSpec,
	}
}

func AsFlight(resource *unstructured.Unstructured) (*Flight, error) {
	if resource.GetKind() != FlightKind || resource.GetAPIVersion() != FlightAPIVersion {
		return nil, nil
	}

	data, err := resource.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var flight Flight
	if err := json.Unmarshal(data, &flight); err != nil {
		return nil, err
	}

	return &flight, nil
}
