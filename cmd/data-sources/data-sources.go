// Code generated from OpenAPI specs by Databricks SDK Generator. DO NOT EDIT.

package data_sources

import (
	"github.com/databricks/bricks/cmd/root"
	"github.com/databricks/bricks/lib/sdk"
	"github.com/databricks/bricks/lib/ui"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "data-sources",
	Short: `This API is provided to assist you in making new query objects.`,
	Long: `This API is provided to assist you in making new query objects. When creating
  a query object, you may optionally specify a data_source_id for the SQL
  warehouse against which it will run. If you don't already know the
  data_source_id for your desired SQL warehouse, this API will help you find
  it.
  
  This API does not support searches. It returns the full list of SQL warehouses
  in your workspace. We advise you to use any text editor, REST client, or
  grep to search the response from this API for the name of your SQL warehouse
  as it appears in Databricks SQL.`,
}

// start list command

func init() {
	Cmd.AddCommand(listCmd)

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: `Get a list of SQL warehouses.`,
	Long: `Get a list of SQL warehouses.
  
  Retrieves a full list of SQL warehouses available in this workspace. All
  fields that appear in this API response are enumerated for clarity. However,
  you need only a SQL warehouse's id to create new queries against it.`,

	Annotations: map[string]string{},
	PreRunE:     root.TryWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.DataSources.List(ctx)
		if err != nil {
			return err
		}
		return ui.Render(cmd, response)
	},
}

// end service DataSources
