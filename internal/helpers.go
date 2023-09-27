package internal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/databricks/cli/cmd"
	_ "github.com/databricks/cli/cmd/version"
	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/apierr"
	"github.com/databricks/databricks-sdk-go/service/compute"
	"github.com/databricks/databricks-sdk-go/service/files"
	"github.com/databricks/databricks-sdk-go/service/jobs"
	"github.com/databricks/databricks-sdk-go/service/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	_ "github.com/databricks/cli/cmd/workspace"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GetEnvOrSkipTest proceeds with test only with that env variable
func GetEnvOrSkipTest(t *testing.T, name string) string {
	value := os.Getenv(name)
	if value == "" {
		t.Skipf("Environment variable %s is missing", name)
	}
	return value
}

// RandomName gives random name with optional prefix. e.g. qa.RandomName("tf-")
func RandomName(prefix ...string) string {
	randLen := 12
	b := make([]byte, randLen)
	for i := range b {
		b[i] = charset[rand.Intn(randLen)]
	}
	if len(prefix) > 0 {
		return fmt.Sprintf("%s%s", strings.Join(prefix, ""), b)
	}
	return string(b)
}

// Helper for running the root command in the background.
// It ensures that the background goroutine terminates upon
// test completion through cancelling the command context.
type cobraTestRunner struct {
	*testing.T

	args   []string
	stdout bytes.Buffer
	stderr bytes.Buffer

	ctx context.Context

	// Line-by-line output.
	// Background goroutines populate these channels by reading from stdout/stderr pipes.
	stdoutLines <-chan string
	stderrLines <-chan string

	errch <-chan error
}

func consumeLines(ctx context.Context, wg *sync.WaitGroup, r io.Reader) <-chan string {
	ch := make(chan string, 1000)
	wg.Add(1)
	go func() {
		defer close(ch)
		defer wg.Done()
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- scanner.Text():
			}
		}
	}()
	return ch
}

