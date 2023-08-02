package bundle

import (
	"encoding/json"
	"reflect"

	"github.com/databricks/cli/bundle/config"
	"github.com/databricks/cli/bundle/schema"
	"github.com/spf13/cobra"
)

func newSchemaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Generate JSON Schema for bundle configuration",
	}

	var openapi string
	var onlyDocs bool
	cmd.Flags().StringVar(&openapi, "openapi", "", "path to a databricks openapi spec")
	cmd.Flags().BoolVar(&onlyDocs, "only-docs", false, "only generate descriptions for the schema")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		docs, err := schema.BundleDocs(openapi)
		if err != nil {
			return err
		}
		schema, err := schema.New(reflect.TypeOf(config.Root{}), docs)
		if err != nil {
			return err
		}
		result, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return err
		}
		if onlyDocs {
			result, err = json.MarshalIndent(docs, "", "  ")
			if err != nil {
				return err
			}
		}
		cmd.OutOrStdout().Write(result)
		return nil
	}

	return cmd
}
