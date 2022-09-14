package git

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/databricks/bricks/folders"
	"github.com/databricks/bricks/utilities"
	"github.com/databricks/databricks-sdk-go/workspaces"
	giturls "github.com/whilp/git-urls"
	"gopkg.in/ini.v1"
)

func Root() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return folders.FindDirWithLeaf(wd, ".git")
}

// Origin finds the git repository the project is cloned from, so that
// we could automatically verify if this project is checked out in repos
// home folder of the user according to recommended best practices. Can
// also be used to determine a good enough default project name.
func Origin() (*url.URL, error) {
	root, err := Root()
	if err != nil {
		return nil, err
	}
	file := fmt.Sprintf("%s/.git/config", root)
	gitConfig, err := ini.Load(file)
	if err != nil {
		return nil, err
	}
	section := gitConfig.Section(`remote "origin"`)
	if section == nil {
		return nil, fmt.Errorf("remote `origin` is not defined in %s", file)
	}
	url := section.Key("url")
	if url == nil {
		return nil, fmt.Errorf("git origin url is not defined")
	}
	return giturls.Parse(url.Value())
}

// HttpsOrigin returns URL in the format expected by Databricks Repos
// platform functionality. Gradually expand implementation to work with
// other formats of git URLs.
func HttpsOrigin() (string, error) {
	origin, err := Origin()
	if err != nil {
		return "", err
	}
	// if current repo is checked out with a SSH key
	if origin.Scheme != "https" {
		origin.Scheme = "https"
	}
	// `git@` is not required for HTTPS, as Databricks Repos are checked
	// out using an API token instead of username. But does it hold true
	// for all of the git implementations?
	if origin.User != nil {
		origin.User = nil
	}
	// Remove `.git` suffix, if present.
	origin.Path = strings.TrimSuffix(origin.Path, ".git")
	return origin.String(), nil
}

// RepositoryName returns repository name as last path entry from detected
// git repository up the tree or returns error if it fails to do so.
func RepositoryName() (string, error) {
	origin, err := Origin()
	if err != nil {
		return "", err
	}
	base := path.Base(origin.Path)
	return strings.TrimSuffix(base, ".git"), nil
}

func RepoExists(remotePath string, ctx context.Context, wsc *workspaces.WorkspacesClient) (bool, error) {
	repos, err := utilities.GetAllRepos(ctx, wsc, remotePath)
	if err != nil {
		return false, fmt.Errorf("could not get repos: %s", err)
	}
	foundRepo := false
	for _, repo := range repos {
		if repo.Path == remotePath {
			foundRepo = true
			break
		}
	}
	return foundRepo, nil
}
