package yoke

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"

	"github.com/davidmdm/x/xerr"
)

func LoadWasm(ctx context.Context, path string) (wasm []byte, err error) {
	uri, _ := url.Parse(path)
	if uri.Scheme == "" {
		return os.ReadFile(path)
	}

	if !slices.Contains([]string{"http", "https"}, uri.Scheme) {
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

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("unexpected statuscode fetching %s: %w", uri.String(), err)
	}

	return io.ReadAll(resp.Body)
}
