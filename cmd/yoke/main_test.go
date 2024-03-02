package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	appsv1config "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1config "k8s.io/client-go/applyconfigurations/core/v1"
	metav1config "k8s.io/client-go/applyconfigurations/meta/v1"

	"github.com/davidmdm/yoke/internal/home"
	"github.com/davidmdm/yoke/internal/k8s"
)

func createBasicDeployment(t *testing.T, name, namespace string) io.Reader {
	deployment := appsv1config.Deployment(name, namespace).
		WithLabels(map[string]string{"app": name}).
		WithSpec(
			appsv1config.DeploymentSpec().
				WithSelector(metav1config.LabelSelector().
					WithMatchLabels(map[string]string{"app": name}),
				).
				WithTemplate(
					corev1config.PodTemplateSpec().
						WithLabels(map[string]string{"app": name}).
						WithSpec(
							corev1config.PodSpec().WithContainers(
								corev1config.Container().
									WithName(name).
									WithImage("alpine:latest").
									WithCommand("watch", "echo", "hello", "world"),
							)),
				),
		)

	data, err := json.Marshal(deployment)
	require.NoError(t, err)

	return bytes.NewReader(data)
}

func TestCreateDeleteCycle(t *testing.T) {
	settings := GlobalSettings{KubeConfigPath: home.Kubeconfig}
	params := TakeoffParams{
		GlobalSettings: settings,
		Release:        "foo",
		Flight: TakeoffFlightParams{
			Input: createBasicDeployment(t, "sample-app", "default"),
		},
	}

	restcfg, err := clientcmd.BuildConfigFromFlags("", home.Kubeconfig)
	require.NoError(t, err)

	clientset, err := kubernetes.NewForConfig(restcfg)
	require.NoError(t, err)

	client, err := k8s.NewClient(restcfg)
	require.NoError(t, err)

	mappings, err := client.GetResourceReleaseMapping(context.Background())
	require.NoError(t, err)
	require.Equal(t, map[string]string{}, mappings)

	revisions, err := client.GetAllRevisions(context.Background())
	require.NoError(t, err)
	require.Len(t, revisions, 0)

	defaultDeployments := clientset.AppsV1().Deployments("default")

	deployments, err := defaultDeployments.List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)

	require.Len(t, deployments.Items, 0)

	require.NoError(t, TakeOff(context.Background(), params))

	mappings, err = client.GetResourceReleaseMapping(context.Background())
	require.NoError(t, err)
	require.Equal(t, map[string]string{"default.apps.v1.deployment.sample-app": "foo"}, mappings)

	deployments, err = defaultDeployments.List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)

	require.Len(t, deployments.Items, 1)

	require.NoError(t, Mayday(context.Background(), MaydayParams{
		GlobalSettings: settings,
		Release:        "foo",
	}))

	mappings, err = client.GetResourceReleaseMapping(context.Background())
	require.NoError(t, err)
	require.Equal(t, map[string]string{}, mappings)

	deployments, err = defaultDeployments.List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)

	require.Len(t, deployments.Items, 0)
}

func TestFailApplyDryRun(t *testing.T) {
	settings := GlobalSettings{KubeConfigPath: home.Kubeconfig}
	params := TakeoffParams{
		GlobalSettings: settings,
		Release:        "foo",
		Flight: TakeoffFlightParams{
			Input: createBasicDeployment(t, "sample-app", "does-not-exist"),
		},
	}

	require.EqualError(
		t,
		TakeOff(context.Background(), params),
		`failed to apply resources: dry run: does-not-exist.apps.v1.deployment.sample-app: namespaces "does-not-exist" not found`,
	)
}
