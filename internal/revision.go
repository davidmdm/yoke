package internal

import (
	"cmp"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Revisions struct {
	Release     string     `json:"release"`
	Total       int        `json:"total"`
	History     []Revision `json:"history"`
	ActiveIndex int        `json:"activeIndex"`
}

func (revisions *Revisions) Add(resources []*unstructured.Unstructured, name, sha string) {
	revisions.History = append(revisions.History, Revision{
		ID:         revisions.Total + 1,
		Platter:    name,
		PlatterSHA: sha,
		CreatedAt:  time.Now(),
		Resources:  resources,
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
	ID         int                          `json:"id"`
	Platter    string                       `json:"platter"`
	PlatterSHA string                       `json:"platterSHA"`
	CreatedAt  time.Time                    `json:"createdAt"`
	Resources  []*unstructured.Unstructured `json:"resources"`
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
	gvk := resource.GetObjectKind().GroupVersionKind()

	return strings.ToLower(strings.Join(
		[]string{
			Namespace(resource),
			cmp.Or(gvk.Group, "core"),
			gvk.Version,
			resource.GetKind(),
			resource.GetName(),
		},
		".",
	))
}

func Namespace(resource *unstructured.Unstructured) string {
	return cmp.Or(resource.GetNamespace(), "default")
}

func CanonicalNameList(resources []*unstructured.Unstructured) []string {
	result := make([]string, len(resources))
	for i, resource := range resources {
		result[i] = Canonical(resource)
	}
	return result
}
