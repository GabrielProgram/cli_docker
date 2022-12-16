package mutator_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/databricks/bricks/bundle"
	"github.com/databricks/bricks/bundle/config"
	"github.com/databricks/bricks/bundle/config/mutator"
	"github.com/databricks/bricks/bundle/config/resources"
	"github.com/databricks/databricks-sdk-go/service/jobs"
	"github.com/databricks/databricks-sdk-go/service/pipelines"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func touchFile(t *testing.T, path string) {
	f, err := os.Create(path)
	require.NoError(t, err)
	f.Close()
}

func TestNotebookPaths(t *testing.T) {
	dir := t.TempDir()
	touchFile(t, filepath.Join(dir, "my_job_notebook.py"))
	touchFile(t, filepath.Join(dir, "my_pipeline_notebook.py"))

	bundle := &bundle.Bundle{
		Config: config.Root{
			Path: dir,
			Resources: config.Resources{
				Jobs: map[string]*resources.Job{
					"job": {
						JobSettings: &jobs.JobSettings{
							Tasks: []jobs.JobTaskSettings{
								{
									NotebookTask: &jobs.NotebookTask{
										NotebookPath: "./my_job_notebook.py",
									},
								},
								{
									NotebookTask: &jobs.NotebookTask{
										NotebookPath: "./doesnt_exist.py",
									},
								},
								{
									NotebookTask: &jobs.NotebookTask{
										NotebookPath: "./my_job_notebook.py",
									},
								},
								{
									PythonWheelTask: &jobs.PythonWheelTask{
										PackageName: "foo",
									},
								},
							},
						},
					},
				},
				Pipelines: map[string]*resources.Pipeline{
					"pipeline": {
						PipelineSpec: &pipelines.PipelineSpec{
							Libraries: []pipelines.PipelineLibrary{
								{
									Notebook: &pipelines.NotebookLibrary{
										Path: "./my_pipeline_notebook.py",
									},
								},
								{
									Notebook: &pipelines.NotebookLibrary{
										Path: "./doesnt_exist.py",
									},
								},
								{
									Notebook: &pipelines.NotebookLibrary{
										Path: "./my_pipeline_notebook.py",
									},
								},
								{
									Jar: "foo",
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := mutator.TranslateNotebookPaths().Apply(context.Background(), bundle)
	require.NoError(t, err)

	// Assert that the notebook artifact was defined.
	assert.Len(t, bundle.Config.Artifacts, 2)
	for _, artifact := range bundle.Config.Artifacts {
		assert.Contains(t, artifact.Notebook.Path, "notebook.py")
	}

	// Assert that the path in the tasks now refer to the artifact.
	assert.Equal(
		t,
		"${artifacts.my_job_notebook_py.notebook.remote_path}",
		bundle.Config.Resources.Jobs["job"].Tasks[0].NotebookTask.NotebookPath,
	)
	assert.Equal(
		t,
		"./doesnt_exist.py",
		bundle.Config.Resources.Jobs["job"].Tasks[1].NotebookTask.NotebookPath,
	)
	assert.Equal(
		t,
		"${artifacts.my_job_notebook_py.notebook.remote_path}",
		bundle.Config.Resources.Jobs["job"].Tasks[2].NotebookTask.NotebookPath,
	)

	// Assert that the path in the libraries now refer to the artifact.
	assert.Equal(
		t,
		"${artifacts.my_pipeline_notebook_py.notebook.remote_path}",
		bundle.Config.Resources.Pipelines["pipeline"].Libraries[0].Notebook.Path,
	)
	assert.Equal(
		t,
		"./doesnt_exist.py",
		bundle.Config.Resources.Pipelines["pipeline"].Libraries[1].Notebook.Path,
	)
	assert.Equal(
		t,
		"${artifacts.my_pipeline_notebook_py.notebook.remote_path}",
		bundle.Config.Resources.Pipelines["pipeline"].Libraries[2].Notebook.Path,
	)
}
