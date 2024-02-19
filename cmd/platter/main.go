package main

import (
	"encoding/json"
	"fmt"
	"os"

	k8 "github.com/davidmdm/halloumi/pkg/utils/resource"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	json.NewEncoder(os.Stdout).Encode([]any{
		k8.Deployment{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
			Metadata: k8.Metadata{
				Name:      "sample-app-prod",
				Namespace: "default",
			},
			Spec: k8.DeploymentSpec{
				Replicas: 3,
				Selector: k8.Selector{
					MatchLabels: map[string]string{"app": "sample-app"},
				},
				Template: k8.PodTemplateSpec{
					Metadata: k8.TemplateMetadata{
						Labels: map[string]string{"app": "sample-app"},
					},
					Spec: k8.PodSpec{
						Containers: []k8.Container{
							{
								Name:    "web-app",
								Image:   "alpine:latest",
								Command: []string{"watch", "echo", "hello", "world"},
							},
						},
					},
				},
			},
		},
	})

	return nil
}
