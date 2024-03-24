package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/davidmdm/yoke/internal"
)

type ExecParams struct {
	GlobalSettings
	Release string
	Flight  TakeoffFlightParams
}

func init() {
	takeoffHelp = strings.TrimSpace(internal.Colorize(takeoffHelp))
}

func GetExecParams(settings GlobalSettings, source io.Reader, args []string) (*ExecParams, error) {
	flagset := flag.NewFlagSet("exec", flag.ExitOnError)

	flagset.Usage = func() {
		// fmt.Fprintln(flagset.Output(), takeoffHelp)
		flagset.PrintDefaults()
	}

	params := ExecParams{
		GlobalSettings: settings,
		Flight:         TakeoffFlightParams{Input: source},
	}

	RegisterGlobalFlags(flagset, &params.GlobalSettings)

	args, params.Flight.Args = internal.CutArgs(args)

	flagset.Parse(args)

	params.Release = flagset.Arg(0)
	params.Flight.Path = flagset.Arg(1)

	if params.Release == "" {
		return nil, fmt.Errorf("release is required as first positional arg")
	}
	if params.Flight.Input == nil && params.Flight.Path == "" {
		return nil, fmt.Errorf("flight-path is required as second position arg")
	}

	return &params, nil
}

func Exec(ctx context.Context, params ExecParams) error {
	output, _, err := EvalFlight(ctx, params.Release, params.Flight)
	if err != nil {
		return fmt.Errorf("failed to evaluate flight: %w", err)
	}
	_, err = fmt.Print(string(output))
	return err
}
