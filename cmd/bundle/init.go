package bundle

import (
	"path/filepath"

	"github.com/databricks/cli/libs/template"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Template",
	Long:  `Initialize template`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// initialize default value for config file path
		if configFile == "" {
			configFile = filepath.Join(targetDir, template.ConfigFileName)
		}
		return template.Materialize(args[0], targetDir)
	},
}

var targetDir string
var configFile string

func init() {
	initCmd.Flags().StringVar(&targetDir, "target-dir", ".", "path to directory template will be initialized in")
	initCmd.Flags().StringVar(&configFile, "config-file", "", "path to config to use for template initialization")
	AddCommand(initCmd)
}
