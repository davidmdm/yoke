package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/davidmdm/halloumi/pkg/helm"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//go:embed postgresql-14.2.3.tgz
var pg []byte

func run() error {
	chart, err := helm.LoadChartFromZippedArchive(pg)
	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	resources, err := chart.Render(os.Args[0], "default", nil)
	if err != nil {
		return fmt.Errorf("failed to render chart resources: %w", err)
	}

	return json.NewEncoder(os.Stdout).Encode(resources)
}
