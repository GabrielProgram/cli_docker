package fs

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/databricks/cli/cmd/root"
	"github.com/databricks/cli/libs/filer"
	"github.com/spf13/cobra"
)

func filerForPath(ctx context.Context, fullPath string) (filer.Filer, string, error) {
	// Split path at : to detect any file schemes
	parts := strings.SplitN(fullPath, ":", 2)

	// If no scheme is specified, then local path
	if len(parts) < 2 {
		f, err := filer.NewLocalClient("")
		return f, fullPath, err
	}

	// On windows systems, paths start with a drive letter. If the scheme
	// is a single letter and the OS is windows, then we conclude the path
	// is meant to be a local path.
	if runtime.GOOS == "windows" && len(parts[0]) == 1 {
		f, err := filer.NewLocalClient("")
		return f, fullPath, err
	}

	if parts[0] != "dbfs" {
		return nil, "", fmt.Errorf("invalid scheme: %s", parts[0])
	}

	path := parts[1]
	w := root.WorkspaceClient(ctx)

	// If the specified path has the "Volumes" prefix, use the Files API.
	if strings.HasPrefix(path, "/Volumes/") {
		f, err := filer.NewFilesClient(w, "/")
		return f, path, err
	}

	// The file is a dbfs file, and uses the DBFS APIs
	f, err := filer.NewDbfsClient(w, "/")
	return f, path, err
}

const DbfsPrefix string = "dbfs:/"

func isDbfsPath(path string) bool {
	return strings.HasPrefix(path, DbfsPrefix)
}

func getValidArgsFunction(pathArgCount int) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		cmd.SetContext(root.SkipPrompt(cmd.Context()))

		err := mustWorkspaceClient(cmd, args)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		isValidPrefix := isDbfsPath(toComplete)

		if !isValidPrefix {
			return []string{DbfsPrefix}, cobra.ShellCompDirectiveNoSpace
		}

		f, path, err := filerForPath(cmd.Context(), toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		wsc := root.WorkspaceClient(cmd.Context())
		completer := filer.NewCompleter(cmd.Context(), wsc, f)

		if len(args) < pathArgCount {
			completions, directive := completer.CompleteRemotePath(path)

			// DbfsPrefix without trailing "/"
			prefix := DbfsPrefix[:len(DbfsPrefix)-1]

			// Add the prefix to the completions
			for i, completion := range completions {
				completions[i] = fmt.Sprintf("%s%s", prefix, completion)
			}

			return completions, directive
		} else {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
	}
}

// Wrapper for [root.MustWorkspaceClient] that disables loading authentication configuration from a bundle.
func mustWorkspaceClient(cmd *cobra.Command, args []string) error {
	cmd.SetContext(root.SkipLoadBundle(cmd.Context()))
	return root.MustWorkspaceClient(cmd, args)
}
