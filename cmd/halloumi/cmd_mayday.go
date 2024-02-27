package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/davidmdm/halloumi/internal"
	"github.com/davidmdm/halloumi/internal/k8s"
)

type MaydayParams struct {
	GlobalSettings
	Release string
}

//go:embed cmd_mayday_help.txt
var maydayHelp string

func init() {
	maydayHelp = strings.TrimSpace(internal.Colorize(maydayHelp))
}

func GetMaydayParams(settings GlobalSettings, args []string) (*MaydayParams, error) {
	flagset := flag.NewFlagSet("mayday", flag.ExitOnError)

	flagset.Usage = func() {
		fmt.Fprintln(flagset.Output(), maydayHelp)
		flagset.PrintDefaults()
	}

	params := MaydayParams{GlobalSettings: settings}

	RegisterGlobalFlags(flagset, &params.GlobalSettings)

	flagset.Parse(args)

	params.Release = flagset.Arg(0)
	if params.Release == "" {
		return nil, fmt.Errorf("release is required")
	}

	return &params, nil
}

func Mayday(ctx context.Context, params MaydayParams) error {
	restcfg, err := clientcmd.BuildConfigFromFlags("", params.KubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to build k8 config: %w", err)
	}

	client, err := k8s.NewClient(restcfg)
	if err != nil {
		return fmt.Errorf("failed to instantiate k8 client: %w", err)
	}

	revisions, err := client.GetRevisions(ctx, params.Release)
	if err != nil {
		return fmt.Errorf("failed to get revision history for release: %w", err)
	}

	removed, err := client.RemoveOrphans(ctx, revisions.CurrentResources(), nil)
	if err != nil {
		return fmt.Errorf("failed to delete resources: %w", err)
	}

	if err := client.UpdateResourceReleaseMapping(ctx, params.Release, nil, internal.CanonicalNameList(removed)); err != nil {
		return fmt.Errorf("failed to update resource to release mapping: %w", err)
	}

	if err := client.DeleteRevisions(ctx, params.Release); err != nil {
		return fmt.Errorf("failed to delete revision history: %w", err)
	}

	return nil
}
