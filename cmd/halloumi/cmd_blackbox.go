package main

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/davidmdm/halloumi/internal"
	"github.com/davidmdm/halloumi/internal/k8"
)

type BlackboxParams struct {
	GlobalSettings
	Release    string
	RevisionID int
}

//go:embed cmd_blackbox_help.txt
var blackboxHelp string

func init() {
	blackboxHelp = strings.TrimSpace(internal.Colorize(blackboxHelp))
}

func GetBlackBoxParams(settings GlobalSettings, args []string) (*BlackboxParams, error) {
	flagset := flag.NewFlagSet("blackbox", flag.ExitOnError)

	flagset.Usage = func() {
		fmt.Fprintln(flagset.Output(), blackboxHelp)
		flagset.PrintDefaults()
	}

	params := BlackboxParams{GlobalSettings: settings}

	RegisterGlobalFlags(flagset, &params.GlobalSettings)

	flagset.Parse(args)

	params.Release = flagset.Arg(0)

	if revision := flagset.Arg(1); revision != "" {
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

	allReleases, err := client.GetAllRevisions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get revisions: %w", err)
	}
	if params.Release == "" {
		tbl := table.NewWriter()
		tbl.SetStyle(table.StyleRounded)

		tbl.AppendHeader(table.Row{"release", "revision id"})
		for _, revisions := range allReleases {
			tbl.AppendRow(table.Row{revisions.Release, revisions.History[revisions.ActiveIndex].ID})
		}

		_, err = io.WriteString(os.Stdout, tbl.Render()+"\n")
		return err
	}

	revisions, ok := Find(allReleases, func(revisions internal.Revisions) bool {
		return revisions.Release == params.Release
	})
	if !ok {
		return fmt.Errorf("release %q not found", params.Release)
	}

	if params.RevisionID == 0 {
		tbl := table.NewWriter()
		tbl.SetStyle(table.StyleRounded)

		history := revisions.History
		slices.Reverse(history)

		tbl.AppendHeader(table.Row{"id", "resources", "platter", "sha", "created at"})
		for _, version := range history {
			tbl.AppendRow(table.Row{version.ID, len(version.Resources), version.Source.Ref, version.Source.Checksum, version.CreatedAt})
		}

		_, err = io.WriteString(os.Stdout, tbl.Render()+"\n")
		return err
	}

	revision, ok := Find(revisions.History, func(revision internal.Revision) bool {
		return revision.ID == params.RevisionID
	})
	if !ok {
		return fmt.Errorf("revision %d not found", params.RevisionID)
	}

	output := make(map[string]any, len(revision.Resources))
	for _, resource := range revision.Resources {
		output[internal.Canonical(resource)] = resource.Object
	}

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)

	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to encode resources: %w", err)
	}

	_, err = fmt.Fprint(os.Stdout, buffer.String())
	return err
}

func Find[S ~[]E, E any](slice S, fn func(E) bool) (E, bool) {
	switch idx := slices.IndexFunc(slice, fn); idx {
	case -1:
		var zero E
		return zero, false
	default:
		return slice[idx], true
	}
}
