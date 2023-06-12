package fs

import (
	"github.com/databricks/cli/cmd/root"
	"github.com/databricks/cli/libs/filer"
	"github.com/spf13/cobra"
)

var mkdirCmd = &cobra.Command{
	Use: "mkdir DIR_PATH",
	// Alias `mkdirs` for this command exists for legacy purposes. This command
	// is called databricks fs mkdirs in our legacy CLI: https://github.com/databricks/databricks-cli
	Aliases: []string{"mkdirs"},
	Short:   "Make directories",
	Long:    `Mkdir will create directories along the path to the argument directory.`,
	Args:    cobra.ExactArgs(1),
	PreRunE: root.MustWorkspaceClient,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		path, err := resolveDbfsPath(args[0])
		if err != nil {
			return err
		}

		f, err := filer.NewDbfsClient(w, "/")
		if err != nil {
			return err
		}

		return f.Mkdir(ctx, path)
	},
}

func init() {
	fsCmd.AddCommand(mkdirCmd)
}
