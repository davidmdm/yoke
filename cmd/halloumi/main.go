package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/davidmdm/x/xcontext"
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

	var settings GlobalSettings
	RegisterGlobalFlags(flag.CommandLine, &settings)

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		return fmt.Errorf("no command provided")
	}

	switch cmd := flag.Arg(0); cmd {
	case "takeoff", "up":
		{
			params, err := GetTakeoffParams(settings, flag.Args()[1:])
			if err != nil {
				return err
			}
			return TakeOff(ctx, *params)
		}
	case "descent", "down", "rollback":
		{
			return Descent(ctx)
		}
	case "runway", "render":
		{
			return Runway(ctx)
		}
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}
