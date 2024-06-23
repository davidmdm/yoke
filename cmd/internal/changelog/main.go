package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	out := flag.String("o", "./CHANGELOG.md", "location to write file to")
	flag.Parse()

	repo, err := git.PlainOpen(".")
	if err != nil {
		return fmt.Errorf("failed to open git repo: %w", err)
	}

	iter, err := repo.Tags()
	if err != nil {
		return fmt.Errorf("failed to read tags: %w", err)
	}

	tags := map[string]string{}

	iter.ForEach(func(r *plumbing.Reference) error {
		hash := r.Hash().String()
		current := strings.TrimPrefix(r.Name().String(), "refs/tags/")

		if previous := tags[hash]; previous != "" {
			tags[hash] = previous + " " + current
			return nil
		}

		tags[hash] = current
		return nil
	})

	commits, err := repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		return fmt.Errorf("failed to log: %w", err)
	}

	var changelog Changelog

	commits.ForEach(func(c *object.Commit) error {
		if tag := tags[c.Hash.String()]; tag != "" {
			changelog = append(changelog, Entry{Tag: tag})
		}
		if len(changelog) == 0 {
			return nil
		}
		changelog[len(changelog)-1].Commits = append(changelog[len(changelog)-1].Commits, Commit{
			Msg: strings.TrimSpace(c.Message),
			Sha: c.Hash.String(),
		})
		return nil
	})

	if *out == "-" {
		fmt.Println(changelog)
		return nil
	}

	return os.WriteFile(*out, []byte(changelog.String()), 0o644)
}

type (
	Changelog []Entry
	Commit    struct {
		Msg string
		Sha string
	}
	Entry struct {
		Tag     string
		Commits []Commit
	}
)

func (changelog Changelog) String() string {
	var builder strings.Builder
	builder.WriteString("# Changelog\n\n")

	builder.WriteString("> [!IMPORTANT]\n")
	builder.WriteString("> This project has not reached v1.0.0 and as such provides no backwards compatibility guarantees between versions.\n")
	builder.WriteString("> Pre v1.0.0 minor bumps will repesent breaking changes.\n\n")

	for _, entry := range changelog {
		builder.WriteString("## " + entry.Tag + "\n\n")

		if slices.ContainsFunc(entry.Commits, func(commit Commit) bool {
			msg := strings.ToLower(commit.Msg)
			return strings.Contains(msg, "breaking change")
		}) {
			builder.WriteString("> [!CAUTION]\n")
			builder.WriteString("> This version contains breaking changes, and is not expected to be compatible with previous versions\n\n")
		}

		for _, commit := range entry.Commits {
			builder.WriteString(fmt.Sprintf("* %s ([%s](https://github.com/davidmdm/yoke/commit/%s))\n", commit.Msg, commit.Sha[:7], commit.Sha))
		}
		builder.WriteByte('\n')
	}

	return builder.String()
}
