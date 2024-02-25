package helm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"embed"
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
	"github.com/davidmdm/x/xerr"
)

func LoadChartFromZippedArchive(data []byte) (chart *Chart, err error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() {
		err = xerr.MultiErrFrom("", err, gz.Close())
	}()

	archive := tar.NewReader(gz)

	var files []*loader.BufferedFile
	for {
		header, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate through archive: %w", err)
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

	underlyingChart, err := loader.LoadFiles(files)
	if err != nil {
		return nil, err
	}

	return &Chart{underlyingChart}, nil
}

func LoadChartFromFS(fs embed.FS) (*Chart, error) {
	files, err := getAllFilesFromDir(fs, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to get files from FS: %w", err)
	}

	underlyingChart, err := loader.LoadFiles(files)
	if err != nil {
		return nil, err
	}

	return &Chart{underlyingChart}, nil
}

func getAllFilesFromDir(fs embed.FS, p string) ([]*loader.BufferedFile, error) {
	entries, err := fs.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir at %s: %w", p, err)
	}

	var results []*loader.BufferedFile
	for _, entry := range entries {
		filepath := path.Join(p, entry.Name())
		if entry.IsDir() {
			subEntries, err := getAllFilesFromDir(fs, filepath)
			if err != nil {
				return nil, err
			}
			results = append(results, subEntries...)
			continue
		}

		content, err := fs.ReadFile(filepath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file at %s: %w", filepath, err)
		}

		results = append(results, &loader.BufferedFile{
			Name: path.Join(p, entry.Name()),
			Data: content,
		})
	}

	return results, nil
}

type Chart struct {
	*chart.Chart
}

func (chart Chart) Render(release, namespace string, values any) ([]*unstructured.Unstructured, error) {
	opts := chartutil.ReleaseOptions{
		Name:      release,
		Namespace: namespace,
	}

	capabilities := chartutil.DefaultCapabilities.Copy()

	valueMap, err := asMap(values)
	if err != nil {
		return nil, fmt.Errorf("failed to convert values to map: %w", err)
	}

	valueMap, err = chartutil.ToRenderValues(chart.Chart, valueMap, opts, capabilities)
	if err != nil {
		return nil, err
	}

	rendered, err := engine.Engine{}.Render(chart.Chart, valueMap)
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
