package mutator

import (
	"context"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/libs/dyn"
	"github.com/databricks/cli/libs/dyn/merge"
)

type mergeJobClusters struct{}

func MergeJobClusters() bundle.Mutator {
	return &mergeJobClusters{}
}

func (m *mergeJobClusters) Name() string {
	return "MergeJobClusters"
}

func (m *mergeJobClusters) jobClusterKey(v dyn.Value) string {
	switch v.Kind() {
	case dyn.KindNil:
		return ""
	case dyn.KindString:
		return v.MustString()
	default:
		panic("job cluster key must be a string")
	}
}

func (m *mergeJobClusters) Apply(ctx context.Context, b *bundle.Bundle) error {
	return b.Config.Mutate(func(v dyn.Value) (dyn.Value, error) {
		if v == dyn.NilValue {
			return v, nil
		}

		return dyn.Map(v, "resources.jobs", dyn.Foreach(func(job dyn.Value) (dyn.Value, error) {
			return dyn.Map(job, "job_clusters", merge.ElementsByKey("job_cluster_key", m.jobClusterKey))
		}))
	})
}
