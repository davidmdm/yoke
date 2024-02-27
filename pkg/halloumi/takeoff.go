package halloumi

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/davidmdm/halloumi/internal"
)

type TakeoffParams struct {
	Release   string
	Resources []*unstructured.Unstructured
	PlatterID string
	Wasm      []byte
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

	if err := client.k8s.ApplyResources(ctx, params.Resources); err != nil {
		return fmt.Errorf("failed to apply resources: %w", err)
	}

	revisions.Add(params.Resources, params.PlatterID, params.Wasm)

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
