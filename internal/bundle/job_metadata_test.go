package bundle

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strconv"
	"testing"

	"github.com/databricks/cli/bundle/config"
	"github.com/databricks/cli/bundle/metadata"
	"github.com/databricks/cli/internal"
	"github.com/databricks/cli/libs/filer"
	"github.com/databricks/databricks-sdk-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccJobsMetadataFile(t *testing.T) {
	env := internal.GetEnvOrSkipTest(t, "CLOUD_ENV")
	t.Log(env)

	w, err := databricks.NewWorkspaceClient()
	require.NoError(t, err)

	nodeTypeId := internal.GetNodeTypeId(env)
	uniqueId := uuid.New().String()
	bundleRoot, err := initTestTemplate(t, "job_metadata", map[string]any{
		"unique_id":     uniqueId,
		"node_type_id":  nodeTypeId,
		"spark_version": "13.2.x-snapshot-scala2.12",
	})
	require.NoError(t, err)

	// deploy bundle
	err = deployBundle(t, bundleRoot)
	require.NoError(t, err)

	// Cleanup the deployed bundle
	t.Cleanup(func() {
		err = destroyBundle(t, bundleRoot)
		require.NoError(t, err)
	})

	// assert job 1 is created
	jobName := "test-job-metadata-1-" + uniqueId
	job1, err := w.Jobs.GetBySettingsName(context.Background(), jobName)
	require.NoError(t, err)
	assert.Equal(t, job1.Settings.Name, jobName)

	// assert job 2 is created
	jobName = "test-job-metadata-2-" + uniqueId
	job2, err := w.Jobs.GetBySettingsName(context.Background(), jobName)
	require.NoError(t, err)
	assert.Equal(t, job2.Settings.Name, jobName)

	// Compute root path for the bundle deployment
	me, err := w.CurrentUser.Me(context.Background())
	require.NoError(t, err)
	root := fmt.Sprintf("/Users/%s/.bundle/%s", me.UserName, uniqueId)
	f, err := filer.NewWorkspaceFilesClient(w, root)
	require.NoError(t, err)

	// Read metadata object from the workspace
	r, err := f.Read(context.Background(), "state/metadata.json")
	require.NoError(t, err)
	b, err := io.ReadAll(r)
	require.NoError(t, err)
	actualMetadata := metadata.Metadata{}
	err = json.Unmarshal(b, &actualMetadata)
	require.NoError(t, err)

	// expected value for the metadata
	expectedMetadata := metadata.Metadata{
		Version: metadata.Version,
		Config: metadata.Config{
			Bundle: metadata.Bundle{
				Git: config.Git{
					BundleRootPath: ".",
				},
			},
			Workspace: metadata.Workspace{
				FilePath: path.Join(root, "files"),
			},
			Resources: metadata.Resources{
				Jobs: map[string]*metadata.Job{
					"foo": {
						ID:           strconv.FormatInt(job1.JobId, 10),
						RelativePath: "databricks.yml",
					},
					"bar": {
						ID:           strconv.FormatInt(job2.JobId, 10),
						RelativePath: "a/b/resources.yml",
					},
				},
			},
		},
	}

	// Assert metadata matches what we expected.
	assert.Equal(t, expectedMetadata, actualMetadata)
}
