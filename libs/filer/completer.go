package filer

import (
	"context"
	"fmt"
	"path"
	"sort"

	"github.com/databricks/databricks-sdk-go"
	"github.com/spf13/cobra"
)

type completer struct {
	ctx   context.Context
	wsc   *databricks.WorkspaceClient
	filer Filer
}

// General completer that takes a Filer to complete remote paths when TAB-ing through a path.
func NewCompleter(ctx context.Context, wsc *databricks.WorkspaceClient, filer Filer) *completer {
	return &completer{ctx: ctx, wsc: wsc, filer: filer}
}

func (c *completer) CompleteRemotePath(remotePath string) ([]string, cobra.ShellCompDirective) {
	_, err := c.wsc.CurrentUser.Me(c.ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// If the user is TAB-ing their way through a path, the path in `toComplete`
	// is valid and we should list nested directories.
	// If the path in `toComplete` is incomplete, however,
	// then we should list adjacent directories.
	nested := fetchDirs(c.ctx, c.filer, remotePath)
	dirs := <-nested

	if dirs == nil {
		adjacent := fetchDirs(c.ctx, c.filer, path.Dir(remotePath))
		dirs = <-adjacent
	}

	return dirs, cobra.ShellCompDirectiveNoSpace
}

func fetchDirs(
	ctx context.Context,
	filer Filer,
	remotePath string,
) <-chan []string {
	ch := make(chan []string, 1)
	go func() {
		defer close(ch)

		entries, err := filer.ReadDir(ctx, remotePath)
		if err != nil {
			return
		}

		dirs := []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				separator := ""

				// Ensure that the path and name have a "/" separating them
				if remotePath[len(remotePath)-1] != '/' {
					separator = "/"
				}

				completion := fmt.Sprintf("%s%s%s", remotePath, separator, entry.Name())
				dirs = append(dirs, completion)
			}
		}

		// Sort completions alphabetically
		sort.Slice(dirs, func(i, j int) bool {
			return dirs[i] < dirs[j]
		})

		ch <- dirs
	}()

	return ch
}
