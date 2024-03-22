package main

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	FlightAPIVersion = "yoke.cloud/v1alpha1"
	FlightKind       = "Flight"
)

type Resource[T any] struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       T        `json:"spec"`
}

type Metadata struct {
	Labels      map[string]any `json:"labels"`
	Annotations map[string]any `json:"annotations"`
	Name        string         `json:"name"`
	Namespace   string         `json:"namespace"`
}

type FlightSpec struct {
	WasmURL string   `json:"wasmUrl"`
	Args    []string `json:"args"`
	Input   string   `json:"input"`
}

type Flight Resource[FlightSpec]

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

type SourcePlugin struct {
	Name string      `json:"name"`
	Env  []PluginEnv `json:"env"`
}

type ApplicationSpec struct {
	Source ApplicationSource `json:"source"`
}

type App Resource[ApplicationSpec]

func (flight Flight) AsArgoApplication(argo ArgoConfig) App {
	data, _ := yaml.Marshal(flight)

	return App{
		APIVersion: "argoproj.io/v1alpha1",
		Kind:       "Application",
		Metadata:   flight.Metadata,
		Spec: ApplicationSpec{
			Source: ApplicationSource{
				RepoURL:        argo.RepoURL,
				Path:           argo.Path,
				TargetRevision: argo.Revision,
				Plugin: SourcePlugin{
					Name: argo.PluginName,
					Env: []PluginEnv{
						{Name: "flight", Value: string(data)},
					},
				},
			},
		},
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
