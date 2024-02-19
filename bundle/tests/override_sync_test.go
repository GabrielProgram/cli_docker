package config_tests

import (
	"path/filepath"
	"testing"

	"github.com/databricks/cli/bundle"
	"github.com/stretchr/testify/assert"
)

func TestOverrideSyncTarget(t *testing.T) {
	var b *bundle.Bundle

	b = loadTarget(t, "./override_sync", "development")
	assert.ElementsMatch(t, []string{filepath.FromSlash("src/*"), filepath.FromSlash("tests/*")}, b.Config.Sync.Include)
	assert.ElementsMatch(t, []string{filepath.FromSlash("dist")}, b.Config.Sync.Exclude)

	b = loadTarget(t, "./override_sync", "staging")
	assert.ElementsMatch(t, []string{filepath.FromSlash("src/*"), filepath.FromSlash("fixtures/*")}, b.Config.Sync.Include)
	assert.ElementsMatch(t, []string{}, b.Config.Sync.Exclude)

	b = loadTarget(t, "./override_sync", "prod")
	assert.ElementsMatch(t, []string{filepath.FromSlash("src/*")}, b.Config.Sync.Include)
	assert.ElementsMatch(t, []string{}, b.Config.Sync.Exclude)
}

func TestOverrideSyncTargetNoRootSync(t *testing.T) {
	var b *bundle.Bundle

	b = loadTarget(t, "./override_sync_no_root", "development")
	assert.ElementsMatch(t, []string{filepath.FromSlash("tests/*")}, b.Config.Sync.Include)
	assert.ElementsMatch(t, []string{filepath.FromSlash("dist")}, b.Config.Sync.Exclude)

	b = loadTarget(t, "./override_sync_no_root", "staging")
	assert.ElementsMatch(t, []string{filepath.FromSlash("fixtures/*")}, b.Config.Sync.Include)
	assert.ElementsMatch(t, []string{}, b.Config.Sync.Exclude)

	b = loadTarget(t, "./override_sync_no_root", "prod")
	assert.ElementsMatch(t, []string{}, b.Config.Sync.Include)
	assert.ElementsMatch(t, []string{}, b.Config.Sync.Exclude)
}
