package template

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertFileContent(t *testing.T, path string, content string) {
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, content, string(b))
}

func assertFilePermissions(t *testing.T, path string, perm fs.FileMode) {
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, perm, info.Mode().Perm())
}

func TestRendererWithAssociatedTemplateInLibrary(t *testing.T) {
	tmpDir := t.TempDir()

	r, err := newRenderer(context.Background(), nil, "./testdata/email/template", "./testdata/email/library", tmpDir)
	require.NoError(t, err)

	err = r.walk()
	require.NoError(t, err)

	err = r.persistToDisk()
	require.NoError(t, err)

	b, err := os.ReadFile(filepath.Join(tmpDir, "my_email"))
	require.NoError(t, err)
	assert.Equal(t, "shreyas.goenka@databricks.com", strings.Trim(string(b), "\n\r"))
}

func TestRendererExecuteTemplate(t *testing.T) {
	templateText :=
		`"{{.count}} items are made of {{.Material}}".
{{if eq .Animal "sheep" }}
Sheep wool is the best!
{{else}}
{{.Animal}} wool is not too bad...
{{end}}
My email is {{template "email"}}
`

	r := renderer{
		config: map[string]any{
			"Material": "wool",
			"count":    1,
			"Animal":   "sheep",
		},
		baseTemplate: template.Must(template.New("base").Parse(`{{define "email"}}shreyas.goenka@databricks.com{{end}}`)),
	}

	statement, err := r.executeTemplate(templateText)
	require.NoError(t, err)
	assert.Contains(t, statement, `"1 items are made of wool"`)
	assert.NotContains(t, statement, `cat wool is not too bad.."`)
	assert.Contains(t, statement, "Sheep wool is the best!")
	assert.Contains(t, statement, `My email is shreyas.goenka@databricks.com`)

	r = renderer{
		config: map[string]any{
			"Material": "wool",
			"count":    1,
			"Animal":   "cat",
		},
		baseTemplate: template.Must(template.New("base").Parse(`{{define "email"}}hrithik.roshan@databricks.com{{end}}`)),
	}

	statement, err = r.executeTemplate(templateText)
	require.NoError(t, err)
	assert.Contains(t, statement, `"1 items are made of wool"`)
	assert.Contains(t, statement, `cat wool is not too bad...`)
	assert.NotContains(t, statement, "Sheep wool is the best!")
	assert.Contains(t, statement, `My email is hrithik.roshan@databricks.com`)
}

func TestRendererIsSkipped(t *testing.T) {
	r := renderer{
		skipPatterns: []string{"a*", "*yz", "def", "a/b/*"},
	}

	// skipped paths
	isSkipped, err := r.isSkipped("abc")
	require.NoError(t, err)
	assert.True(t, isSkipped)

	isSkipped, err = r.isSkipped("abcd")
	require.NoError(t, err)
	assert.True(t, isSkipped)

	isSkipped, err = r.isSkipped("a")
	require.NoError(t, err)
	assert.True(t, isSkipped)

	isSkipped, err = r.isSkipped("xxyz")
	require.NoError(t, err)
	assert.True(t, isSkipped)

	isSkipped, err = r.isSkipped("yz")
	require.NoError(t, err)
	assert.True(t, isSkipped)

	isSkipped, err = r.isSkipped("a/b/c")
	require.NoError(t, err)
	assert.True(t, isSkipped)

	// NOT skipped paths
	isSkipped, err = r.isSkipped(".")
	require.NoError(t, err)
	assert.False(t, isSkipped)

	isSkipped, err = r.isSkipped("y")
	require.NoError(t, err)
	assert.False(t, isSkipped)

	isSkipped, err = r.isSkipped("z")
	require.NoError(t, err)
	assert.False(t, isSkipped)

	isSkipped, err = r.isSkipped("defg")
	require.NoError(t, err)
	assert.False(t, isSkipped)

	isSkipped, err = r.isSkipped("cat")
	require.NoError(t, err)
	assert.False(t, isSkipped)

	isSkipped, err = r.isSkipped("a/b/c/d")
	require.NoError(t, err)
	assert.False(t, isSkipped)
}

