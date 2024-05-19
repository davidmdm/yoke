package internal

import (
	"context"
	"io"
	"os"
)

type (
	stdoutKey struct{}
	stderrKey struct{}
	stdinKey  struct{}
)

func WithStdout(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, stdoutKey{}, w)
}

func Stdout(ctx context.Context) io.Writer {
	w, ok := ctx.Value(stdoutKey{}).(io.Writer)
	if !ok {
		return os.Stdout
	}
	return w
}

func WithStderr(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, stderrKey{}, w)
}

func Stderr(ctx context.Context) io.Writer {
	w, ok := ctx.Value(stderrKey{}).(io.Writer)
	if !ok {
		return os.Stderr
	}
	return w
}

func WithStdin(ctx context.Context, r io.Reader) context.Context {
	return context.WithValue(ctx, stdinKey{}, r)
}

func Stdin(ctx context.Context) io.Reader {
	r, ok := ctx.Value(stdinKey{}).(io.Reader)
	if !ok {
		return os.Stdin
	}
	return r
}
