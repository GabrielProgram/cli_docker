package python

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/config/mutator"
	jobs_utils "github.com/databricks/cli/libs/jobs"
	"github.com/databricks/databricks-sdk-go/service/jobs"
)

//go:embed trampoline_data/notebook.py
var notebookTrampolineData string

//go:embed trampoline_data/python.py
var pyTrampolineData string

func TransforNotebookTask() bundle.Mutator {
	return mutator.NewTrampoline(
		"python_notebook",
		&notebookTrampoline{},
	)
}

type notebookTrampoline struct{}

func localNotebookPath(b *bundle.Bundle, task *jobs.Task) (string, error) {
	remotePath := task.NotebookTask.NotebookPath
	relRemotePath, err := filepath.Rel(b.Config.Workspace.FilesPath, remotePath)
	if err != nil {
		return "", err
	}
	localPath := filepath.Join(b.Config.Path, filepath.FromSlash(relRemotePath))
	_, err = os.Stat(fmt.Sprintf("%s.ipynb", localPath))
	if err == nil {
		return fmt.Sprintf("%s.ipynb", localPath), nil
	}

	_, err = os.Stat(fmt.Sprintf("%s.py", localPath))
	if err == nil {
		return fmt.Sprintf("%s.py", localPath), nil
	}

	return "", fmt.Errorf("notebook %s not found", localPath)
}

func (n *notebookTrampoline) GetTasks(b *bundle.Bundle) []jobs_utils.TaskWithJobKey {
	return jobs_utils.GetTasksWithJobKeyBy(b, func(task *jobs.Task) bool {
		if task.NotebookTask == nil ||
			task.NotebookTask.Source == jobs.SourceGit {
			return false
		}
		_, err := localNotebookPath(b, task)
		// We assume if the notebook is not available locally in the bundle
		// then the user has it somewhere in the workspace. For these
		// out of bundle notebooks we do not want to write a trampoline.
		return err == nil
	})
}

func (n *notebookTrampoline) CleanUp(task *jobs.Task) error {
	return nil
}

func (n *notebookTrampoline) GetTemplate(b *bundle.Bundle, task *jobs.Task) (string, error) {
	localPath, err := localNotebookPath(b, task)
	if err != nil {
		return "", err
	}

	bytesData, err := os.ReadFile(localPath)
	if err != nil {
		return "", err
	}
	s := strings.TrimSpace(string(bytesData))
	if strings.HasSuffix(localPath, ".ipynb") {
		return getIpynbTemplate(s)
	}

	lines := strings.Split(s, "\n")
	if strings.HasPrefix(lines[0], "# Databricks notebook source") {
		return getDbnbTemplate(s)
	}

	return pyTrampolineData, nil
}

func getDbnbTemplate(s string) (string, error) {
	s = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "# Databricks notebook source"))
	return fmt.Sprintf(`# Databricks notebook source
%s
# Command ----------
%s
`, notebookTrampolineData, s), nil
}

type IpynbData struct {
	Cells []IpynbCell `json:"cells"`
}

type IpynbCell struct {
	CellType string   `json:"cell_type"`
	Source   []string `json:"source"`
}

func getIpynbTemplate(s string) (string, error) {
	var data IpynbData
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		return "", err
	}

	if data.Cells == nil {
		data.Cells = []IpynbCell{}
	}

	data.Cells = append([]IpynbCell{
		{
			CellType: "code",
			Source:   []string{"# Databricks notebook source"},
		},
	}, data.Cells...)

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (n *notebookTrampoline) GetTemplateData(b *bundle.Bundle, task *jobs.Task) (map[string]any, error) {
	return map[string]any{
		"ProjectRoot": b.Config.Workspace.FilesPath,
		"SourceFile":  task.NotebookTask.NotebookPath,
	}, nil
}
