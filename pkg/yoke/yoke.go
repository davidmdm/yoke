package yoke

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/davidmdm/yoke/internal"
	"github.com/davidmdm/yoke/internal/k8s"
)

type Client struct {
	k8s *k8s.Client
}

func FromKubeConfig(path string) (*Client, error) {
	client, err := k8s.NewClientFromKubeConfig(path)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize k8s client: %w", err)
	}
	return &Client{client}, nil
}

type TakeoffParams struct {
	Release    string
	Resources  []*unstructured.Unstructured
	FlightID   string
	Wasm       []byte
	SkipDryRun bool
}

func (client Client) Takeoff(ctx context.Context, params TakeoffParams) error {
	revisions, err := client.k8s.GetRevisions(ctx, params.Release)
	if err != nil {
		return fmt.Errorf("failed to get revision history: %w", err)
	}

	previous := revisions.CurrentResources()

	if reflect.DeepEqual(previous, params.Resources) {
		return internal.Warning("resources are the same as previous revision: skipping takeoff")
	}

	if err := client.k8s.ValidateOwnership(ctx, params.Release, params.Resources); err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}

	applyOpts := k8s.ApplyResourcesOpts{SkipDryRun: params.SkipDryRun}
	if err := client.k8s.ApplyResources(ctx, params.Resources, applyOpts); err != nil {
		return fmt.Errorf("failed to apply resources: %w", err)
	}

	revisions.Add(params.Resources, params.FlightID, params.Wasm)

	if err := client.k8s.UpsertRevisions(ctx, params.Release, revisions); err != nil {
		return fmt.Errorf("failed to create revision: %w", err)
	}

	removed, err := client.k8s.RemoveOrphans(ctx, previous, params.Resources)
	if err != nil {
		return fmt.Errorf("failed to remove orhpans: %w", err)
	}

	var (
		createdNames = internal.CanonicalNameList(params.Resources)
		removedNames = internal.CanonicalNameList(removed)
	)

	if err := client.k8s.UpdateResourceReleaseMapping(ctx, params.Release, createdNames, removedNames); err != nil {
		return fmt.Errorf("failed to update resource release mapping: %w", err)
	}

	return nil
}

type DescentParams struct {
	Release    string
	RevisionID int
}

func (client Client) Descent(ctx context.Context, params DescentParams) error {
	revisions, err := client.k8s.GetRevisions(ctx, params.Release)
	if err != nil {
		return fmt.Errorf("failed to get revisions: %w", err)
	}

	index, next := func() (int, *internal.Revision) {
		for i, revision := range revisions.History {
			if revision.ID == params.RevisionID {
				return i, &revision
			}
		}
		return 0, nil
	}()

	if next == nil {
		return fmt.Errorf("revision %d is not within history", params.RevisionID)
	}

	if err := client.k8s.ValidateOwnership(ctx, params.Release, next.Resources); err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}

	previous := revisions.CurrentResources()

	if err := client.k8s.ApplyResources(ctx, next.Resources, k8s.ApplyResourcesOpts{SkipDryRun: true}); err != nil {
		return fmt.Errorf("failed to apply resources: %w", err)
	}

	revisions.ActiveIndex = index

	if err := client.k8s.UpsertRevisions(ctx, params.Release, revisions); err != nil {
		return fmt.Errorf("failed to update revision history: %w", err)
	}

	removed, err := client.k8s.RemoveOrphans(ctx, previous, next.Resources)
	if err != nil {
		return fmt.Errorf("failed to remove orphaned resources: %w", err)
	}

	var (
		createdNames = internal.CanonicalNameList(next.Resources)
		removedNames = internal.CanonicalNameList(removed)
	)

	if err := client.k8s.UpdateResourceReleaseMapping(ctx, params.Release, createdNames, removedNames); err != nil {
		return fmt.Errorf("failed to update resource release mapping: %w", err)
	}

	return nil
}

func (client Client) Mayday(ctx context.Context, release string) error {
	revisions, err := client.k8s.GetRevisions(ctx, release)
	if err != nil {
		return fmt.Errorf("failed to get revision history for release: %w", err)
	}

	removed, err := client.k8s.RemoveOrphans(ctx, revisions.CurrentResources(), nil)
	if err != nil {
		return fmt.Errorf("failed to delete resources: %w", err)
	}

	if err := client.k8s.UpdateResourceReleaseMapping(ctx, release, nil, internal.CanonicalNameList(removed)); err != nil {
		return fmt.Errorf("failed to update resource to release mapping: %w", err)
	}

	if err := client.k8s.DeleteRevisions(ctx, release); err != nil {
		return fmt.Errorf("failed to delete revision history: %w", err)
	}

	return nil
}
