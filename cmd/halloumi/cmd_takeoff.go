package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/davidmdm/x/xerr"

	"github.com/davidmdm/halloumi/internal"
	"github.com/davidmdm/halloumi/internal/k8"
	"github.com/davidmdm/halloumi/internal/wasi"
)

type TakeoffParams struct {
	GlobalSettings
	ReleaseName string
	PlatterPath string
	PlatterArgs []string
}

//go:embed cmd_takeoff_help.txt
var takeoffHelp string

func init() {
	takeoffHelp = strings.TrimSpace(internal.Colorize(takeoffHelp))
}

func GetTakeoffParams(settings GlobalSettings, args []string) (*TakeoffParams, error) {
	flagset := flag.NewFlagSet("takeoff", flag.ExitOnError)

	flagset.Usage = func() {
		fmt.Fprintln(flagset.Output(), takeoffHelp)
		flagset.PrintDefaults()
	}

	params := TakeoffParams{GlobalSettings: settings}

	RegisterGlobalFlags(flagset, &params.GlobalSettings)

	args, params.PlatterArgs = internal.CutArgs(args)

	flagset.Parse(args)

	params.ReleaseName = flagset.Arg(0)
	params.PlatterPath = flagset.Arg(1)

	if params.ReleaseName == "" {
		return nil, fmt.Errorf("release is required as first positional arg")
	}
	if params.PlatterPath == "" {
		return nil, fmt.Errorf("platter-path is required as second position arg")
	}

	return &params, nil
}

func TakeOff(ctx context.Context, params TakeoffParams) error {
	wasm, err := LoadWasm(ctx, params.PlatterPath)
	if err != nil {
		return fmt.Errorf("failed to read wasm program: %w", err)
	}

	output, err := wasi.Execute(ctx, wasm, params.ReleaseName, params.PlatterArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute wasm: %w", err)
	}

	var resources []*unstructured.Unstructured
	if err := json.Unmarshal(output, &resources); err != nil {
		return fmt.Errorf("failed to unmarshal raw resources: %w", err)
	}

	internal.AddHallmouiMetadata(resources, params.ReleaseName)

	restcfg, err := clientcmd.BuildConfigFromFlags("", params.KubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to build k8 config: %w", err)
	}

	client, err := k8.NewClient(restcfg)
	if err != nil {
		return fmt.Errorf("failed to instantiate k8 client: %w", err)
	}

	if err := client.ApplyResources(ctx, resources); err != nil {
		return fmt.Errorf("failed to apply resources: %w", err)
	}

	if err := client.MakeRevision(ctx, params.ReleaseName, resources); err != nil {
		return fmt.Errorf("failed to create revision: %w", err)
	}

	if err := client.RemoveOrphans(ctx, params.ReleaseName); err != nil {
		return fmt.Errorf("failed to remove orhpans: %w", err)
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
