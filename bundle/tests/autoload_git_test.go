package config_tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitConfig(t *testing.T) {
	b := load(t, "./autoload_git")
	assert.Equal(t, "foo", b.Config.Bundle.Git.Branch)
	sshUrl := "git@github.com:databricks/bricks.git"
	httpsUrl := "https://github.com/databricks/bricks"
	assert.Contains(t, []string{sshUrl, httpsUrl}, b.Config.Bundle.Git.OriginURL)
}
