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
	resources, err := helm.Render(helm.Params{
		Source:      pg,
		ReleaseName: "foo",
		Namespace:   "default",
		Values:      nil,
	})
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(resources)
}
