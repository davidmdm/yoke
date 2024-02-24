package main

import (
	"context"
	"crypto/sha1"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/davidmdm/x/xerr"

	"github.com/davidmdm/halloumi/internal"
	"github.com/davidmdm/halloumi/internal/k8"
	"github.com/davidmdm/halloumi/internal/wasi"
)

type TakeoffPlatterParams struct {
	Path  string
	Input io.Reader
	Args  []string
}

type TakeoffParams struct {
	GlobalSettings
	Release   string
	Platter   TakeoffPlatterParams
	OutputDir string
}

//go:embed cmd_takeoff_help.txt
var takeoffHelp string

func init() {
	takeoffHelp = strings.TrimSpace(internal.Colorize(takeoffHelp))
}

func GetTakeoffParams(settings GlobalSettings, source io.Reader, args []string) (*TakeoffParams, error) {
	flagset := flag.NewFlagSet("takeoff", flag.ExitOnError)

	flagset.Usage = func() {
		fmt.Fprintln(flagset.Output(), takeoffHelp)
		flagset.PrintDefaults()
	}

	params := TakeoffParams{
		GlobalSettings: settings,
		Platter:        TakeoffPlatterParams{Input: source},
	}

	RegisterGlobalFlags(flagset, &params.GlobalSettings)
	flagset.StringVar(&params.OutputDir, "outDir", "", "if present outputs platter resources to outDir instead of applying to k8")

	args, params.Platter.Args = internal.CutArgs(args)

	flagset.Parse(args)

	params.Release = flagset.Arg(0)
	params.Platter.Path = flagset.Arg(1)

	if params.Release == "" {
		return nil, fmt.Errorf("release is required as first positional arg")
	}
	if params.Platter.Input == nil && params.Platter.Path == "" {
		return nil, fmt.Errorf("platter-path is required as second position arg")
	}

	return &params, nil
}

func TakeOff(ctx context.Context, params TakeoffParams) error {
	output, wasm, err := func() ([]byte, []byte, error) {
		if params.Platter.Input != nil && params.Platter.Path == "" {
			output, err := io.ReadAll(params.Platter.Input)
			return output, nil, err
		}

		wasm, err := LoadWasm(ctx, params.Platter.Path)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read wasm program: %w", err)
		}

		output, err := wasi.Execute(ctx, wasm, params.Release, params.Platter.Input, params.Platter.Args...)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to execute wasm: %w", err)
		}

		return output, wasm, nil
	}()
	if err != nil {
		return fmt.Errorf("failed to load platter: %w", err)
	}

	var resources internal.List[*unstructured.Unstructured]
	if err := yaml.Unmarshal(output, &resources); err != nil {
		return fmt.Errorf("failed to unmarshal raw resources: %w", err)
	}

	internal.AddHallmouiMetadata(resources, params.Release)

	if params.OutputDir != "" {
		if err := ExportToFS(params.OutputDir, params.Release, resources); err != nil {
			return fmt.Errorf("failed to export release: %w", err)
		}
		return nil
	}

	restcfg, err := clientcmd.BuildConfigFromFlags("", params.KubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to build k8 config: %w", err)
	}

	client, err := k8.NewClient(restcfg)
	if err != nil {
		return fmt.Errorf("failed to instantiate k8 client: %w", err)
	}

	revisions, err := client.GetRevisions(ctx, params.Release)
	if err != nil {
		return fmt.Errorf("failed to get revision history: %w", err)
	}

	previous := revisions.CurrentResources()

	if reflect.DeepEqual(previous, []*unstructured.Unstructured(resources)) {
		return internal.Warning("resources are the same as previous revision: skipping takeoff")
	}

	if err := client.ValidateOwnership(ctx, params.Release, resources); err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}

	if err := client.ApplyResources(ctx, resources); err != nil {
		return fmt.Errorf("failed to apply resources: %w", err)
	}

	revisions.Add(resources, path.Clean(params.Platter.Path), fmt.Sprintf("%x", sha1.Sum(wasm)))

	if err := client.UpsertRevisions(ctx, params.Release, revisions); err != nil {
		return fmt.Errorf("failed to create revision: %w", err)
	}

	removed, err := client.RemoveOrphans(ctx, previous, resources)
	if err != nil {
		return fmt.Errorf("failed to remove orhpans: %w", err)
	}

	var (
		createdNames = internal.CanonicalNameList(resources)
		removedNames = internal.CanonicalNameList(removed)
	)

	if err := client.UpdateResourceReleaseMapping(ctx, params.Release, createdNames, removedNames); err != nil {
		return fmt.Errorf("failed to update resource release mapping: %w", err)
	}

	return nil
}

func LoadWasm(ctx context.Context, path string) (wasm []byte, err error) {
	uri, _ := url.Parse(path)
	if uri.Scheme == "" {
		return os.ReadFile(path)
	}

	if !slices.Contains([]string{":http", ":https"}, uri.Scheme) {
		return nil, errors.New("unsupported protocol: %s - http(s) supported only")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer func() {
		err = xerr.MultiErrFrom("", err, resp.Body.Close())
	}()

	return io.ReadAll(resp.Body)
}

func ExportToFS(dir, release string, resources []*unstructured.Unstructured) error {
	root := filepath.Join(dir, release)

	if err := os.RemoveAll(root); err != nil {
		return fmt.Errorf("failed remove previous platter export: %w", err)
	}

	if err := os.MkdirAll(root, 0o755); err != nil {
		return fmt.Errorf("failed to create release output directory: %w", err)
	}

	var errs []error
	for _, resource := range resources {
		path := filepath.Join(root, internal.Canonical(resource)+".yaml")
		if err := internal.WriteYAML(path, resource.Object); err != nil {
			errs = append(errs, err)
		}
	}

	return xerr.MultiErrFrom("", errs...)
}
