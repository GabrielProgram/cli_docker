package command_execution

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "command-execution",
	Short: `This API allows execution of Python, Scala, SQL, or R commands on running Databricks Clusters.`, // TODO: fix FirstSentence logic and append dot to summary
}
