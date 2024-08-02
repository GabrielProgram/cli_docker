package completer

import (
	"context"
	"path"

	"github.com/databricks/cli/libs/filer"
	"github.com/spf13/cobra"
)

type completer struct {
	ctx      context.Context
	filer    filer.Filer
	onlyDirs bool
}

// General completer that takes a Filer to complete remote paths when TAB-ing through a path.
func NewCompleter(ctx context.Context, filer filer.Filer, onlyDirs bool) *completer {
	return &completer{ctx: ctx, filer: filer, onlyDirs: onlyDirs}
}

func (c *completer) CompletePath(p string) ([]string, string, cobra.ShellCompDirective) {
	// If the user is TAB-ing their way through a path, the path in `toComplete`
	// is valid and we should list nested directories.
	// If the path in `toComplete` is incomplete, however,
	// then we should list adjacent directories.
	dirPath := p
	nested := fetchCompletions(c, dirPath)
	completions := <-nested
	if completions == nil {
		dirPath = path.Dir(p)
		adjacent := fetchCompletions(c, dirPath)
		completions = <-adjacent
	}

	return completions, dirPath, cobra.ShellCompDirectiveNoSpace
}

func fetchCompletions(
	c *completer,
	path string,
) <-chan []string {
	ch := make(chan []string, 1)
	go func() {
		defer close(ch)

		entries, err := c.filer.ReadDir(c.ctx, path)
		if err != nil {
			return
		}

		completions := []string{}
		for _, entry := range entries {
			if !c.onlyDirs || entry.IsDir() {
				completions = append(completions, entry.Name())
			}
		}

		ch <- completions
	}()

	return ch
}
