package k8

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/davidmdm/x/xerr"

	"github.com/davidmdm/halloumi/internal"
)

type Client struct {
	dynamic   *dynamic.DynamicClient
	discovery *discovery.DiscoveryClient
	clientset *kubernetes.Clientset
}

func NewClient(cfg *rest.Config) (*Client, error) {
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client component: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8 clientset: %w", err)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client component: %w", err)
	}

	return &Client{
		dynamic:   dynamicClient,
		discovery: discoveryClient,
		clientset: clientset,
	}, nil
}

func (client Client) ApplyResources(ctx context.Context, resources []*unstructured.Unstructured) error {
	var errs []error
	for _, resource := range resources {
		if err := client.ApplyResource(ctx, resource, ApplyOpts{DryRun: true}); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", internal.Canonical(resource), err))
		}
	}

	if err := xerr.MultiErrOrderedFrom("dry run", errs...); err != nil {
		return err
	}

	for _, resource := range resources {
		if err := client.ApplyResource(ctx, resource, ApplyOpts{DryRun: false}); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", internal.Canonical(resource), err))
		}
	}

	return xerr.MultiErrOrderedFrom("", errs...)
}

type ApplyOpts struct {
	DryRun bool
}

func (client Client) ApplyResource(ctx context.Context, resource *unstructured.Unstructured, opts ApplyOpts) error {
	resourceInterface, err := client.getDynamicResourceInterface(resource)
	if err != nil {
		return fmt.Errorf("failed to resolve resource: %w", err)
	}

	_, err = resourceInterface.Apply(
		ctx,
		resource.GetName(),
		resource,
		metav1.ApplyOptions{
			FieldManager: "halloumi",
			DryRun: func() []string {
				if opts.DryRun {
					return []string{metav1.DryRunAll}
				}
				return nil
			}(),
		},
	)
	return err
}

func (client Client) MakeRevision(ctx context.Context, release string, resources []*unstructured.Unstructured) error {
	configmaps := client.clientset.
		CoreV1().
		ConfigMaps("kube-system")

	name := "halloumi-" + release

	configMap, err := configmaps.Get(ctx, name, metav1.GetOptions{})
	if kerrors.IsNotFound(err) {
		var revisions internal.Revisions
		revisions.Add(resources)

		data, err := json.Marshal(revisions)
		if err != nil {
			return err
		}

		config := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Data:       map[string]string{"revisions": string(data)},
		}

		_, err = configmaps.Create(ctx, config, metav1.CreateOptions{FieldManager: "halloumi"})
		return err
	}
	if err != nil {
		return fmt.Errorf("failed to lookup revision for %s: %w", release, err)
	}

	var revisions internal.Revisions
	if err := json.Unmarshal([]byte(configMap.Data["revisions"]), &revisions); err != nil {
		return fmt.Errorf("failed to parse revision history: %w", err)
	}

	revisions.Add(resources)

	data, err := json.Marshal(revisions)
	if err != nil {
		return err
	}

	configMap.Data["revisions"] = string(data)

	_, err = configmaps.Update(ctx, configMap, metav1.UpdateOptions{FieldManager: "halloumi"})
	return err
}

func (client Client) RemoveOrphans(ctx context.Context, previous, current []*unstructured.Unstructured) error {
	set := make(map[string]struct{})
	for _, resource := range current {
		set[internal.Canonical(resource)] = struct{}{}
	}

	var errs []error
	for _, resource := range previous {
		if _, ok := set[internal.Canonical(resource)]; ok {
			continue
		}

		resourceInterface, err := client.getDynamicResourceInterface(resource)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to resolve resource %s: %w", internal.Canonical(resource), err))
			continue
		}

		if err := resourceInterface.Delete(ctx, resource.GetName(), metav1.DeleteOptions{}); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete %s: %w", internal.Canonical(resource), err))
			continue
		}
	}

	return xerr.MultiErrOrderedFrom("", errs...)
}

func (client Client) GetCurrentResources(ctx context.Context, release string) ([]*unstructured.Unstructured, error) {
	revisions, err := client.getRevisions(ctx, release)
	if err != nil {
		return nil, fmt.Errorf("failed to get revision history: %w", err)
	}
	return revisions.CurrentResources(), nil
}

func (client Client) getDynamicResourceInterface(resource *unstructured.Unstructured) (dynamic.ResourceInterface, error) {
	gvk := schema.FromAPIVersionAndKind(resource.GetAPIVersion(), resource.GetKind())

	resources, err := client.discovery.ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		return nil, fmt.Errorf("failed to discover resources for %s: %w", gvk.GroupVersion().String(), err)
	}

	resourceName := func() string {
		for _, api := range resources.APIResources {
			if api.Kind == gvk.Kind && !strings.Contains(api.Name, "/") {
				return api.Name
			}
		}
		return ""
	}()

	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: resourceName,
	}

	namespace := internal.Namespace(resource)

	return client.dynamic.Resource(gvr).Namespace(namespace), nil
}

func (client Client) getRevisions(ctx context.Context, release string) (*internal.Revisions, error) {
	name := "halloumi-" + release

	configMap, err := client.clientset.CoreV1().ConfigMaps("kube-system").Get(ctx, name, metav1.GetOptions{})
	if kerrors.IsNotFound(err) {
		return new(internal.Revisions), nil
	}
	if err != nil {
		return nil, err
	}

	var revisions internal.Revisions
	if err := json.Unmarshal([]byte(configMap.Data["revisions"]), &revisions); err != nil {
		return nil, err
	}

	return &revisions, nil
}
