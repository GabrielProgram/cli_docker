package root

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/internal/testutil"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDatabricksCfg(t *testing.T) {
	tempHomeDir := t.TempDir()
	homeEnvVar := "HOME"
	if runtime.GOOS == "windows" {
		homeEnvVar = "USERPROFILE"
	}

	cfg := []byte("[PROFILE-1]\nhost = https://a.com\ntoken = a\n[PROFILE-2]\nhost = https://a.com\ntoken = b\n")
	err := os.WriteFile(filepath.Join(tempHomeDir, ".databrickscfg"), cfg, 0644)
	assert.NoError(t, err)

	t.Setenv("DATABRICKS_CONFIG_FILE", "")
	t.Setenv(homeEnvVar, tempHomeDir)
}

func emptyCommand(t *testing.T) *cobra.Command {
	ctx := context.Background()
	cmd := &cobra.Command{}
	cmd.SetContext(ctx)
	initProfileFlag(cmd)
	return cmd
}

func setup(t *testing.T, cmd *cobra.Command, host string) *bundle.Bundle {
	setupDatabricksCfg(t)

	rootPath := t.TempDir()
	testutil.Chdir(t, rootPath)

	contents := fmt.Sprintf(`
workspace:
  host: %q
`, host)
	err := os.WriteFile(filepath.Join(rootPath, "databricks.yml"), []byte(contents), 0644)
	require.NoError(t, err)

	b, diags := MustConfigureBundle(cmd)
	require.NoError(t, diags.Error())
	return b
}

func TestBundleConfigureDefault(t *testing.T) {
	testutil.CleanupEnvironment(t)

	cmd := emptyCommand(t)
	b := setup(t, cmd, "https://x.com")

	client, err := b.InitializeWorkspaceClient()
	require.NoError(t, err)
	assert.Equal(t, "https://x.com", client.Config.Host)
}

func TestBundleConfigureWithMultipleMatches(t *testing.T) {
	testutil.CleanupEnvironment(t)

	cmd := emptyCommand(t)
	b := setup(t, cmd, "https://a.com")

	_, err := b.InitializeWorkspaceClient()
	assert.ErrorContains(t, err, "multiple profiles matched: PROFILE-1, PROFILE-2")
}

func TestBundleConfigureWithNonExistentProfileFlag(t *testing.T) {
	testutil.CleanupEnvironment(t)

	cmd := emptyCommand(t)
	cmd.Flag("profile").Value.Set("NOEXIST")
	b := setup(t, cmd, "https://x.com")

	_, err := b.InitializeWorkspaceClient()
	assert.ErrorContains(t, err, "has no NOEXIST profile configured")
}

func TestBundleConfigureWithMismatchedProfile(t *testing.T) {
	testutil.CleanupEnvironment(t)

	cmd := emptyCommand(t)
	cmd.Flag("profile").Value.Set("PROFILE-1")
	b := setup(t, cmd, "https://x.com")

	_, err := b.InitializeWorkspaceClient()
	assert.ErrorContains(t, err, "config host mismatch: profile uses host https://a.com, but CLI configured to use https://x.com")
}

func TestBundleConfigureWithCorrectProfile(t *testing.T) {
	testutil.CleanupEnvironment(t)

	cmd := emptyCommand(t)
	cmd.Flag("profile").Value.Set("PROFILE-1")
	b := setup(t, cmd, "https://a.com")

	client, err := b.InitializeWorkspaceClient()
	require.NoError(t, err)
	assert.Equal(t, "https://a.com", client.Config.Host)
	assert.Equal(t, "PROFILE-1", client.Config.Profile)
}

func TestBundleConfigureWithMismatchedProfileEnvVariable(t *testing.T) {
	testutil.CleanupEnvironment(t)

	t.Setenv("DATABRICKS_CONFIG_PROFILE", "PROFILE-1")
	cmd := emptyCommand(t)
	b := setup(t, cmd, "https://x.com")

	_, err := b.InitializeWorkspaceClient()
	assert.ErrorContains(t, err, "config host mismatch: profile uses host https://a.com, but CLI configured to use https://x.com")
}

func TestBundleConfigureWithProfileFlagAndEnvVariable(t *testing.T) {
	testutil.CleanupEnvironment(t)

	t.Setenv("DATABRICKS_CONFIG_PROFILE", "NOEXIST")
	cmd := emptyCommand(t)
	cmd.Flag("profile").Value.Set("PROFILE-1")
	b := setup(t, cmd, "https://a.com")

	client, err := b.InitializeWorkspaceClient()
	require.NoError(t, err)
	assert.Equal(t, "https://a.com", client.Config.Host)
	assert.Equal(t, "PROFILE-1", client.Config.Profile)
}

func TestTargetFlagFull(t *testing.T) {
	cmd := emptyCommand(t)
	initTargetFlag(cmd)
	cmd.SetArgs([]string{"version", "--target", "development"})

	ctx := context.Background()
	err := cmd.ExecuteContext(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "development", getTarget(cmd))
}

func TestTargetFlagShort(t *testing.T) {
	cmd := emptyCommand(t)
	initTargetFlag(cmd)
	cmd.SetArgs([]string{"version", "-t", "production"})

	ctx := context.Background()
	err := cmd.ExecuteContext(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "production", getTarget(cmd))
}

// TODO: remove when environment flag is fully deprecated
func TestTargetEnvironmentFlag(t *testing.T) {
	cmd := emptyCommand(t)
	initTargetFlag(cmd)
	initEnvironmentFlag(cmd)
	cmd.SetArgs([]string{"version", "--environment", "development"})

	ctx := context.Background()
	err := cmd.ExecuteContext(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "development", getTarget(cmd))
}