func (t *cobraTestRunner) registerFlagCleanup(c *cobra.Command) {
	// Find target command that will be run. Example: if the command run is `databricks fs cp`,
	// target command corresponds to `cp`
	targetCmd, _, err := c.Find(t.args)
	require.NoError(t, err)

	// Force initialization of default flags.
	// These are initialized by cobra at execution time and would otherwise
	// not be cleaned up by the cleanup function below.
	targetCmd.InitDefaultHelpFlag()
	targetCmd.InitDefaultVersionFlag()

	// Restore flag values to their original value on test completion.
	targetCmd.Flags().VisitAll(func(f *pflag.Flag) {
		v := reflect.ValueOf(f.Value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		// Store copy of the current flag value.
		reset := reflect.New(v.Type()).Elem()
		reset.Set(v)
		t.Cleanup(func() {
			v.Set(reset)
		})
	})
}

func (t *cobraTestRunner) RunBackground() {
	var stdoutR, stderrR io.Reader
	var stdoutW, stderrW io.WriteCloser
	stdoutR, stdoutW = io.Pipe()
	stderrR, stderrW = io.Pipe()
	root := cmd.New(context.Background())
	root.SetOut(stdoutW)
	root.SetErr(stderrW)
	root.SetArgs(t.args)

	// Register cleanup function to restore flags to their original values
	// once test has been executed. This is needed because flag values reside
	// in a global singleton data-structure, and thus subsequent tests might
	// otherwise interfere with each other
	t.registerFlagCleanup(root)

	errch := make(chan error)
	ctx, cancel := context.WithCancel(t.ctx)

	// Tee stdout/stderr to buffers.
	stdoutR = io.TeeReader(stdoutR, &t.stdout)
	stderrR = io.TeeReader(stderrR, &t.stderr)

	// Consume stdout/stderr line-by-line.
	var wg sync.WaitGroup
	t.stdoutLines = consumeLines(ctx, &wg, stdoutR)
	t.stderrLines = consumeLines(ctx, &wg, stderrR)

	// Run command in background.
	go func() {
		cmd, err := root.ExecuteContextC(ctx)
		if err != nil {
			t.Logf("Error running command: %s", err)
		}

		// Close pipes to signal EOF.
		stdoutW.Close()
		stderrW.Close()

		// Wait for the [consumeLines] routines to finish now that
		// the pipes they're reading from have closed.
		wg.Wait()

		if t.stdout.Len() > 0 {
			// Make a copy of the buffer such that it remains "unread".
			scanner := bufio.NewScanner(bytes.NewBuffer(t.stdout.Bytes()))
			for scanner.Scan() {
				t.Logf("[databricks stdout]: %s", scanner.Text())
			}
		}

		if t.stderr.Len() > 0 {
			// Make a copy of the buffer such that it remains "unread".
			scanner := bufio.NewScanner(bytes.NewBuffer(t.stderr.Bytes()))
			for scanner.Scan() {
				t.Logf("[databricks stderr]: %s", scanner.Text())
			}
		}

		// Reset context on command for the next test.
		// These commands are globals so we have to clean up to the best of our ability after each run.
		// See https://github.com/spf13/cobra/blob/a6f198b635c4b18fff81930c40d464904e55b161/command.go#L1062-L1066
		//lint:ignore SA1012 cobra sets the context and doesn't clear it
		cmd.SetContext(nil)

		// Make caller aware of error.
		errch <- err
		close(errch)
	}()

	// Ensure command terminates upon test completion (success or failure).
	t.Cleanup(func() {
		// Signal termination of command.
		cancel()
		// Wait for goroutine to finish.
		<-errch
	})

	t.errch = errch
}

func (t *cobraTestRunner) Run() (bytes.Buffer, bytes.Buffer, error) {
	t.RunBackground()
	err := <-t.errch
	return t.stdout, t.stderr, err
}

// Like [require.Eventually] but errors if the underlying command has failed.
func (c *cobraTestRunner) Eventually(condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) {
	ch := make(chan bool, 1)

	timer := time.NewTimer(waitFor)
	defer timer.Stop()

	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	// Kick off condition check immediately.
	go func() { ch <- condition() }()

	for tick := ticker.C; ; {
		select {
		case err := <-c.errch:
			require.Fail(c, "Command failed", err)
			return
		case <-timer.C:
			require.Fail(c, "Condition never satisfied", msgAndArgs...)
			return
		case <-tick:
			tick = nil
			go func() { ch <- condition() }()
		case v := <-ch:
			if v {
				return
			}
			tick = ticker.C
		}
	}
}

func NewCobraTestRunner(t *testing.T, args ...string) *cobraTestRunner {
	return &cobraTestRunner{
		T:    t,
		ctx:  context.Background(),
		args: args,
	}
}

func NewCobraTestRunnerWithContext(t *testing.T, ctx context.Context, args ...string) *cobraTestRunner {
	return &cobraTestRunner{
		T:    t,
		ctx:  ctx,
		args: args,
	}
}

func RequireSuccessfulRun(t *testing.T, args ...string) (bytes.Buffer, bytes.Buffer) {
	t.Logf("run args: [%s]", strings.Join(args, ", "))
	c := NewCobraTestRunner(t, args...)
	stdout, stderr, err := c.Run()
	require.NoError(t, err)
	return stdout, stderr
}

func RequireErrorRun(t *testing.T, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	c := NewCobraTestRunner(t, args...)
	stdout, stderr, err := c.Run()
	require.Error(t, err)
	return stdout, stderr, err
}

func writeFile(t *testing.T, name string, body string) string {
	f, err := os.Create(filepath.Join(t.TempDir(), name))
	require.NoError(t, err)
	_, err = f.WriteString(body)
	require.NoError(t, err)
	f.Close()
	return f.Name()
}

func GenerateNotebookTasks(notebookPath string, versions []string, nodeTypeId string) []jobs.SubmitTask {
	tasks := make([]jobs.SubmitTask, 0)
	for i := 0; i < len(versions); i++ {
		task := jobs.SubmitTask{
			TaskKey: fmt.Sprintf("notebook_%s", strings.ReplaceAll(versions[i], ".", "_")),
			NotebookTask: &jobs.NotebookTask{
				NotebookPath: notebookPath,
			},
			NewCluster: &compute.ClusterSpec{
				SparkVersion:     versions[i],
				NumWorkers:       1,
				NodeTypeId:       nodeTypeId,
				DataSecurityMode: compute.DataSecurityModeUserIsolation,
			},
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func GenerateSparkPythonTasks(notebookPath string, versions []string, nodeTypeId string) []jobs.SubmitTask {
	tasks := make([]jobs.SubmitTask, 0)
	for i := 0; i < len(versions); i++ {
		task := jobs.SubmitTask{
			TaskKey: fmt.Sprintf("spark_%s", strings.ReplaceAll(versions[i], ".", "_")),
			SparkPythonTask: &jobs.SparkPythonTask{
				PythonFile: notebookPath,
			},
			NewCluster: &compute.ClusterSpec{
				SparkVersion:     versions[i],
				NumWorkers:       1,
				NodeTypeId:       nodeTypeId,
				DataSecurityMode: compute.DataSecurityModeUserIsolation,
			},
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func GenerateWheelTasks(wheelPath string, versions []string, nodeTypeId string) []jobs.SubmitTask {
	tasks := make([]jobs.SubmitTask, 0)
	for i := 0; i < len(versions); i++ {
		task := jobs.SubmitTask{
			TaskKey: fmt.Sprintf("whl_%s", strings.ReplaceAll(versions[i], ".", "_")),
			PythonWheelTask: &jobs.PythonWheelTask{
				PackageName: "my_test_code",
				EntryPoint:  "run",
			},
			NewCluster: &compute.ClusterSpec{
				SparkVersion:     versions[i],
				NumWorkers:       1,
				NodeTypeId:       nodeTypeId,
				DataSecurityMode: compute.DataSecurityModeUserIsolation,
			},
			Libraries: []compute.Library{
				{Whl: wheelPath},
			},
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func TemporaryWorkspaceDir(t *testing.T, w *databricks.WorkspaceClient) string {
	ctx := context.Background()
	me, err := w.CurrentUser.Me(ctx)
	require.NoError(t, err)

	basePath := fmt.Sprintf("/Users/%s/%s", me.UserName, RandomName("integration-test-wsfs-"))

	t.Logf("Creating %s", basePath)
	err = w.Workspace.MkdirsByPath(ctx, basePath)
	require.NoError(t, err)

	// Remove test directory on test completion.
	t.Cleanup(func() {
		t.Logf("Removing %s", basePath)
		err := w.Workspace.Delete(ctx, workspace.Delete{
			Path:      basePath,
			Recursive: true,
		})
		if err == nil || apierr.IsMissing(err) {
			return
		}
		t.Logf("Unable to remove temporary workspace directory %s: %#v", basePath, err)
	})

	return basePath
}

func TemporaryDbfsDir(t *testing.T, w *databricks.WorkspaceClient) string {
	ctx := context.Background()
	path := fmt.Sprintf("/tmp/%s", RandomName("integration-test-dbfs-"))

	t.Logf("Creating DBFS folder:%s", path)
	err := w.Dbfs.MkdirsByPath(ctx, path)
	require.NoError(t, err)

	t.Cleanup(func() {
		t.Logf("Removing DBFS folder:%s", path)
		err := w.Dbfs.Delete(ctx, files.Delete{
			Path:      path,
			Recursive: true,
		})
		if err == nil || apierr.IsMissing(err) {
			return
		}
		t.Logf("unable to remove temporary dbfs directory %s: %#v", path, err)
	})

	return path
}

func TemporaryRepo(t *testing.T, w *databricks.WorkspaceClient) string {
	ctx := context.Background()
	me, err := w.CurrentUser.Me(ctx)
	require.NoError(t, err)

	repoPath := fmt.Sprintf("/Repos/%s/%s", me.UserName, RandomName("integration-test-repo-"))

	t.Logf("Creating repo:%s", repoPath)
	repoInfo, err := w.Repos.Create(ctx, workspace.CreateRepo{
		Url:      "https://github.com/databricks/cli",
		Provider: "github",
		Path:     repoPath,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		t.Logf("Removing repo: %s", repoPath)
		err := w.Repos.Delete(ctx, workspace.DeleteRepoRequest{
			RepoId: repoInfo.Id,
		})
		if err == nil || apierr.IsMissing(err) {
			return
		}
		t.Logf("unable to remove repo %s: %#v", repoPath, err)
	})

	return repoPath
}

func GetNodeTypeId(env string) string {
	if env == "gcp" {
		return "n1-standard-4"
	} else if env == "aws" {
		return "i3.xlarge"
	}
	return "Standard_DS4_v2"
}
