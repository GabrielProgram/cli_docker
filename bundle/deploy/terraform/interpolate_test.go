package terraform

import (
	"context"
	"testing"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/config"
	"github.com/databricks/cli/bundle/config/resources"
	"github.com/databricks/databricks-sdk-go/service/jobs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterpolate(t *testing.T) {
	b := &bundle.Bundle{
		Config: config.Root{
			Bundle: config.Bundle{
				Name: "example",
			},
			Resources: config.Resources{
				Jobs: map[string]*resources.Job{
					"my_job": {
						JobSettings: &jobs.JobSettings{
							Tags: map[string]string{
								"other_pipeline":         "${resources.pipelines.other_pipeline.id}",
								"other_job":              "${resources.jobs.other_job.id}",
								"other_model":            "${resources.models.other_model.id}",
								"other_experiment":       "${resources.experiments.other_experiment.id}",
								"other_model_serving":    "${resources.model_serving_endpoints.other_model_serving.id}",
								"other_registered_model": "${resources.registered_models.other_registered_model.id}",
							},
						},
					},
				},
			},
		},
	}

	err := bundle.Apply(context.Background(), b, Interpolate())
	require.NoError(t, err)

	j := b.Config.Resources.Jobs["my_job"]
	assert.Equal(t, "${databricks_pipeline.other_pipeline.id}", j.Tags["other_pipeline"])
	assert.Equal(t, "${databricks_job.other_job.id}", j.Tags["other_job"])
	assert.Equal(t, "${databricks_mlflow_model.other_model.id}", j.Tags["other_model"])
	assert.Equal(t, "${databricks_mlflow_experiment.other_experiment.id}", j.Tags["other_experiment"])
	assert.Equal(t, "${databricks_model_serving.other_model_serving.id}", j.Tags["other_model_serving"])
	assert.Equal(t, "${databricks_registered_model.other_registered_model.id}", j.Tags["other_registered_model"])
}

func TestInterpolateUnknownResourceType(t *testing.T) {
	b := &bundle.Bundle{
		Config: config.Root{
			Resources: config.Resources{
				Jobs: map[string]*resources.Job{
					"my_job": {
						JobSettings: &jobs.JobSettings{
							Tags: map[string]string{
								"other_unknown": "${resources.unknown.other_unknown.id}",
							},
						},
					},
				},
			},
		},
	}

	err := bundle.Apply(context.Background(), b, Interpolate())
	assert.Contains(t, err.Error(), `reference does not exist: ${resources.unknown.other_unknown.id}`)
}
