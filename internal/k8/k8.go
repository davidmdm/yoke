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
	"k8s.io/client-go/util/retry"

	"github.com/davidmdm/x/xerr"

	"github.com/davidmdm/halloumi/internal"
)

const (
	ResourceReleaseMapping = "halloumi-resource-release-mapping"
	NSKubeSystem           = "kube-system"
	Halloumi               = "halloumi"
	KeyRevisions           = "revisions"
	KeyRelease             = "release"
)

func releaseName(release string) string { return Halloumi + "-" + release }

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
			FieldManager: Halloumi,
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

func (client Client) RemoveOrphans(ctx context.Context, previous, current []*unstructured.Unstructured) ([]*unstructured.Unstructured, error) {
	set := make(map[string]struct{})
	for _, resource := range current {
		set[internal.Canonical(resource)] = struct{}{}
	}

	var errs []error
	var removedResources []*unstructured.Unstructured
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

		removedResources = append(removedResources, resource)
	}

	return removedResources, xerr.MultiErrOrderedFrom("", errs...)
}

func (client Client) GetRevisions(ctx context.Context, release string) (*internal.Revisions, error) {
	name := releaseName(release)

	configMap, err := client.clientset.CoreV1().ConfigMaps(NSKubeSystem).Get(ctx, name, metav1.GetOptions{})
	if kerrors.IsNotFound(err) {
		return &internal.Revisions{Release: release}, nil
	}
	if err != nil {
		return nil, err
	}

	var revisions internal.Revisions
	if err := json.Unmarshal([]byte(configMap.Data[KeyRevisions]), &revisions); err != nil {
		return nil, err
	}

	return &revisions, nil
}

func (client Client) UpsertRevisions(ctx context.Context, release string, revisions *internal.Revisions) error {
	name := releaseName(release)

	configMaps := client.clientset.CoreV1().ConfigMaps(NSKubeSystem)

	data, err := json.Marshal(revisions)
	if err != nil {
		return err
	}

	configMap, err := configMaps.Get(ctx, name, metav1.GetOptions{})
	if kerrors.IsNotFound(err) {
		_, err := configMaps.Create(
			ctx,
			&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:   name,
					Labels: map[string]string{"internal.halloumi/kind": "revisions"},
				},
				Data: map[string]string{
					KeyRelease:   release,
					KeyRevisions: string(data),
				},
			},
			metav1.CreateOptions{FieldManager: Halloumi},
		)
		return err
	}

	if err != nil {
		return fmt.Errorf("failed to get revisions: %w", err)
	}

	configMap.Data[KeyRevisions] = string(data)

	_, err = configMaps.Update(ctx, configMap, metav1.UpdateOptions{FieldManager: "halloumu"})
	return err
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

func (client Client) UpdateResourceReleaseMapping(ctx context.Context, release string, create, remove []string) error {
	configMaps := client.clientset.CoreV1().ConfigMaps(NSKubeSystem)

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		configMap, err := configMaps.Get(ctx, ResourceReleaseMapping, metav1.GetOptions{})
		if kerrors.IsNotFound(err) {
			mapping := map[string]string{}
			for _, value := range create {
				mapping[value] = release
			}

			_, err := configMaps.Create(
				ctx,
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:   ResourceReleaseMapping,
						Labels: map[string]string{"internal.halloumi/kind": "resource-mapping"},
					},
					Data: mapping,
				},
				metav1.CreateOptions{FieldManager: Halloumi},
			)
			return err
		}

		if err != nil {
			return fmt.Errorf("failed to get resource to release mapping: %w", err)
		}

		for _, value := range remove {
			delete(configMap.Data, value)
		}
		for _, value := range create {
			configMap.Data[value] = release
		}

		_, err = configMaps.Update(ctx, configMap, metav1.UpdateOptions{FieldManager: Halloumi})
		return err
	})
}

func (client Client) GetResourceReleaseMapping(ctx context.Context) (map[string]string, error) {
	configMaps := client.clientset.CoreV1().ConfigMaps(NSKubeSystem)

	configMap, err := configMaps.Get(ctx, ResourceReleaseMapping, metav1.GetOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	return configMap.Data, nil
}

func (client Client) ValidateOwnership(ctx context.Context, release string, resources []*unstructured.Unstructured) error {
	resourceReleaseMapping, err := client.GetResourceReleaseMapping(ctx)
	if err != nil {
		return fmt.Errorf("failed to get release to resource mapping: %w", err)
	}

	var errs []error
	for _, resource := range internal.CanonicalNameList(resources) {
		if currentRelease, ok := resourceReleaseMapping[resource]; ok && currentRelease != release {
			errs = append(errs, fmt.Errorf("resource %+q is owned by release %+q", resource, currentRelease))
		}
	}

	return xerr.MultiErrOrderedFrom("conflict(s)", errs...)
}

func (client Client) GetAllRevisions(ctx context.Context) ([]internal.Revisions, error) {
	configMaps := client.clientset.CoreV1().ConfigMaps(NSKubeSystem)

	configs, err := configMaps.List(ctx, metav1.ListOptions{LabelSelector: "internal.halloumi/kind=revisions"})
	if err != nil {
		return nil, fmt.Errorf("failed to list revisions: %w", err)
	}

	results := make([]internal.Revisions, len(configs.Items))
	for i, cfg := range configs.Items {
		var revisions internal.Revisions
		if err := json.Unmarshal([]byte(cfg.Data[KeyRevisions]), &revisions); err != nil {
			return nil, fmt.Errorf("could not parse release %q state: %w", cfg.Data[KeyRelease], err)
		}
		results[i] = revisions
	}

	return results, nil
}
