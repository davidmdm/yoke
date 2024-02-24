package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/davidmdm/halloumi/internal/k8"
)

type BlackboxParams struct {
	GlobalSettings
	Release    string
	RevisionID int
}

func GetBlackBoxParams(settings GlobalSettings, args []string) (*BlackboxParams, error) {
	flagset := flag.NewFlagSet("blackbox", flag.ExitOnError)

	flagset.Usage = func() {
		// fmt.Fprintln(flagset.Output(), descentHelp)
		flagset.PrintDefaults()
	}

	params := BlackboxParams{GlobalSettings: settings}

	RegisterGlobalFlags(flagset, &params.GlobalSettings)

	flagset.Parse(args)

	params.Release = flagset.Arg(0)

	if revision := flag.Arg(1); revision != "" {
		revisionID, err := strconv.Atoi(flagset.Arg(1))
		if err != nil {
			return nil, fmt.Errorf("revision must be an integer ID: %w", err)
		}
		params.RevisionID = revisionID
	}

	return &params, nil
}

func Blackbox(ctx context.Context, params BlackboxParams) error {
	restcfg, err := clientcmd.BuildConfigFromFlags("", params.KubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to build k8 config: %w", err)
	}

	client, err := k8.NewClient(restcfg)
	if err != nil {
		return fmt.Errorf("failed to instantiate k8 client: %w", err)
	}

	if params.Release == "" {
		revisions, err := client.GetAllRevisions(ctx)
		if err != nil {
			return fmt.Errorf("failed to get revisions: %w", err)
		}

		tbl := table.NewWriter()
		tbl.SetStyle(table.StyleRounded)

		tbl.AppendHeader(table.Row{"release", "revision id"})
		for _, revision := range revisions {
			tbl.AppendRow(table.Row{revision.Release, revision.History[revision.ActiveIndex].ID})
		}

		_, err = io.WriteString(os.Stdout, tbl.Render()+"\n")
		return err
	}

	return fmt.Errorf("not implemented")
}
