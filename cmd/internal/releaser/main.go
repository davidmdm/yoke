package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/mod/semver"

	"github.com/davidmdm/x/xerr"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	dry := flag.Bool("dry", false, "dry-run")

	var commands []string
	flag.Func("cmd", "commands to buid as wasm and release", func(value string) error {
		commands = append(commands, strings.Split(value, ",")...)
		return nil
	})

	flag.Parse()

	repo, err := git.PlainOpen(".")
	if err != nil {
		return fmt.Errorf("failed to open git repo: %w", err)
	}

	iter, err := repo.Tags()
	if err != nil {
		return fmt.Errorf("failed to read tags: %w", err)
	}

	versions := map[string]string{}

	iter.ForEach(func(r *plumbing.Reference) error {
		release, version := path.Split(r.Name()[len("refs/tags/"):].String())
		if !semver.IsValid(version) {
			return nil
		}
		release = path.Clean(release)
		if semver.Compare(version, versions[release]) > 0 {
			versions[release] = version
		}
		return nil
	})

	releaser := Releaser{
		Versions: versions,
		Repo:     repo,
		DryRun:   *dry,
	}

	var errs []error
	for _, cmd := range commands {
		if err := releaser.ReleaseWasmBinary(cmd); err != nil {
			errs = append(errs, fmt.Errorf("failed to release %s: %w", cmd, err))
		}
	}

	return xerr.MultiErrOrderedFrom("", errs...)
}

type Releaser struct {
	Versions map[string]string
	Repo     *git.Repository
	DryRun   bool
}

func (releaser Releaser) ReleaseWasmBinary(name string) (err error) {
	version := releaser.Versions[name]

	if version == "" {
		fmt.Println("No version found for", name)
	} else if semver.Compare(version, "v0.0.1") < 0 {
		fmt.Printf("%s is pre v0.0.1... tagging new release\n", name)
	} else if version != "" {
		diff, err := releaser.HasDiff(name)
		if err != nil {
			return fmt.Errorf("failed to check for diff: %w", err)
		}

		if !diff {
			return nil
		}
	}

	outputPath, err := build(filepath.Join("cmd", name))
	if err != nil {
		return fmt.Errorf("failed to build wasm: %w", err)
	}

	outputPath, err = compress(outputPath)
	if err != nil {
		return fmt.Errorf("failed to compress wasm: %w", err)
	}

	version = func() string {
		patch := semver.Canonical(version)[len(semver.MajorMinor(version)):]
		patchNum, _ := strconv.Atoi(patch)
		patchNum++
		return fmt.Sprintf("v0.0.%d", patchNum)
	}()

	tag := fmt.Sprintf("%s/%s", name, version)

	if releaser.DryRun {
		fmt.Println("dry-run: create realease", tag)
		return nil
	}
	if err := release(tag, outputPath); err != nil {
		return fmt.Errorf("failed to release: %w", err)
	}

	return nil
}

func (releaser Releaser) HasDiff(name string) (bool, error) {
	tag := path.Join(name, releaser.Versions[name])

	hash, err := releaser.Repo.ResolveRevision(plumbing.Revision(plumbing.NewTagReferenceName(tag)))
	if err != nil {
		return false, fmt.Errorf("failed to resolve: %s: %w", tag, err)
	}

	commit, err := releaser.Repo.CommitObject(*hash)
	if err != nil {
		return false, fmt.Errorf("failed to get commit for tag %s: %w", tag, err)
	}

	h, err := releaser.Repo.Head()
	if err != nil {
		return false, fmt.Errorf("failed to resolve head: %w", err)
	}

	head, err := releaser.Repo.CommitObject(h.Hash())
	if err != nil {
		return false, fmt.Errorf("failed to get head commit: %w", err)
	}

	headTree, err := head.Tree()
	if err != nil {
		return false, fmt.Errorf("failed to get tree for head commit: %w", err)
	}

	commitTree, err := commit.Tree()
	if err != nil {
		return false, fmt.Errorf("failed to get tree for tag commit: %w", err)
	}

	changes, err := headTree.Diff(commitTree)
	if err != nil {
		return false, fmt.Errorf("failed to diff trees: %w", err)
	}

	var changed bool
	for _, change := range changes {
		var (
			from = change.From.Name
			to   = change.To.Name
		)
		if strings.HasPrefix(from, "cmd/"+name+"/") || strings.HasPrefix(to, "cmd/"+name+"/") {
			fmt.Printf("detected change from %s to %s\n", from, to)
			changed = true
			break
		}
	}

	return changed, nil
}

func build(path string) (string, error) {
	_, name := filepath.Split(path)

	out := name + ".wasm"

	cmd := exec.Command("go", "build", "-o", out, "./"+path)
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, output)
	}
	return out, nil
}

func release(tag, path string) error {
	out, err := exec.Command("gh", "release", "create", tag, path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, out)
	}
	return nil
}

func compress(path string) (out string, err error) {
	output := path + ".gz"

	destination, err := os.Create(output)
	if err != nil {
		return "", err
	}
	defer func() {
		err = xerr.MultiErrFrom("", err, destination.Close())
	}()

	compressor, err := gzip.NewWriterLevel(destination, gzip.BestCompression)
	if err != nil {
		return "", fmt.Errorf("could not create gzip writer: %w", err)
	}
	defer func() {
		err = xerr.MultiErrFrom("", err, compressor.Close())
	}()

	source, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		err = xerr.MultiErrFrom("", err, source.Close())
	}()

	if _, err := io.Copy(compressor, source); err != nil {
		return "", err
	}

	return output, nil
}
