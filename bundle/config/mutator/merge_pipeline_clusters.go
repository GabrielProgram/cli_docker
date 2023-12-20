package mutator

import (
	"context"
	"strings"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/libs/config"
	"github.com/databricks/cli/libs/config/merge"
)

type mergePipelineClusters struct{}

func MergePipelineClusters() bundle.Mutator {
	return &mergePipelineClusters{}
}

func (m *mergePipelineClusters) Name() string {
	return "MergePipelineClusters"
}

func (m *mergePipelineClusters) clusterLabel(cluster config.Value) (label string) {
	v := cluster.Get("label")
	if v == config.NilValue {
		return "default"
	}

	if v.Kind() != config.KindString {
		panic("cluster label must be a string")
	}

	return strings.ToLower(v.MustString())
}

// mergeClustersForPipeline merges cluster definitions with same label.
// The clusters field is a slice, and as such, overrides are appended to it.
// We can identify a cluster by its label, however, so we can use this label
// to figure out which definitions are actually overrides and merge them.
//
// Note: the cluster label is optional and defaults to 'default'.
// We therefore ALSO merge all clusters without a label.
func (m *mergePipelineClusters) mergeClustersForPipeline(v config.Value) (config.Value, error) {
	// We know the type of this value is a sequence.
	// For additional defence, return self if it is not.
	clusters, ok := v.AsSequence()
	if !ok {
		return v, nil
	}

	seen := make(map[string]config.Value, len(clusters))
	labels := make([]string, 0, len(clusters))

	// Target overrides are always appended, so we can iterate in natural order to
	// first find the base definition, and merge instances we encounter later.
	for i := range clusters {
		label := m.clusterLabel(clusters[i])

		// Register pipeline cluster with label if not yet seen before.
		ref, ok := seen[label]
		if !ok {
			labels = append(labels, label)
			seen[label] = clusters[i]
			continue
		}

		// Merge this instance into the reference.
		var err error
		seen[label], err = merge.Merge(ref, clusters[i])
		if err != nil {
			return v, err
		}
	}

	// Gather resulting clusters in natural order.
	out := make([]config.Value, 0, len(labels))
	for _, label := range labels {
		// Overwrite the label with the normalized version.
		nv, err := seen[label].Set("label", config.V(label))
		if err != nil {
			return config.InvalidValue, err
		}
		out = append(out, nv)
	}

	return config.NewValue(out, v.Location()), nil
}

func (m *mergePipelineClusters) foreachPipeline(v config.Value) (config.Value, error) {
	pipelines, ok := v.AsMap()
	if !ok {
		return v, nil
	}

	out := make(map[string]config.Value)
	for key, pipeline := range pipelines {
		var err error
		out[key], err = pipeline.Transform("clusters", m.mergeClustersForPipeline)
		if err != nil {
			return v, err
		}
	}

	return config.NewValue(out, v.Location()), nil
}

func (m *mergePipelineClusters) Apply(ctx context.Context, b *bundle.Bundle) error {
	return b.Config.Mutate(func(v config.Value) (config.Value, error) {
		if v == config.NilValue {
			return v, nil
		}

		nv, err := v.Transform("resources.pipelines", m.foreachPipeline)

		// It is not a problem if the pipelines key is not set.
		if config.IsNoSuchKeyError(err) {
			return v, nil
		}

		if err != nil {
			return v, err
		}

		return nv, nil
	})
}
