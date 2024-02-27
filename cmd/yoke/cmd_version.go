package main

import (
	"fmt"
	"runtime/debug"
	"slices"

	"github.com/jedib0t/go-pretty/v6/table"
)

func Version() error {
	info, _ := debug.ReadBuildInfo()

	tbl := table.NewWriter()
	tbl.SetStyle(table.StyleRounded)

	tbl.AppendRow(table.Row{"yoke", info.Main.Version})

	for _, mod := range info.Deps {
		if !slices.Contains([]string{"k8s.io/client-go"}, mod.Path) {
			continue
		}
		tbl.AppendRow(table.Row{mod.Path, mod.Version})
	}

	fmt.Println(tbl.Render())

	return nil
}
