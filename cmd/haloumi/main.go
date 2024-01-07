package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"syscall"

	"wasiplaytime/internal/wasi"

	"github.com/davidmdm/x/xcontext"
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

	out, err := wasi.Execute(ctx, wasm)
	if err != nil {
		return fmt.Errorf("failed to execute wasm: %w", err)
	}

	fmt.Println(out)

	return nil
}
