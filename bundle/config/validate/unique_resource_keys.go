package validate

import (
	"context"
	"fmt"
	"slices"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/libs/diag"
	"github.com/databricks/cli/libs/dyn"
)

// This mutator validates that:
//
//  1. Each resource key is unique across different resource types. No two resources
//     of the same type can have the same key. This is because command like "bundle run"
//     rely on the resource key to identify the resource to run.
//     Eg: jobs.foo and pipelines.foo are not allowed simultaneously.
//
//  2. Each resource definition is contained within a single file, and is not spread
//     across multiple files. Note: This is not applicable to resource configuration
//     defined in a target override. That is why this mutator MUST run before the target
//     overrides are merged.
func UniqueResourceKeys() bundle.Mutator {
	return &uniqueResourceKeys{}
}

type uniqueResourceKeys struct{}

func (m *uniqueResourceKeys) Name() string {
	return "validate:unique_resource_keys"
}

func (m *uniqueResourceKeys) Apply(ctx context.Context, b *bundle.Bundle) diag.Diagnostics {
	diags := diag.Diagnostics{}

	// Maps of resource key to the paths and locations the resource is defined at.
	pathsByKey := map[string][]dyn.Path{}
	locationsByKey := map[string][]dyn.Location{}

	rv := b.Config.Value().Get("resources")

	// return early if no resources are defined or the resources block is empty.
	if rv.Kind() == dyn.KindInvalid || rv.Kind() == dyn.KindNil {
		return diags
	}

	// Gather the paths and locations of all resources.
	_, err := dyn.MapByPattern(
		rv,
		dyn.NewPattern(dyn.AnyKey(), dyn.AnyKey()),
		func(p dyn.Path, v dyn.Value) (dyn.Value, error) {
			// The key for the resource. Eg: "my_job" for jobs.my_job.
			k := p[1].Key()

			// dyn.Path under the hood is a slice. So, we need to clone it.
			pathsByKey[k] = append(pathsByKey[k], slices.Clone(p))

			locationsByKey[k] = append(locationsByKey[k], v.Locations()...)
			return v, nil
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, locations := range locationsByKey {
		if len(locations) <= 1 {
			continue
		}

		// If there are multiple resources with the same key, report an error.
		diags = append(diags, diag.Diagnostic{
			Severity:  diag.Error,
			Summary:   fmt.Sprintf("multiple resources have been defined with the same key: %s", k),
			Locations: locations,
			Paths:     pathsByKey[k],
		})
	}

	return diags
}
