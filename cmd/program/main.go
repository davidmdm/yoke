package main

import (
	"encoding/json"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

	return encoder.Encode([]runtime.Object{
		&v1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "app",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:    "echo",
						Image:   "alpine:latest",
						Command: []string{"watch", "echo", "hello"},
					},
				},
			},
		},
	})
}
