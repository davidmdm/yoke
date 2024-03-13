package argocd

import (
	_ "embed"
	"fmt"

	"github.com/davidmdm/yoke/pkg/helm"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//go:embed argo-cd-5.10.4.tgz
var archive []byte

func RenderChart(release, namespace string, values *Values) ([]*unstructured.Unstructured, error) {
	chart, err := helm.LoadChartFromZippedArchive(archive)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart from zipped archive: %w", err)
	}

	return chart.Render(release, namespace, values)
}
