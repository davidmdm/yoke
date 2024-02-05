package k8

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type Client struct {
	dynamic   *dynamic.DynamicClient
	discovery *discovery.DiscoveryClient
}

func NewClient(cfg *rest.Config) (*Client, error) {
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client component: %w", err)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client component: %w", err)
	}

	return &Client{dynamic: dynamicClient, discovery: discoveryClient}, nil
}

func (client Client) ApplyResource(ctx context.Context, resource *unstructured.Unstructured) error {
	gvk := schema.FromAPIVersionAndKind(resource.GetAPIVersion(), resource.GetKind())
	if gvk.Group == "" {
		gvk.Group = "core"
	}

	resources, err := client.discovery.ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		return fmt.Errorf("failed to discover resources for %s: %w", gvk.GroupVersion().String(), err)
	}

	gvr := schema.GroupVersionResource{
		Group:   gvk.Group,
		Version: gvk.Version,
		Resource: func() string {
			for _, api := range resources.APIResources {
				if api.Kind == gvk.Kind && !strings.Contains(api.Name, "/") {
					return api.Name
				}
			}
			return ""
		}(),
	}

	rc := client.dynamic.Resource(gvr)

	_, err = rc.Apply(ctx, resource.GetName(), resource, v1.ApplyOptions{FieldManager: "halloumi"})
	return err
}
