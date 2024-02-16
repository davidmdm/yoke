package internal

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Revision struct {
	Resources []*unstructured.Unstructured `json:"resources"`
}
