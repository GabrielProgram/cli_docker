package repofiles

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/apierr"
	"github.com/databricks/databricks-sdk-go/client"
	"github.com/databricks/databricks-sdk-go/service/workspace"
)

// Use this class to do file upload/delete on a workspace repo
//
// This class comes with safeguards when mutating remote files to prevent
// accidental deletion of repos and more robust methods to overwrite workspace files
type RepoFiles struct {
	repoRoot        string
	localRoot       string
	workspaceClient *databricks.WorkspaceClient
}

func Create(repoRoot, localRoot string, workspaceClient *databricks.WorkspaceClient) *RepoFiles {
	return &RepoFiles{
		repoRoot:        repoRoot,
		localRoot:       localRoot,
		workspaceClient: workspaceClient,
	}
}

func cleanPath(relativePath string) (string, error) {
	cleanRelativePath := path.Clean(relativePath)
	if strings.Contains(cleanRelativePath, `..`) {
		return "", fmt.Errorf(`file relative path %s contains forbidden pattern ".."`, relativePath)
	}
	if cleanRelativePath == "" || cleanRelativePath == "/" || cleanRelativePath == "." {
		return "", fmt.Errorf("file path relative to repo root cannot be empty: %s", relativePath)
	}
	return cleanRelativePath, nil
}

func (r *RepoFiles) remotePath(relativePath string) (string, error) {
	cleanRelativePath, err := cleanPath(relativePath)
	if err != nil {
		return "", err
	}
	return path.Join(r.repoRoot, cleanRelativePath), nil
}

func (r *RepoFiles) readLocal(relativePath string) ([]byte, error) {
	localPath := filepath.Join(r.localRoot, relativePath)
	return os.ReadFile(localPath)
}

func (r *RepoFiles) writeRemote(ctx context.Context, relativePath string, content []byte) error {
	apiClient, err := client.New(r.workspaceClient.Config)
	if err != nil {
		return err
	}
	remotePath, err := r.remotePath(relativePath)
	if err != nil {
		return err
	}
	apiPath := fmt.Sprintf(
		"/api/2.0/workspace-files/import-file/%s?overwrite=true",
		strings.TrimLeft(remotePath, "/"))

	err = apiClient.Do(ctx, http.MethodPost, apiPath, content, nil)

	// Handling some edge cases when an upload might fail
	//
	// We cannot do more precise error scoping here because the API does not
	// provide descriptive errors yet
	//
	// TODO: narrow down the error condition scope of this "if" block to only
	// trigger for the specific edge cases instead of all errors once the API
	// implements them
	if err != nil {
		// Delete any artifact files incase non overwriteable by the current file
		// type and thus are failing the PUT request.
		// files, folders and notebooks might not have been cleaned up and they
		// can't overwrite each other. If a folder `foo` exists, then attempts to
		// PUT a file `foo` will fail
		err := r.workspaceClient.Workspace.Delete(ctx,
			workspace.Delete{
				Path:      remotePath,
				Recursive: true,
			},
		)
		// ignore RESOURCE_DOES_NOT_EXIST here incase nothing existed at remotePath
		if val, ok := err.(apierr.APIError); ok && val.ErrorCode == "RESOURCE_DOES_NOT_EXIST" {
			err = nil
		}
		if err != nil {
			return err
		}

		// Mkdir parent dirs incase they are what's causing the PUT request to
		// fail
		err = r.workspaceClient.Workspace.MkdirsByPath(ctx, path.Dir(remotePath))
		if err != nil {
			return fmt.Errorf("could not mkdir to put file: %s", err)
		}

		// Attempt to upload file again after cleanup/setup
		err = apiClient.Do(ctx, http.MethodPost, apiPath, content, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RepoFiles) deleteRemote(ctx context.Context, relativePath string) error {
	remotePath, err := r.remotePath(relativePath)
	if err != nil {
		return err
	}
	return r.workspaceClient.Workspace.Delete(ctx,
		workspace.Delete{
			Path:      remotePath,
			Recursive: false,
		},
	)
}

// The API calls for a python script foo.py would be
// `PUT foo.py`
// `DELETE foo.py`
//
// The API calls for a python notebook foo.py would be
// `PUT foo.py`
// `DELETE foo`
//
// The workspace file system backend strips .py from the file name if the python
// file is a notebook
func (r *RepoFiles) PutFile(ctx context.Context, relativePath string) error {
	content, err := r.readLocal(relativePath)
	if err != nil {
		return err
	}

	return r.writeRemote(ctx, relativePath, content)
}

func (r *RepoFiles) DeleteFile(ctx context.Context, relativePath string) error {
	err := r.deleteRemote(ctx, relativePath)

	// We explictly ignore RESOURCE_DOES_NOT_EXIST error to make delete idempotent
	if val, ok := err.(apierr.APIError); ok && val.ErrorCode == "RESOURCE_DOES_NOT_EXIST" {
		return nil
	}
	return nil
}
