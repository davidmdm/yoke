package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	y3 "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/davidmdm/x/xerr"

	"github.com/davidmdm/yoke/internal"
	"github.com/davidmdm/yoke/internal/wasi"
	"github.com/davidmdm/yoke/pkg/yoke"
)

type TakeoffPlatterParams struct {
	Path  string
	Input io.Reader
	Args  []string
}

type TakeoffParams struct {
	GlobalSettings
	Release string
	Platter TakeoffPlatterParams
	Out     string
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
	flagset.StringVar(&params.Out, "out", "", "if present outputs platter resources to directory specified, if out is - outputs to standard out")

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
	output, wasm, err := EvalPlatter(ctx, params.Release, params.Platter)
	if err != nil {
		return fmt.Errorf("failed to evaluate platter: %w", err)
	}

	var resources internal.List[*unstructured.Unstructured]
	if err := yaml.Unmarshal(output, &resources); err != nil {
		return fmt.Errorf("failed to unmarshal raw resources: %w", err)
	}

	internal.AddHallmouiMetadata(resources, params.Release)

	if params.Out != "" {
		if params.Out == "-" {
			return ExportToStdout(resources)
		}
		return ExportToFS(params.Out, params.Release, resources)
	}

	client, err := yoke.FromKubeConfig(params.KubeConfigPath)
	if err != nil {
		return err
	}

	return client.Takeoff(ctx, yoke.TakeoffParams{
		Release:   params.Release,
		Resources: resources,
		PlatterID: params.Platter.Path,
		Wasm:      wasm,
	})
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

	return xerr.MultiErrFrom("failed to write resource(s)", errs...)
}

func ExportToStdout(resources []*unstructured.Unstructured) error {
	output := make(map[string]any, len(resources))
	for _, resource := range resources {
		output[internal.Canonical(resource)] = resource.Object
	}

	encoder := y3.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	return encoder.Encode(output)
}

func EvalPlatter(ctx context.Context, release string, platter TakeoffPlatterParams) ([]byte, []byte, error) {
	if platter.Input != nil && platter.Path == "" {
		output, err := io.ReadAll(platter.Input)
		return output, nil, err
	}

	wasm, err := yoke.LoadWasm(ctx, platter.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read wasm program: %w", err)
	}

	output, err := wasi.Execute(ctx, wasm, release, platter.Input, platter.Args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute wasm: %w", err)
	}

	return output, wasm, nil
}
