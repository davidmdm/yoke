package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"github.com/davidmdm/halloumi/internal"
	"github.com/davidmdm/halloumi/internal/k8"
	"k8s.io/client-go/tools/clientcmd"
)

type DescentParams struct {
	GlobalSettings
	Release    string
	RevisionID int
}

func GetDescentfParams(settings GlobalSettings, args []string) (*DescentParams, error) {
	flagset := flag.NewFlagSet("descent", flag.ExitOnError)

	flagset.Usage = func() {
		// fmt.Fprintln(flagset.Output(), takeoffHelp)
		flagset.PrintDefaults()
	}

	params := DescentParams{
		GlobalSettings: settings,
	}

	RegisterGlobalFlags(flagset, &params.GlobalSettings)

	flagset.Parse(args)

	params.Release = flag.Arg(0)
	if params.Release == "" {
		return nil, fmt.Errorf("release is required as first positional arg")
	}

	if len(flag.Args()) < 2 {
		return nil, fmt.Errorf("revision is required as second positional arg")
	}

	revisionID, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		return nil, fmt.Errorf("revision must be an integer ID: %w", err)
	}

	params.RevisionID = revisionID

	return &params, nil
}

func Descent(ctx context.Context, params DescentParams) error {
	restcfg, err := clientcmd.BuildConfigFromFlags("", params.KubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to build k8 config: %w", err)
	}

	client, err := k8.NewClient(restcfg)
	if err != nil {
		return fmt.Errorf("failed to instantiate k8 client: %w", err)
	}

	revisions, err := client.GetRevisions(ctx, params.Release)
	if err != nil {
		return fmt.Errorf("failed to get revisions: %w", err)
	}

	index, next := func() (int, *internal.Revision) {
		for i, revision := range revisions.History {
			if revision.ID == params.RevisionID {
				return i, &revision
			}
		}
		return 0, nil
	}()

	if next == nil {
		return fmt.Errorf("revision %d is not within history", params.RevisionID)
	}

	previous := revisions.CurrentResources()

	if err := client.ApplyResources(ctx, next.Resources); err != nil {
		return fmt.Errorf("failed to apply resources: %w", err)
	}

	revisions.ActiveIndex = index

	if err := client.UpsertRevisions(ctx, params.Release, revisions); err != nil {
		return fmt.Errorf("failed to update revision history: %w", err)
	}

	if err := client.RemoveOrphans(ctx, previous, next.Resources); err != nil {
		return fmt.Errorf("failed to remove orphaned resources: %w", err)
	}

	return nil
}
