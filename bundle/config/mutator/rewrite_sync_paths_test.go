package mutator_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/config"
	"github.com/databricks/cli/bundle/config/mutator"
	"github.com/databricks/cli/bundle/internal/bundletest"
	"github.com/stretchr/testify/assert"
)

func TestRewriteSyncPathsRelative(t *testing.T) {
	b := &bundle.Bundle{
		Config: config.Root{
			Path: ".",
			Sync: config.Sync{
				Include: []string{
					"foo",
					"bar",
				},
				Exclude: []string{
					"baz",
					"qux",
				},
			},
		},
	}

	bundletest.SetLocation(b, "sync.include[0]", "./file.yml")
	bundletest.SetLocation(b, "sync.include[1]", "./a/file.yml")
	bundletest.SetLocation(b, "sync.exclude[0]", "./a/b/file.yml")
	bundletest.SetLocation(b, "sync.exclude[1]", "./a/b/c/file.yml")

	err := bundle.Apply(context.Background(), b, mutator.RewriteSyncPaths())
	assert.NoError(t, err)

	assert.Equal(t, filepath.Clean("foo"), b.Config.Sync.Include[0])
	assert.Equal(t, filepath.Clean("a/bar"), b.Config.Sync.Include[1])
	assert.Equal(t, filepath.Clean("a/b/baz"), b.Config.Sync.Exclude[0])
	assert.Equal(t, filepath.Clean("a/b/c/qux"), b.Config.Sync.Exclude[1])
}

func TestRewriteSyncPathsAbsolute(t *testing.T) {
	b := &bundle.Bundle{
		Config: config.Root{
			Path: "/tmp/dir",
			Sync: config.Sync{
				Include: []string{
					"foo",
					"bar",
				},
				Exclude: []string{
					"baz",
					"qux",
				},
			},
		},
	}

	bundletest.SetLocation(b, "sync.include[0]", "/tmp/dir/file.yml")
	bundletest.SetLocation(b, "sync.include[1]", "/tmp/dir/a/file.yml")
	bundletest.SetLocation(b, "sync.exclude[0]", "/tmp/dir/a/b/file.yml")
	bundletest.SetLocation(b, "sync.exclude[1]", "/tmp/dir/a/b/c/file.yml")

	err := bundle.Apply(context.Background(), b, mutator.RewriteSyncPaths())
	assert.NoError(t, err)

	assert.Equal(t, filepath.Clean("foo"), b.Config.Sync.Include[0])
	assert.Equal(t, filepath.Clean("a/bar"), b.Config.Sync.Include[1])
	assert.Equal(t, filepath.Clean("a/b/baz"), b.Config.Sync.Exclude[0])
	assert.Equal(t, filepath.Clean("a/b/c/qux"), b.Config.Sync.Exclude[1])
}

func TestRewriteSyncPathsErrorPaths(t *testing.T) {
	t.Run("no sync block", func(t *testing.T) {
		b := &bundle.Bundle{
			Config: config.Root{
				Path: ".",
			},
		}

		err := bundle.Apply(context.Background(), b, mutator.RewriteSyncPaths())
		assert.NoError(t, err)
	})

	t.Run("empty include/exclude blocks", func(t *testing.T) {
		b := &bundle.Bundle{
			Config: config.Root{
				Path: ".",
				Sync: config.Sync{
					Include: []string{},
					Exclude: []string{},
				},
			},
		}

		err := bundle.Apply(context.Background(), b, mutator.RewriteSyncPaths())
		assert.NoError(t, err)
	})
}
