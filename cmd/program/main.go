package main

import (
	"encoding/json"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	return encoder.Encode([]any{
		&v1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example-pod",
				Namespace: "default",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:    "pod-a",
						Image:   "my-image:tag",
						Command: []string{"do", "the", "thing"},
						Ports: []v1.ContainerPort{
							{
								HostPort:      3000,
								ContainerPort: 3000,
							},
						},
						Env: []v1.EnvVar{
							{
								Name:  "ENV",
								Value: "VAR",
							},
						},
					},
				},
			},
		},
	})
}