func TestRendererPersistToDisk(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	r := &renderer{
		ctx:          ctx,
		instanceRoot: tmpDir,
		skipPatterns: []string{"a/b/c", "mn*"},
		files: []*inMemoryFile{
			{
				root:    tmpDir,
				relPath: "a/b/c",
				content: nil,
				perm:    0444,
			},
			{
				root:    tmpDir,
				relPath: "mno",
				content: nil,
				perm:    0444,
			},
			{
				root:    tmpDir,
				relPath: "a/b/d",
				content: []byte("123"),
				perm:    0444,
			},
			{
				root:    tmpDir,
				relPath: "mmnn",
				content: []byte("456"),
				perm:    0444,
			},
		},
	}

	err := r.persistToDisk()
	require.NoError(t, err)

	assert.NoFileExists(t, filepath.Join(tmpDir, "a", "b", "c"))
	assert.NoFileExists(t, filepath.Join(tmpDir, "mno"))

	assertFileContent(t, filepath.Join(tmpDir, "a", "b", "d"), "123")
	assertFilePermissions(t, filepath.Join(tmpDir, "a", "b", "d"), 0444)
	assertFileContent(t, filepath.Join(tmpDir, "mmnn"), "456")
	assertFilePermissions(t, filepath.Join(tmpDir, "mmnn"), 0444)
}

func TestRendererWalk(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	r, err := newRenderer(ctx, nil, "./testdata/walk/template", "./testdata/walk/library", tmpDir)
	require.NoError(t, err)

	err = r.walk()
	assert.NoError(t, err)

	getContent := func(r *renderer, path string) string {
		for _, f := range r.files {
			if f.relPath == path {
				return strings.Trim(string(f.content), "\r\n")
			}
		}
		require.FailNow(t, "file is absent: "+path)
		return ""
	}

	assert.Len(t, r.files, 4)
	assert.Equal(t, "file one", getContent(r, "file1"))
	assert.Equal(t, "file two", getContent(r, "file2"))
	assert.Equal(t, "file three", getContent(r, "dir1/dir3/file3"))
	assert.Equal(t, "file four", getContent(r, "dir2/file4"))
}

func TestRendererFailFunction(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	r, err := newRenderer(ctx, nil, "./testdata/fail/template", "./testdata/fail/library", tmpDir)
	require.NoError(t, err)

	err = r.walk()
	assert.Equal(t, "I am an error message", err.Error())
}

func TestRendererSkipsDirsEagerly(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	r, err := newRenderer(ctx, nil, "./testdata/skip-dir-eagerly/template", "./testdata/skip-dir-eagerly/library", tmpDir)
	require.NoError(t, err)

	err = r.walk()
	assert.NoError(t, err)

	assert.Len(t, r.files, 1)
	content := string(r.files[0].content)
	assert.Equal(t, "I should be the only file created", strings.Trim(content, "\r\n"))
}

func TestRendererSkipAllFilesInCurrentDirectory(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	r, err := newRenderer(ctx, nil, "./testdata/skip-all-files-in-cwd/template", "./testdata/skip-all-files-in-cwd/library", tmpDir)
	require.NoError(t, err)

	err = r.walk()
	assert.NoError(t, err)
	// All 3 files are executed and have in memory representations
	require.Len(t, r.files, 3)

	err = r.persistToDisk()
	require.NoError(t, err)

	entries, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	// Assert none of the files are persisted to disk, because of {{skip "*"}}
	assert.Len(t, entries, 0)
}

func TestRendererSkipPatternsAreRelativeToFileDirectory(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	r, err := newRenderer(ctx, nil, "./testdata/skip-is-relative/template", "./testdata/skip-is-relative/library", tmpDir)
	require.NoError(t, err)

	err = r.walk()
	assert.NoError(t, err)

	assert.Len(t, r.skipPatterns, 3)
	assert.Contains(t, r.skipPatterns, "a")
	assert.Contains(t, r.skipPatterns, "dir1/b")
	assert.Contains(t, r.skipPatterns, "dir1/dir2/c")
}

