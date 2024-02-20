package internal

import (
	"cmp"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Revision struct {
	Resources []*unstructured.Unstructured `json:"resources"`
}

func AddHallmouiMetadata(resources []*unstructured.Unstructured, release string) {
	for _, resource := range resources {
		annotations := resource.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}
		annotations["halloumi/managed-by"] = "halloumi"
		resource.SetAnnotations(annotations)

		labels := resource.GetLabels()
		if labels == nil {
			labels = make(map[string]string)
		}
		labels["halloumi/release"] = release

		resource.SetLabels(labels)
	}
}

func Canonical(resource *unstructured.Unstructured) string {
	return strings.ToLower(strings.Join(
		[]string{
			Namespace(resource),
			strings.ReplaceAll(resource.GetAPIVersion(), "/", "."),
			resource.GetKind(),
			resource.GetName(),
		},
		".",
	))
}

func Namespace(resource *unstructured.Unstructured) string {
	return cmp.Or(resource.GetNamespace(), "default")
}
