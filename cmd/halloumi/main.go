package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/davidmdm/halloumi/internal"
	"github.com/davidmdm/x/xcontext"
	"golang.org/x/term"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error()+"\n")
		if internal.IsWarning(err) {
			return
		}
		os.Exit(1)
	}
}

//go:embed cmd_help.txt
var rootHelp string

func init() {
	rootHelp = strings.TrimSpace(internal.Colorize(rootHelp))
}

func run() error {
	ctx, done := xcontext.WithSignalCancelation(context.Background(), syscall.SIGINT)
	defer done()

	var settings GlobalSettings
	RegisterGlobalFlags(flag.CommandLine, &settings)

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), rootHelp)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		return fmt.Errorf("no command provided")
	}

	switch cmd := flag.Arg(0); cmd {
	case "takeoff", "up", "apply":
		{
			var source io.Reader
			if !term.IsTerminal(int(os.Stdin.Fd())) {
				source = os.Stdin
			}
			params, err := GetTakeoffParams(settings, source, flag.Args()[1:])
			if err != nil {
				return err
			}
			return TakeOff(ctx, *params)
		}
	case "descent", "down", "restore":
		{
			params, err := GetDescentfParams(settings, flag.Args()[1:])
			if err != nil {
				return err
			}
			return Descent(ctx, *params)
		}
	case "mayday", "delete":
		{
			return fmt.Errorf("not implemented")
		}
	case "version":
		{
			fmt.Println(version())
			return nil
		}
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func version() string {
	info, _ := debug.ReadBuildInfo()
	return info.Main.Version
}