func TestRendererSkip(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	r, err := newRenderer(ctx, nil, "./testdata/skip/template", "./testdata/skip/library", tmpDir)
	require.NoError(t, err)

	err = r.walk()
	assert.NoError(t, err)
	// All 6 files are computed, even though "dir2/*" is present as a skip pattern
	// This is because "dir2/*" matches the files in dir2, but not dir2 itself
	assert.Len(t, r.files, 6)

	err = r.persistToDisk()
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(tmpDir, "file1"))
	assert.FileExists(t, filepath.Join(tmpDir, "file2"))
	assert.FileExists(t, filepath.Join(tmpDir, "dir1/file5"))

	// These files have been skipped
	assert.NoFileExists(t, filepath.Join(tmpDir, "file3"))
	assert.NoFileExists(t, filepath.Join(tmpDir, "dir1/file4"))
	assert.NoDirExists(t, filepath.Join(tmpDir, "dir2"))
	assert.NoFileExists(t, filepath.Join(tmpDir, "dir2/file6"))
}

func TestRendererInMemoryFileFullPathForWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.SkipNow()
	}
	f := &inMemoryFile{
		root:    `c:\a\b\c`,
		relPath: "d/e",
	}
	assert.Equal(t, `c:\a\b\c\d\e`, f.fullPath())
}

func TestRendererInMemoryFilePersistToDiskSetsExecutableBit(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.SkipNow()
	}
	tmpDir := t.TempDir()

	f := &inMemoryFile{
		root:    tmpDir,
		relPath: "a/b/c",
		content: []byte("123"),
		perm:    0755,
	}
	err := f.persistToDisk()
	assert.NoError(t, err)

	assertFileContent(t, filepath.Join(tmpDir, "a/b/c"), "123")
	assertFilePermissions(t, filepath.Join(tmpDir, "a/b/c"), 0755)
}

func TestRendererInMemoryFilePersistToDiskForWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.SkipNow()
	}
	tmpDir := t.TempDir()

	f := &inMemoryFile{
		root:    tmpDir,
		relPath: "a/b/c",
		content: []byte("123"),
		perm:    0666,
	}
	err := f.persistToDisk()
	assert.NoError(t, err)

	assertFileContent(t, filepath.Join(tmpDir, "a/b/c"), "123")
	assertFilePermissions(t, filepath.Join(tmpDir, "a/b/c"), 0666)
}

func TestRendererReadsPermissionsBits(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.SkipNow()
	}
	tmpDir := t.TempDir()
	ctx := context.Background()

	r, err := newRenderer(ctx, nil, "./testdata/executable-bit-read/template", "./testdata/executable-bit-read/library", tmpDir)
	require.NoError(t, err)

	err = r.walk()
	assert.NoError(t, err)

	getPermissions := func(r *renderer, path string) fs.FileMode {
		for _, f := range r.files {
			if f.relPath == path {
				return f.perm
			}
		}
		require.FailNow(t, "file is absent: "+path)
		return 0
	}

	assert.Len(t, r.files, 2)
	assert.Equal(t, getPermissions(r, "script.sh"), fs.FileMode(0755))
	assert.Equal(t, getPermissions(r, "not-a-script"), fs.FileMode(0644))
}

func TestRendererErrorOnConflictingFile(t *testing.T) {
	tmpDir := t.TempDir()

	f, err := os.Create(filepath.Join(tmpDir, "a"))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)

	r := renderer{
		skipPatterns: []string{},
		files: []*inMemoryFile{
			{
				root:    tmpDir,
				relPath: "a",
				content: []byte("123"),
				perm:    0444,
			},
		},
	}
	err = r.persistToDisk()
	assert.EqualError(t, err, fmt.Sprintf("failed to persist to disk, conflict with existing file: %s", filepath.Join(tmpDir, "a")))
}

func TestRendererNoErrorOnConflictingFileIfSkipped(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	f, err := os.Create(filepath.Join(tmpDir, "a"))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)

	r := renderer{
		ctx:          ctx,
		skipPatterns: []string{"a"},
		files: []*inMemoryFile{
			{
				root:    tmpDir,
				relPath: "a",
				content: []byte("123"),
				perm:    0444,
			},
		},
	}
	err = r.persistToDisk()
	// No error is returned even though a conflicting file exists. This is because
	// the generated file is being skipped
	assert.NoError(t, err)
	assert.Len(t, r.files, 1)
}
