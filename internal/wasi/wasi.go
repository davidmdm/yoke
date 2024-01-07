package wasi

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/davidmdm/x/xerr"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func Execute(ctx context.Context, wasm []byte) (output string, err error) {
	cfg := wazero.
		NewRuntimeConfig().
		WithCloseOnContextDone(true)

	// Create a new WebAssembly Runtime.
	wasi := wazero.NewRuntimeWithConfig(ctx, cfg)
	defer func() {
		err = xerr.MultiErrFrom("", err, wasi.Close(ctx))
	}()

	wasi_snapshot_preview1.MustInstantiate(ctx, wasi)

	// Because we are running a binary directly rather than embedding in an application,
	// we default to wiring up commonly used OS functionality.
	mod, err := wasi.CompileModule(ctx, wasm)
	if err != nil {
		return "", fmt.Errorf("failed to compile module: %w", err)
	}

	var (
		stdout strings.Builder
		stderr strings.Builder
	)

	moduleCfg := wazero.
		NewModuleConfig().
		WithStdout(&stdout).
		WithStderr(&stderr).
		WithRandSource(rand.Reader).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithArgs("exe")

	if _, err := wasi.InstantiateModule(ctx, mod, moduleCfg); err != nil {
		details := stderr.String()
		if details == "" {
			details = "(no output captured on stderr)"
		}
		return "", fmt.Errorf("failed to instantiate module: %w: stderr: %s", err, details)
	}

	return stdout.String(), nil
}
