package config_tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentOverridesDev(t *testing.T) {
	development := loadEnvironment(t, "./environment_overrides", "development")
	assert.Equal(t, "https://development.acme.cloud.databricks.com/", development.Workspace.Host)
	staging := loadEnvironment(t, "./environment_overrides", "staging")
	assert.Equal(t, "https://staging.acme.cloud.databricks.com/", staging.Workspace.Host)
}
