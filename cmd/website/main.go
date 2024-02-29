package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/davidmdm/x/xcontext"
)

//go:embed content
var content embed.FS

//go:embed content/index.html
var root []byte

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := xcontext.WithSignalCancelation(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) { w.Write(root) })
	mux.Handle("GET /content/", http.FileServerFS(content))

	server := http.Server{
		Addr:              ":3000",
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			errCh <- fmt.Errorf("failed to start server: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	if err := server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return context.Cause(ctx)
}
