package gitcredentials

import (
	git_credentials "github.com/databricks/bricks/cmd/gitcredentials/git-credentials"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "gitcredentials",
	Short: `Registers personal access token for Databricks to do operations on behalf of the user.`,
	Long: `Registers personal access token for Databricks to do operations on behalf of
  the user.`,
}

func init() {
	Cmd.PersistentFlags().String("profile", "", "~/.databrickscfg profile")

	Cmd.AddCommand(git_credentials.Cmd)
}
