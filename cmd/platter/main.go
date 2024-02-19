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
	name := "sample-app"
	labels := map[string]string{"app": name}

	deployment := k8.Deployment{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: k8.Metadata{
			Name:      name,
			Namespace: "default",
		},
		Spec: k8.DeploymentSpec{
			Replicas: 3,
			Selector: k8.Selector{MatchLabels: labels},
			Template: k8.PodTemplateSpec{
				Metadata: k8.TemplateMetadata{Labels: labels},
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
	}

	service := k8.Service{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata:   k8.Metadata{Name: name},
		Spec: k8.ServiceSpec{
			Selector: labels,
			Ports: []k8.ServicePort{
				{
					Protocol:   "TCP",
					Port:       80,
					TargetPort: 3000,
				},
			},
		},
	}

	return json.
		NewEncoder(os.Stdout).
		Encode([]any{deployment, service})
}
