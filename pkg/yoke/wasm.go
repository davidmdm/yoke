package yoke

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/davidmdm/x/xerr"
	"github.com/davidmdm/yoke/internal"
)

func LoadWasm(ctx context.Context, path string) (wasm []byte, err error) {
	defer internal.DebugTimer(ctx, "load wasm")()

	uri, _ := url.Parse(path)
	if uri.Scheme == "" {
		wasm, err := loadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load file: %s: %w", path, err)
		}
		return wasm, nil
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
		return nil, fmt.Errorf("unexpected statuscode fetching %s: %d", uri.String(), resp.StatusCode)
	}

	if resp.Header.Get("Content-Encoding") == "gzip" || strings.HasSuffix(req.URL.Path, ".gz") {
		return io.ReadAll(gzipReader(resp.Body))
	}

	return io.ReadAll(resp.Body)
}

func loadFile(path string) (result []byte, err error) {
	if filepath.Ext(path) != ".gz" {
		return os.ReadFile(path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = xerr.MultiErrFrom("", err, file.Close())
	}()

	return io.ReadAll(gzipReader(file))
}

func gzipReader(r io.Reader) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		gr, err := gzip.NewReader(r)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err := io.Copy(pw, gr); err != nil {
			pw.CloseWithError(err)
		}
		pw.CloseWithError(gr.Close())
	}()

	return pr
}
