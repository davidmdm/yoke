package main

import (
	"encoding/json"
	"os"

	"github.com/davidmdm/yoke/cmd/examples/internal/flights/argocd"
)

func main() {
	resources, err := argocd.RenderChart(os.Args[0], "argocd", nil)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(os.Stdout).Encode(resources)
}
