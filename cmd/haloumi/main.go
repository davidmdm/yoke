package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"syscall"

	"github.com/davidmdm/x/xcontext"
	"github.com/davidmdm/x/xerr"

	"github.com/davidmdm/haloumi/internal/wasi"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}

func run() error {
	ctx, done := xcontext.WithSignalCancelation(context.Background(), syscall.SIGINT)
	defer done()

	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	wasm, err := LoadWasm(ctx, cfg.PackagePath)
	if err != nil {
		return fmt.Errorf("failed to read wasm program: %w", err)
	}

	output, err := wasi.Execute(ctx, wasm, cfg.ReleaseName, cfg.PackageFlags...)
	if err != nil {
		return fmt.Errorf("failed to execute wasm: %w", err)
	}

	var resources []any
	if err := json.Unmarshal(output, &resources); err != nil {
		return fmt.Errorf("failed to unmarshal raw resources: %w", err)
	}

	if err := ApplyResources(ctx, resources); err != nil {
		return fmt.Errorf("failed to apply resources: %w", err)
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

// TODO
func ApplyResources(ctx context.Context, resources []any) error {
	fmt.Println(resources...)
	return nil
}
