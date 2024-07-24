package fs

import (
	"context"
	"errors"
	"io/fs"
	"runtime"
	"testing"

	"github.com/databricks/cli/cmd/root"
	"github.com/databricks/cli/internal/testutil"
	"github.com/databricks/cli/libs/filer"
	"github.com/databricks/databricks-sdk-go/experimental/mocks"
	"github.com/databricks/databricks-sdk-go/service/iam"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFilerForPathForLocalPaths(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	f, path, err := filerForPath(ctx, tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, tmpDir, path)

	info, err := f.Stat(ctx, path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestFilerForPathForInvalidScheme(t *testing.T) {
	ctx := context.Background()

	_, _, err := filerForPath(ctx, "dbf:/a")
	assert.ErrorContains(t, err, "invalid scheme")

	_, _, err = filerForPath(ctx, "foo:a")
	assert.ErrorContains(t, err, "invalid scheme")

	_, _, err = filerForPath(ctx, "file:/a")
	assert.ErrorContains(t, err, "invalid scheme")
}

func testWindowsFilerForPath(t *testing.T, ctx context.Context, fullPath string) {
	f, path, err := filerForPath(ctx, fullPath)
	assert.NoError(t, err)

	// Assert path remains unchanged
	assert.Equal(t, path, fullPath)

	// Assert local client is created
	_, ok := f.(*filer.LocalClient)
	assert.True(t, ok)
}

func TestFilerForWindowsLocalPaths(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.SkipNow()
	}

	ctx := context.Background()
	testWindowsFilerForPath(t, ctx, `c:\abc`)
	testWindowsFilerForPath(t, ctx, `c:abc`)
	testWindowsFilerForPath(t, ctx, `d:\abc`)
	testWindowsFilerForPath(t, ctx, `d:\abc`)
	testWindowsFilerForPath(t, ctx, `f:\abc\ef`)
}

func mockCurrentUserApi(m *mocks.MockWorkspaceClient, u *iam.User, e error) {
	currentUserApi := m.GetMockCurrentUserAPI()
	currentUserApi.EXPECT().Me(mock.AnythingOfType("*context.valueCtx")).Return(u, e)
}

func setupValidArgsFunctionTest(t *testing.T) (*mocks.MockWorkspaceClient, *cobra.Command) {
	m := mocks.NewMockWorkspaceClient(t)
	ctx := context.Background()
	ctx = root.SetWorkspaceClient(ctx, m.WorkspaceClient)

	cmd := &cobra.Command{}
	cmd.SetContext(ctx)

	return m, cmd
}

func TestGetValidArgsFunctionCompletion(t *testing.T) {
	m, cmd := setupValidArgsFunctionTest(t)

	mockCurrentUserApi(m, nil, nil)

	mockFilerForPath := testutil.GetMockFilerForPath(t, []fs.DirEntry{
		testutil.NewFakeDirEntry("dir", true),
		testutil.NewFakeDirEntry("file", false),
	})

	validArgsFunction := getValidArgsFunction(1, false, mockFilerForPath)
	completions, directive := validArgsFunction(cmd, []string{}, "dbfs:/")

	assert.Equal(t, []string{"dbfs:/dir", "dbfs:/file"}, completions)
	assert.Equal(t, cobra.ShellCompDirectiveNoSpace, directive)
}

func TestGetValidArgsFunctionCompletionOnlyDirs(t *testing.T) {
	m, cmd := setupValidArgsFunctionTest(t)

	mockCurrentUserApi(m, nil, nil)

	mockFilerForPath := testutil.GetMockFilerForPath(t, []fs.DirEntry{
		testutil.NewFakeDirEntry("dir", true),
		testutil.NewFakeDirEntry("file", false),
	})

	validArgsFunction := getValidArgsFunction(1, true, mockFilerForPath)
	completions, directive := validArgsFunction(cmd, []string{}, "dbfs:/")

	assert.Equal(t, []string{"dbfs:/dir"}, completions)
	assert.Equal(t, cobra.ShellCompDirectiveNoSpace, directive)
}

func TestGetValidArgsFunctionNoDbfsPath(t *testing.T) {
	_, cmd := setupValidArgsFunctionTest(t)

	mockFilerForPath := testutil.GetMockFilerForPath(t, nil)

	validArgsFunction := getValidArgsFunction(1, true, mockFilerForPath)
	completions, directive := validArgsFunction(cmd, []string{}, "/")

	assert.Equal(t, []string{"dbfs:/"}, completions)
	assert.Equal(t, cobra.ShellCompDirectiveNoSpace, directive)
}

func TestGetValidArgsFunctionNoCurrentUser(t *testing.T) {
	m, cmd := setupValidArgsFunctionTest(t)

	mockCurrentUserApi(m, nil, errors.New("Current User Error"))
	mockFilerForPath := testutil.GetMockFilerForPath(t, nil)

	validArgsFunction := getValidArgsFunction(1, true, mockFilerForPath)
	completions, directive := validArgsFunction(cmd, []string{}, "dbfs:/")

	assert.Nil(t, completions)
	assert.Equal(t, cobra.ShellCompDirectiveError, directive)
}

func TestGetValidArgsFunctionNotCompletedArgument(t *testing.T) {
	m, cmd := setupValidArgsFunctionTest(t)

	mockCurrentUserApi(m, nil, nil)
	mockFilerForPath := testutil.GetMockFilerForPath(t, nil)

	// 0 here means we don't complete any arguments
	validArgsFunction := getValidArgsFunction(0, true, mockFilerForPath)
	completions, directive := validArgsFunction(cmd, []string{}, "dbfs:/")

	assert.Nil(t, completions)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}
