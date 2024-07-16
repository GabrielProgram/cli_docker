package config_tests

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/libraries"
	"github.com/databricks/cli/bundle/phases"
	"github.com/stretchr/testify/require"
)

func TestPythonWheelBuild(t *testing.T) {
	ctx := context.Background()
	b, err := bundle.Load(ctx, "./python_wheel/python_wheel")
	require.NoError(t, err)

	diags := bundle.Apply(ctx, b, bundle.Seq(phases.Load(), phases.Build()))
	require.NoError(t, diags.Error())

	matches, err := filepath.Glob("./python_wheel/python_wheel/my_test_code/dist/my_test_code-*.whl")
	require.NoError(t, err)
	require.Equal(t, 1, len(matches))

	match := libraries.ValidateLocalLibrariesExist()
	diags = bundle.Apply(ctx, b, match)
	require.NoError(t, diags.Error())
}

func TestPythonWheelBuildAutoDetect(t *testing.T) {
	ctx := context.Background()
	b, err := bundle.Load(ctx, "./python_wheel/python_wheel_no_artifact")
	require.NoError(t, err)

	diags := bundle.Apply(ctx, b, bundle.Seq(phases.Load(), phases.Build()))
	require.NoError(t, diags.Error())

	matches, err := filepath.Glob("./python_wheel/python_wheel_no_artifact/dist/my_test_code-*.whl")
	require.NoError(t, err)
	require.Equal(t, 1, len(matches))

	match := libraries.ValidateLocalLibrariesExist()
	diags = bundle.Apply(ctx, b, match)
	require.NoError(t, diags.Error())
}

func TestPythonWheelWithDBFSLib(t *testing.T) {
	ctx := context.Background()
	b, err := bundle.Load(ctx, "./python_wheel/python_wheel_dbfs_lib")
	require.NoError(t, err)

	diags := bundle.Apply(ctx, b, bundle.Seq(phases.Load(), phases.Build()))
	require.NoError(t, diags.Error())

	match := libraries.ValidateLocalLibrariesExist()
	diags = bundle.Apply(ctx, b, match)
	require.NoError(t, diags.Error())
}

func TestPythonWheelBuildNoBuildJustUpload(t *testing.T) {
	ctx := context.Background()
	b, err := bundle.Load(ctx, "./python_wheel/python_wheel_no_artifact_no_setup")
	require.NoError(t, err)

	diags := bundle.Apply(ctx, b, bundle.Seq(phases.Load(), phases.Build()))
	require.NoError(t, diags.Error())

	match := libraries.ValidateLocalLibrariesExist()
	diags = bundle.Apply(ctx, b, match)
	require.ErrorContains(t, diags.Error(), "./non-existing/*.whl")

	require.NotZero(t, len(b.Config.Artifacts))

	artifact := b.Config.Artifacts["my_test_code-0.0.1-py3-none-any.whl"]
	require.NotNil(t, artifact)
	require.Empty(t, artifact.BuildCommand)
	require.Contains(t, artifact.Files[0].Source, filepath.Join(b.RootPath, "package",
		"my_test_code-0.0.1-py3-none-any.whl",
	))
}

func TestPythonWheelBuildWithEnvironmentKey(t *testing.T) {
	ctx := context.Background()
	b, err := bundle.Load(ctx, "./python_wheel/environment_key")
	require.NoError(t, err)

	diags := bundle.Apply(ctx, b, bundle.Seq(phases.Load(), phases.Build()))
	require.NoError(t, diags.Error())

	matches, err := filepath.Glob("./python_wheel/environment_key/my_test_code/dist/my_test_code-*.whl")
	require.NoError(t, err)
	require.Equal(t, 1, len(matches))

	match := libraries.ValidateLocalLibrariesExist()
	diags = bundle.Apply(ctx, b, match)
	require.NoError(t, diags.Error())
}
