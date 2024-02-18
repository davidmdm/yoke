package k8

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
			errs = append(errs, fmt.Errorf("%s: %w", canonical(resource), err))
		}
	}

	if err := xerr.MultiErrOrderedFrom("dry run", errs...); err != nil {
		return err
	}

	for _, resource := range resources {
		if err := client.ApplyResource(ctx, resource, ApplyOpts{DryRun: false}); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", canonical(resource), err))
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

	data, err := json.Marshal(internal.Revision{Resources: resources})
	if err != nil {
		return err
	}

	name := "halloumi-" + release

	revisions, err := configmaps.Get(ctx, name, metav1.GetOptions{})
	if kerrors.IsNotFound(err) {
		_, err := configmaps.Create(
			ctx,
			&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{Name: name},
				Data: map[string]string{
					"current": "1",
					"1":       string(data),
				},
			},
			metav1.CreateOptions{FieldManager: "halloumi"},
		)
		return err
	}
	if err != nil {
		return fmt.Errorf("failed to lookup revision for %s: %w", release, err)
	}

	var latest int
	for key := range revisions.Data {
		if version, _ := strconv.Atoi(key); version > latest {
			latest = version
		}
	}

	nextVersion := strconv.Itoa(latest + 1)
	revisions.Data[nextVersion] = string(data)
	revisions.Data["current"] = nextVersion

	_, err = configmaps.Update(ctx, revisions, metav1.UpdateOptions{FieldManager: "halloumi"})
	return err
}

func (client Client) RemoveOrphans(ctx context.Context, release string) error {
	name := "halloumi-" + release

	configmap, err := client.clientset.
		CoreV1().
		ConfigMaps("kube-system").
		Get(ctx, name, metav1.GetOptions{})

	if kerrors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	current, _ := strconv.Atoi(configmap.Data["current"])
	if current == 1 {
		return nil
	}

	var currentRevision internal.Revision
	if err := json.Unmarshal([]byte(configmap.Data[strconv.Itoa(current)]), &currentRevision); err != nil {
		return err
	}

	set := map[string]struct{}{}
	for _, resource := range currentRevision.Resources {
		set[canonical(resource)] = struct{}{}
	}

	var previousRevision internal.Revision
	if err := json.Unmarshal([]byte(configmap.Data[strconv.Itoa(current-1)]), &previousRevision); err != nil {
		return err
	}

	var errs []error
	for _, resource := range previousRevision.Resources {
		if _, ok := set[canonical(resource)]; ok {
			continue
		}

		resourceInterface, err := client.getDynamicResourceInterface(resource)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to resolve resource %s: %w", canonical(resource), err))
			continue
		}

		if err := resourceInterface.Delete(ctx, resource.GetName(), metav1.DeleteOptions{}); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete %s: %w", canonical(resource), err))
			continue
		}
	}

	return xerr.MultiErrOrderedFrom("", errs...)
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

	namespace := getNamespace(resource)

	return client.dynamic.Resource(gvr).Namespace(namespace), nil
}

func canonical(resource *unstructured.Unstructured) string {
	return getNamespace(resource) + "/" + resource.GetKind() + "/" + resource.GetName()
}

func getNamespace(resource *unstructured.Unstructured) string {
	return cmp.Or(resource.GetNamespace(), "default")
}
