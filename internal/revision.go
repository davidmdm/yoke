package internal

import (
	"cmp"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Revisions struct {
	Total       int        `json:"total"`
	History     []Revision `json:"history"`
	ActiveIndex int        `json:"activeIndex"`
}

func (revisions *Revisions) Add(resources []*unstructured.Unstructured) {
	revisions.History = append(revisions.History, Revision{
		ID:        revisions.Total + 1,
		Resources: resources,
	})
	revisions.ActiveIndex = len(revisions.History) - 1
	revisions.Total++
}

func (revisions Revisions) CurrentResources() []*unstructured.Unstructured {
	if len(revisions.History) == 0 {
		return nil
	}
	return revisions.History[revisions.ActiveIndex].Resources
}

type Revision struct {
	ID        int                          `json:"id"`
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
