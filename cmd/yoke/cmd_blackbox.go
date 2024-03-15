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
	"github.com/pmezard/go-difflib/difflib"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/davidmdm/ansi"
	"github.com/davidmdm/yoke/internal"
	"github.com/davidmdm/yoke/internal/k8s"
)

type BlackboxParams struct {
	GlobalSettings
	Release          string
	ResourceMappings bool
	RevisionID       int
	DiffRevisionID   int
	Context          int
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
	flagset.IntVar(&params.Context, "context", 4, "number of lines of context in diff (ignored if not comparing revisions)")
	flagset.BoolVar(&params.ResourceMappings, "mapping", false, "print release to resource mappings. If present ignores all other args")
	flagset.Parse(args)

	params.Release = flagset.Arg(0)

	if revision := flagset.Arg(1); revision != "" {
		revisionID, err := strconv.Atoi(flagset.Arg(1))
		if err != nil {
			return nil, fmt.Errorf("revision must be an integer ID: %w", err)
		}
		params.RevisionID = revisionID
	}

	if revision := flagset.Arg(2); revision != "" {
		revisionID, err := strconv.Atoi(flagset.Arg(2))
		if err != nil {
			return nil, fmt.Errorf("revision to diff must be an integer ID: %w", err)
		}
		params.DiffRevisionID = revisionID
	}

	return &params, nil
}

func Blackbox(ctx context.Context, params BlackboxParams) error {
	restcfg, err := clientcmd.BuildConfigFromFlags("", params.KubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to build k8 config: %w", err)
	}

	client, err := k8s.NewClient(restcfg)
	if err != nil {
		return fmt.Errorf("failed to instantiate k8 client: %w", err)
	}

	if params.ResourceMappings {
		mappings, err := client.GetResourceReleaseMapping(ctx)
		if err != nil {
			return fmt.Errorf("failed to lookup resource to release mappings: %w", err)
		}

		relToRes := make(map[string][]string)
		for resource, release := range mappings {
			relToRes[release] = append(relToRes[release], resource)
		}

		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(relToRes)
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

	revisions, ok := internal.Find(allReleases, func(revisions internal.Revisions) bool {
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

		tbl.AppendHeader(table.Row{"id", "resources", "flight", "sha", "created at"})
		for _, version := range history {
			tbl.AppendRow(table.Row{version.ID, len(version.Resources), version.Source.Ref, version.Source.Checksum, version.CreatedAt})
		}

		_, err = io.WriteString(os.Stdout, tbl.Render()+"\n")
		return err
	}

	revision, ok := internal.Find(revisions.History, func(revision internal.Revision) bool {
		return revision.ID == params.RevisionID
	})
	if !ok {
		return fmt.Errorf("revision %d not found", params.RevisionID)
	}

	primaryRevision := make(map[string]any, len(revision.Resources))
	for _, resource := range revision.Resources {
		primaryRevision[internal.Canonical(resource)] = resource.Object
	}

	if params.DiffRevisionID == 0 {
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)

		if err := encoder.Encode(primaryRevision); err != nil {
			return fmt.Errorf("failed to encode resources: %w", err)
		}
		return nil
	}

	revision, ok = internal.Find(revisions.History, func(revision internal.Revision) bool {
		return revision.ID == params.DiffRevisionID
	})
	if !ok {
		return fmt.Errorf("revision %d not found", params.DiffRevisionID)
	}

	diffRevision := make(map[string]any, len(revision.Resources))
	for _, resource := range revision.Resources {
		diffRevision[internal.Canonical(resource)] = resource.Object
	}

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)

	if err := encoder.Encode(primaryRevision); err != nil {
		return err
	}

	a := buffer.String()

	buffer.Reset()

	if err := encoder.Encode(diffRevision); err != nil {
		return err
	}

	b := buffer.String()

	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        indentLines(difflib.SplitLines(a), "  "),
		B:        indentLines(difflib.SplitLines(b), "  "),
		FromFile: fmt.Sprintf("revision %d", params.RevisionID),
		ToFile:   fmt.Sprintf("revision %d", params.DiffRevisionID),
		Context:  params.Context,
	})

	diff = colorizeDiff(diff)

	_, err = fmt.Fprint(os.Stdout, diff)
	return err
}

var (
	green = ansi.MakeStyle(ansi.FgGreen)
	red   = ansi.MakeStyle(ansi.FgRed)
)

func colorizeDiff(value string) string {
	lines := strings.Split(value, "\n")
	colorized := make([]string, len(lines))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case '-':
			colorized[i] = green.Sprint(line)
		case '+':
			colorized[i] = red.Sprint(line)
		default:
			colorized[i] = line
		}
	}

	return strings.Join(colorized, "\n")
}

func indentLines(lines []string, indent string) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = indent + line
	}
	return result
}
