package git

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecusiveListFile(t *testing.T) {
	// Create .gitignore and ignore .gitignore and any files in
	// ignored_dir
	projectDir := t.TempDir()
	f3, err := os.Create(filepath.Join(projectDir, ".gitignore"))
	assert.NoError(t, err)
	defer f3.Close()
	f3.WriteString(".gitignore\nignored_dir")

	// create config file
	f4, err := os.Create(filepath.Join(projectDir, "databricks.yml"))
	assert.NoError(t, err)
	defer f4.Close()

	// config file is returned
	// .gitignore is not because we explictly ignore it in .gitignore
	fileSet := NewFileSet(projectDir)
	files, err := fileSet.RecursiveListFiles(projectDir)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, "databricks.yml", files[0].Relative)

	helloTxtRelativePath := filepath.Join("a/b/c", "hello.txt")

	// Check that newly added files not in .gitignore
	// are being tracked
	dir1 := filepath.Join(projectDir, "a", "b", "c")
	dir2 := filepath.Join(projectDir, "ignored_dir", "e")
	err = os.MkdirAll(dir2, 0o755)
	assert.NoError(t, err)
	err = os.MkdirAll(dir1, 0o755)
	assert.NoError(t, err)
	f1, err := os.Create(filepath.Join(projectDir, helloTxtRelativePath))
	assert.NoError(t, err)
	defer f1.Close()
	f2, err := os.Create(filepath.Join(projectDir, "ignored_dir/e/world.txt"))
	assert.NoError(t, err)
	defer f2.Close()

	files, err = fileSet.RecursiveListFiles(projectDir)
	assert.NoError(t, err)
	assert.Len(t, files, 2)
	var fileNames []string
	for _, v := range files {
		fileNames = append(fileNames, v.Relative)
	}
	assert.Contains(t, fileNames, "databricks.yml")
	assert.Contains(t, fileNames, helloTxtRelativePath)
}

func TestFileSetNonCleanRoot(t *testing.T) {
	// Test what happens if the root directory can be simplified.
	// Path simplification is done by most filepath functions.

	// remove this once equivalent tests for windows have been set up
	// or this test has been fixed for windows
	// date: 28 Nov 2022
	if runtime.GOOS == "windows" {
		t.Skip("skipping temperorilty to make windows unit tests green")
	}

	root := "./../git"
	require.NotEqual(t, root, filepath.Clean(root))
	fs := NewFileSet(root)
	files, err := fs.All()
	require.NoError(t, err)

	found := false
	for _, f := range files {
		if strings.Contains(f.Relative, "fileset_test") {
			assert.Equal(t, "fileset_test.go", f.Relative)
			assert.Equal(t, "../git/fileset_test.go", f.Absolute)
			found = true
		}
	}

	assert.True(t, found)
}
