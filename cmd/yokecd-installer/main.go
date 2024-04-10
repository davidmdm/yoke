package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// install.yaml downloaded from https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

//go:embed install.yaml
var install string

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var resources []*unstructured.Unstructured
	for _, manifest := range strings.Split(install, "\n---\n") {
		var resource unstructured.Unstructured
		if err := yaml.Unmarshal([]byte(manifest), &resource); err != nil {
			return err
		}
		resources = append(resources, &resource)
	}

	repoServer, i := func() (*unstructured.Unstructured, int) {
		for i, resource := range resources {
			if resource.GetName() == "argocd-repo-server" && resource.GetKind() == "Deployment" {
				return resource, i
			}
		}
		return nil, -1
	}()
	if i == -1 {
		return fmt.Errorf("cannot patch argocd: failed to find argocd-repo-server deployment")
	}

	var deployment appsv1.Deployment
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(repoServer.UnstructuredContent(), &deployment); err != nil {
		return fmt.Errorf("failed to convert argocd-repo-server to typed deployment: %w", err)
	}

	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, corev1.Container{
		Name:            "yokecd",
		Command:         []string{"/var/run/argocd/argocd-cmp-server"},
		Image:           "davidmdm/yokecd:latest",
		ImagePullPolicy: corev1.PullAlways,

		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "var-files",
				MountPath: "/var/run/argocd",
			},
			{
				Name:      "plugins",
				MountPath: "/home/argocd/cmp-server/plugins",
			},
			{
				Name:      "cmp-tmp",
				MountPath: "/tmp",
			},
		},

		SecurityContext: &corev1.SecurityContext{
			RunAsNonRoot: ptr(true),
			RunAsUser:    ptr[int64](999),
		},
	})

	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: "cmp-tmp",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	data, err := json.Marshal(deployment)
	if err != nil {
		return err
	}

	var resource *unstructured.Unstructured
	if err := json.Unmarshal(data, &resource); err != nil {
		return err
	}

	resources[i] = resource

	return json.NewEncoder(os.Stdout).Encode(resources)
}

func ptr[T any](value T) *T { return &value }
