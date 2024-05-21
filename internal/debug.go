package internal

import (
	"context"
	"io"
	"time"

	"github.com/davidmdm/ansi"
)

type debugKey struct{}

func WithDebugFlag(ctx context.Context, debug *bool) context.Context {
	return context.WithValue(ctx, debugKey{}, debug)
}

func Debug(ctx context.Context) ansi.Terminal {
	debug, _ := ctx.Value(debugKey{}).(*bool)
	if debug == nil || !*debug {
		return ansi.Terminal{Writer: io.Discard}
	}
	return ansi.Stderr
}

func DebugTimer(ctx context.Context, msg string) func() {
	start := time.Now()
	terminal := Debug(ctx)
	terminal.Printf("start: %s\n", msg)
	return func() { terminal.Printf("done:  %s: %s\n", msg, time.Since(start)) }
}
