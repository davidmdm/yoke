package main

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestTakeoff(t *testing.T) {
	settings := GlobalSettings{KubeConfigPath: kubeconfig}
	params := TakeoffParams{
		GlobalSettings: settings,
		Release:        "foo",
		Platter: TakeoffPlatterParams{
			Input: strings.NewReader(`
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    halloumi/managed-by: halloumi
  labels:
    app: sample-app
    halloumi/release: foo
  name: sample-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: sample-app
  template:
    metadata:
      labels:
        app: sample-app
    spec:
      containers:
        - command:
            - watch
            - echo
            - hello
            - world
          image: alpine:latest
          name: sample-app
`),
		},
		OutputDir: "",
	}

	restcfg, err := clientcmd.BuildConfigFromFlags("", params.KubeConfigPath)
	require.NoError(t, err)

	clientset, err := kubernetes.NewForConfig(restcfg)
	require.NoError(t, err)

	defaultDeployments := clientset.AppsV1().Deployments("default")

	deployments, err := defaultDeployments.List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)

	require.Len(t, deployments.Items, 0)

	require.NoError(t, TakeOff(context.Background(), params))

	deployments, err = defaultDeployments.List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)

	require.Len(t, deployments.Items, 1)

	require.NoError(t, Mayday(context.Background(), MaydayParams{
		GlobalSettings: settings,
		Release:        "foo",
	}))

	deployments, err = defaultDeployments.List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)

	require.Len(t, deployments.Items, 0)
}
