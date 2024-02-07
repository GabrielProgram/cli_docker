package bundle

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/databricks/cli/internal"
	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/service/jobs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAccBindJobToExistingJob(t *testing.T) {
	env := internal.GetEnvOrSkipTest(t, "CLOUD_ENV")
	t.Log(env)

	nodeTypeId := internal.GetNodeTypeId(env)
	uniqueId := uuid.New().String()
	bundleRoot, err := initTestTemplate(t, "basic", map[string]any{
		"unique_id":     uniqueId,
		"spark_version": "13.3.x-scala2.12",
		"node_type_id":  nodeTypeId,
	})
	require.NoError(t, err)

	jobId := createTestJob(t)
	t.Cleanup(func() {
		destroyJob(t, jobId)
		require.NoError(t, err)
	})

	t.Setenv("BUNDLE_ROOT", bundleRoot)
	c := internal.NewCobraTestRunner(t, "bundle", "deployment", "bind", "foo", fmt.Sprint(jobId), "--auto-approve")
	_, _, err = c.Run()
	require.NoError(t, err)

	// Remove .databricks directory to simulate a fresh deployment
	err = os.RemoveAll(filepath.Join(bundleRoot, ".databricks"))
	require.NoError(t, err)

	err = deployBundle(t, bundleRoot)
	require.NoError(t, err)

	w, err := databricks.NewWorkspaceClient()
	require.NoError(t, err)

	ctx := context.Background()
	// Check that job is bound and updated with config from bundle
	job, err := w.Jobs.Get(ctx, jobs.GetJobRequest{
		JobId: jobId,
	})
	require.NoError(t, err)
	require.Equal(t, job.Settings.Name, fmt.Sprintf("test-job-basic-%s", uniqueId))
	require.Contains(t, job.Settings.Tasks[0].SparkPythonTask.PythonFile, "hello_world.py")

	c = internal.NewCobraTestRunner(t, "bundle", "deployment", "unbind", "foo")
	_, _, err = c.Run()
	require.NoError(t, err)

	// Remove .databricks directory to simulate a fresh deployment
	err = os.RemoveAll(filepath.Join(bundleRoot, ".databricks"))
	require.NoError(t, err)

	err = destroyBundle(t, bundleRoot)
	require.NoError(t, err)

	// Check that job is unbound and exists after bundle is destroyed
	job, err = w.Jobs.Get(ctx, jobs.GetJobRequest{
		JobId: jobId,
	})
	require.NoError(t, err)
	require.Equal(t, job.Settings.Name, fmt.Sprintf("test-job-basic-%s", uniqueId))
	require.Contains(t, job.Settings.Tasks[0].SparkPythonTask.PythonFile, "hello_world.py")

}
