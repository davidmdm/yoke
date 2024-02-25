package helm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"slices"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/davidmdm/halloumi/internal"
)

func loadTgz(data []byte) (*chart.Chart, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gzip: %w", err)
	}

	archive := tar.NewReader(gz)

	var files []*loader.BufferedFile
	for {
		header, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		content, err := io.ReadAll(archive)
		if err != nil {
			return nil, err
		}

		name := strings.Join(strings.Split(path.Clean(header.Name), "/")[1:], "/")

		files = append(files, &loader.BufferedFile{
			Name: name,
			Data: content,
		})
	}

	return loader.LoadFiles(files)
}

type Params struct {
	Source      []byte
	ReleaseName string
	Namespace   string
	Values      any
}

func Render(params Params) ([]*unstructured.Unstructured, error) {
	chart, err := loadTgz(params.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %w", err)
	}

	opts := chartutil.ReleaseOptions{
		Name:      params.ReleaseName,
		Namespace: params.Namespace,
	}

	capabilities := chartutil.DefaultCapabilities.Copy()

	values, err := asMap(params.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to convert values to map: %w", err)
	}

	values, err = chartutil.ToRenderValues(chart, values, opts, capabilities)
	if err != nil {
		return nil, err
	}

	rendered, err := engine.Engine{}.Render(chart, values)
	if err != nil {
		return nil, err
	}

	var results []*unstructured.Unstructured

	for name, content := range rendered {
		if ext := path.Ext(name); ext != ".yaml" {
			continue
		}

		var resource unstructured.Unstructured
		if err := yaml.Unmarshal([]byte(content), &resource); err != nil {
			return nil, fmt.Errorf("%s: %w\n%s", name, err, content)
		}
		if resource.Object == nil {
			continue
		}
		results = append(results, &resource)
	}

	slices.SortFunc(results, func(a, b *unstructured.Unstructured) int {
		return strings.Compare(internal.Canonical(a), internal.Canonical(b))
	})

	return results, nil
}

func asMap(values any) (map[string]any, error) {
	if m, ok := values.(map[string]any); ok {
		return m, nil
	}

	type Mappter interface {
		ToMap() (map[string]any, error)
	}
	if v, ok := values.(Mappter); ok {
		return v.ToMap()
	}

	data, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	err = json.Unmarshal(data, &m)
	return m, err
}
