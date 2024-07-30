package bundle

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databricks/cli/internal/acc"
	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/apierr"
	"github.com/databricks/databricks-sdk-go/service/catalog"
	"github.com/databricks/databricks-sdk-go/service/files"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUcSchemaBundle(t *testing.T, ctx context.Context, w *databricks.WorkspaceClient, uniqueId string) string {
	bundleRoot, err := initTestTemplate(t, ctx, "uc_schema", map[string]any{
		"unique_id": uniqueId,
	})
	require.NoError(t, err)

	err = deployBundle(t, ctx, bundleRoot)
	require.NoError(t, err)

	t.Cleanup(func() {
		destroyBundle(t, ctx, bundleRoot)
	})

	// Assert the schema is created
	catalogName := "main"
	schemaName := "test-schema-" + uniqueId
	schema, err := w.Schemas.GetByFullName(ctx, strings.Join([]string{catalogName, schemaName}, "."))
	require.NoError(t, err)
	require.Equal(t, strings.Join([]string{catalogName, schemaName}, "."), schema.FullName)
	require.Equal(t, "This schema was created from DABs", schema.Comment)

	// Assert the pipeline is created
	pipelineName := "test-pipeline-" + uniqueId
	pipeline, err := w.Pipelines.GetByName(ctx, pipelineName)
	require.NoError(t, err)
	require.Equal(t, pipelineName, pipeline.Name)
	id := pipeline.PipelineId

	// Assert the pipeline uses the schema
	i, err := w.Pipelines.GetByPipelineId(ctx, id)
	require.NoError(t, err)
	require.Equal(t, catalogName, i.Spec.Catalog)
	require.Equal(t, strings.Join([]string{catalogName, schemaName}, "."), i.Spec.Target)

	// Create a volume in the schema, and add a file to it. This ensure that the
	// schema as some data in it and deletion will fail unless the generated
	// terraform configuration has force_destroy set to true.
	volumeName := "test-volume-" + uniqueId
	volume, err := w.Volumes.Create(ctx, catalog.CreateVolumeRequestContent{
		CatalogName: catalogName,
		SchemaName:  schemaName,
		Name:        volumeName,
		VolumeType:  catalog.VolumeTypeManaged,
	})
	require.NoError(t, err)
	require.Equal(t, volume.Name, volumeName)

	fileName := "test-file-" + uniqueId
	err = w.Files.Upload(ctx, files.UploadRequest{
		Contents: io.NopCloser(strings.NewReader("Hello, world!")),
		FilePath: fmt.Sprintf("/Volumes/%s/%s/%s/%s", catalogName, schemaName, volumeName, fileName),
	})
	require.NoError(t, err)

	return bundleRoot
}

func TestAccBundleDeployUcSchema(t *testing.T) {
	ctx, wt := acc.UcWorkspaceTest(t)
	w := wt.W

	uniqueId := uuid.New().String()
	schemaName := "test-schema-" + uniqueId
	catalogName := "main"

	bundleRoot := setupUcSchemaBundle(t, ctx, w, uniqueId)

	// Remove the UC schema from the resource configuration.
	err := os.Remove(filepath.Join(bundleRoot, "schema.yml"))
	require.NoError(t, err)

	// Redeploy the bundle
	err = deployBundle(t, ctx, bundleRoot)
	require.NoError(t, err)

	// Assert the schema is deleted
	_, err = w.Schemas.GetByFullName(ctx, strings.Join([]string{catalogName, schemaName}, "."))
	apiErr := &apierr.APIError{}
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, "SCHEMA_DOES_NOT_EXIST", apiErr.ErrorCode)
}
