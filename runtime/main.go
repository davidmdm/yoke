package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/davidmdm/x/xcontext"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// main is an example of how to extend a Go application with an addition
// function defined in WebAssembly.
//
// Since addWasm was compiled with TinyGo's `wasi` target, we need to configure
// WASI host imports.
func run() error {
	// Parse positional arguments.
	flag.Parse()

	wasm, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		return fmt.Errorf("failed to read wasm program: %w", err)
	}

	// Choose the context to use for function calls.
	ctx, done := xcontext.WithSignalCancelation(context.Background(), syscall.SIGINT)
	defer done()

	cfg := wazero.
		NewRuntimeConfig().
		WithCloseOnContextDone(true)

	// Create a new WebAssembly Runtime.
	wasi := wazero.NewRuntimeWithConfig(ctx, cfg)
	defer wasi.Close(ctx) // This closes everything this Runtime created.

	wasi_snapshot_preview1.MustInstantiate(ctx, wasi)

	// Because we are running a binary directly rather than embedding in an application,
	// we default to wiring up commonly used OS functionality.
	mod, err := wasi.CompileModule(ctx, wasm)
	if err != nil {
		return fmt.Errorf("failed to compile module: %w", err)
	}

	moduleCfg := wazero.
		NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithArgs(append([]string{flag.Arg(0)}, flag.Args()[1:]...)...)

	if _, err := wasi.InstantiateModule(ctx, mod, moduleCfg); err != nil {
		return fmt.Errorf("failed to instantiate module: %w", err)
	}

	return nil
}
