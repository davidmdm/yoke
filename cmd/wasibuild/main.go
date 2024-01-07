package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	flag.Parse()

	cmd := exec.Command("go", "build", "-o", flag.Arg(1), flag.Arg(0))
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
