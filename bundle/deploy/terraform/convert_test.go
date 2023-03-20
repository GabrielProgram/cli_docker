package terraform

import (
	"testing"

	"github.com/databricks/bricks/bundle/config"
	"github.com/databricks/bricks/bundle/config/resources"
	"github.com/databricks/databricks-sdk-go/service/clusters"
	"github.com/databricks/databricks-sdk-go/service/jobs"
	"github.com/databricks/databricks-sdk-go/service/libraries"
	"github.com/databricks/databricks-sdk-go/service/mlflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertJob(t *testing.T) {
	var src = resources.Job{
		JobSettings: &jobs.JobSettings{
			Name: "my job",
			JobClusters: []jobs.JobCluster{
				{
					JobClusterKey: "key",
					NewCluster: &clusters.BaseClusterInfo{
						SparkVersion: "10.4.x-scala2.12",
					},
				},
			},
			GitSource: &jobs.GitSource{
				GitProvider: jobs.GitSourceGitProviderGithub,
				GitUrl:      "https://github.com/foo/bar",
			},
		},
	}

	var config = config.Root{
		Resources: config.Resources{
			Jobs: map[string]*resources.Job{
				"my_job": &src,
			},
		},
	}

	out := BundleToTerraform(&config)
	assert.Equal(t, "my job", out.Resource.Job["my_job"].Name)
	assert.Len(t, out.Resource.Job["my_job"].JobCluster, 1)
	assert.Equal(t, "https://github.com/foo/bar", out.Resource.Job["my_job"].GitSource.Url)
	assert.Nil(t, out.Data)
}

func TestConvertJobTaskLibraries(t *testing.T) {
	var src = resources.Job{
		JobSettings: &jobs.JobSettings{
			Name: "my job",
			Tasks: []jobs.JobTaskSettings{
				{
					TaskKey: "key",
					Libraries: []libraries.Library{
						{
							Pypi: &libraries.PythonPyPiLibrary{
								Package: "mlflow",
							},
						},
					},
				},
			},
		},
	}

	var config = config.Root{
		Resources: config.Resources{
			Jobs: map[string]*resources.Job{
				"my_job": &src,
			},
		},
	}

	out := BundleToTerraform(&config)
	assert.Equal(t, "my job", out.Resource.Job["my_job"].Name)
	require.Len(t, out.Resource.Job["my_job"].Task, 1)
	require.Len(t, out.Resource.Job["my_job"].Task[0].Library, 1)
	assert.Equal(t, "mlflow", out.Resource.Job["my_job"].Task[0].Library[0].Pypi.Package)
}

func TestConvertModel(t *testing.T) {
	var src = resources.MlflowModel{
		RegisteredModel: &mlflow.RegisteredModel{
			Name:        "name",
			Description: "description",
			Tags: []mlflow.RegisteredModelTag{
				{
					Key:   "k1",
					Value: "v1",
				},
				{
					Key:   "k2",
					Value: "v2",
				},
			},
		},
	}

	var config = config.Root{
		Resources: config.Resources{
			Models: map[string]*resources.MlflowModel{
				"my_model": &src,
			},
		},
	}

	out := BundleToTerraform(&config)
	assert.Equal(t, "name", out.Resource.MlflowModel["my_model"].Name)
	assert.Equal(t, "description", out.Resource.MlflowModel["my_model"].Description)
	assert.Len(t, out.Resource.MlflowModel["my_model"].Tags, 2)
	assert.Equal(t, "k1", out.Resource.MlflowModel["my_model"].Tags[0].Key)
	assert.Equal(t, "v1", out.Resource.MlflowModel["my_model"].Tags[0].Value)
	assert.Equal(t, "k2", out.Resource.MlflowModel["my_model"].Tags[1].Key)
	assert.Equal(t, "v2", out.Resource.MlflowModel["my_model"].Tags[1].Value)
	assert.Nil(t, out.Data)
}

func TestConvertExperiment(t *testing.T) {
	var src = resources.MlflowExperiment{
		Experiment: &mlflow.Experiment{
			Name: "name",
		},
	}

	var config = config.Root{
		Resources: config.Resources{
			Experiments: map[string]*resources.MlflowExperiment{
				"my_experiment": &src,
			},
		},
	}

	out := BundleToTerraform(&config)
	assert.Equal(t, "name", out.Resource.MlflowExperiment["my_experiment"].Name)
	assert.Nil(t, out.Data)
}
