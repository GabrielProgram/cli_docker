// Code generated from OpenAPI specs by Databricks SDK Generator. DO NOT EDIT.

package vector_search_indexes

import (
	"fmt"

	"github.com/databricks/cli/cmd/root"
	"github.com/databricks/cli/libs/cmdio"
	"github.com/databricks/cli/libs/flags"
	"github.com/databricks/databricks-sdk-go/service/vectorsearch"
	"github.com/spf13/cobra"
)

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var cmdOverrides []func(*cobra.Command)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vector-search-indexes",
		Short: `**Index**: An efficient representation of your embedding vectors that supports real-time and efficient approximate nearest neighbor (ANN) search queries.`,
		Long: `**Index**: An efficient representation of your embedding vectors that supports
  real-time and efficient approximate nearest neighbor (ANN) search queries.
  
  There are 2 types of Vector Search indexes: * **Delta Sync Index**: An index
  that automatically syncs with a source Delta Table, automatically and
  incrementally updating the index as the underlying data in the Delta Table
  changes. * **Direct Vector Access Index**: An index that supports direct read
  and write of vectors and metadata through our REST and SDK APIs. With this
  model, the user manages index updates.`,
		GroupID: "vectorsearch",
		Annotations: map[string]string{
			"package": "vectorsearch",
		},
	}

	// Apply optional overrides to this command.
	for _, fn := range cmdOverrides {
		fn(cmd)
	}

	return cmd
}

// start create-index command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var createIndexOverrides []func(
	*cobra.Command,
	*vectorsearch.CreateVectorIndexRequest,
)

func newCreateIndex() *cobra.Command {
	cmd := &cobra.Command{}

	var createIndexReq vectorsearch.CreateVectorIndexRequest
	var createIndexJson flags.JsonFlag

	// TODO: short flags
	cmd.Flags().Var(&createIndexJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	// TODO: complex arg: delta_sync_vector_index_spec
	// TODO: complex arg: direct_access_index_spec
	cmd.Flags().StringVar(&createIndexReq.EndpointName, "endpoint-name", createIndexReq.EndpointName, `Name of the endpoint to be used for serving the index.`)

	cmd.Use = "create-index NAME PRIMARY_KEY INDEX_TYPE"
	cmd.Short = `Create an index.`
	cmd.Long = `Create an index.
  
  Create a new index.

  Arguments:
    NAME: Name of the index
    PRIMARY_KEY: Primary key of the index
    INDEX_TYPE: There are 2 types of Vector Search indexes:
      
      - DELTA_SYNC: An index that automatically syncs with a source Delta
      Table, automatically and incrementally updating the index as the
      underlying data in the Delta Table changes. - DIRECT_ACCESS: An index
      that supports direct read and write of vectors and metadata through our
      REST and SDK APIs. With this model, the user manages index updates.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("json") {
			err := cobra.ExactArgs(0)(cmd, args)
			if err != nil {
				return fmt.Errorf("when --json flag is specified, no positional arguments are required. Provide 'name', 'primary_key', 'index_type' in your JSON input")
			}
			return nil
		}
		check := cobra.ExactArgs(3)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		if cmd.Flags().Changed("json") {
			err = createIndexJson.Unmarshal(&createIndexReq)
			if err != nil {
				return err
			}
		}
		if !cmd.Flags().Changed("json") {
			createIndexReq.Name = args[0]
		}
		if !cmd.Flags().Changed("json") {
			createIndexReq.PrimaryKey = args[1]
		}
		if !cmd.Flags().Changed("json") {
			_, err = fmt.Sscan(args[2], &createIndexReq.IndexType)
			if err != nil {
				return fmt.Errorf("invalid INDEX_TYPE: %s", args[2])
			}
		}

		response, err := w.VectorSearchIndexes.CreateIndex(ctx, createIndexReq)
		if err != nil {
			return err
		}
		return cmdio.Render(ctx, response)
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range createIndexOverrides {
		fn(cmd, &createIndexReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newCreateIndex())
	})
}

// start delete-data-vector-index command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var deleteDataVectorIndexOverrides []func(
	*cobra.Command,
	*vectorsearch.DeleteDataVectorIndexRequest,
)

func newDeleteDataVectorIndex() *cobra.Command {
	cmd := &cobra.Command{}

	var deleteDataVectorIndexReq vectorsearch.DeleteDataVectorIndexRequest
	var deleteDataVectorIndexJson flags.JsonFlag

	// TODO: short flags
	cmd.Flags().Var(&deleteDataVectorIndexJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	cmd.Use = "delete-data-vector-index NAME"
	cmd.Short = `Delete data from index.`
	cmd.Long = `Delete data from index.
  
  Handles the deletion of data from a specified vector index.

  Arguments:
    NAME: Name of the vector index where data is to be deleted. Must be a Direct
      Vector Access Index.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(1)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		if cmd.Flags().Changed("json") {
			err = deleteDataVectorIndexJson.Unmarshal(&deleteDataVectorIndexReq)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("please provide command input in JSON format by specifying the --json flag")
		}
		deleteDataVectorIndexReq.Name = args[0]

		response, err := w.VectorSearchIndexes.DeleteDataVectorIndex(ctx, deleteDataVectorIndexReq)
		if err != nil {
			return err
		}
		return cmdio.Render(ctx, response)
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range deleteDataVectorIndexOverrides {
		fn(cmd, &deleteDataVectorIndexReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newDeleteDataVectorIndex())
	})
}

// start delete-index command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var deleteIndexOverrides []func(
	*cobra.Command,
	*vectorsearch.DeleteIndexRequest,
)

func newDeleteIndex() *cobra.Command {
	cmd := &cobra.Command{}

	var deleteIndexReq vectorsearch.DeleteIndexRequest

	// TODO: short flags

	cmd.Use = "delete-index INDEX_NAME"
	cmd.Short = `Delete an index.`
	cmd.Long = `Delete an index.
  
  Delete an index.

  Arguments:
    INDEX_NAME: Name of the index`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(1)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		deleteIndexReq.IndexName = args[0]

		err = w.VectorSearchIndexes.DeleteIndex(ctx, deleteIndexReq)
		if err != nil {
			return err
		}
		return nil
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range deleteIndexOverrides {
		fn(cmd, &deleteIndexReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newDeleteIndex())
	})
}

// start get-index command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var getIndexOverrides []func(
	*cobra.Command,
	*vectorsearch.GetIndexRequest,
)

func newGetIndex() *cobra.Command {
	cmd := &cobra.Command{}

	var getIndexReq vectorsearch.GetIndexRequest

	// TODO: short flags

	cmd.Use = "get-index INDEX_NAME"
	cmd.Short = `Get an index.`
	cmd.Long = `Get an index.
  
  Get an index.

  Arguments:
    INDEX_NAME: Name of the index`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(1)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		getIndexReq.IndexName = args[0]

		response, err := w.VectorSearchIndexes.GetIndex(ctx, getIndexReq)
		if err != nil {
			return err
		}
		return cmdio.Render(ctx, response)
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range getIndexOverrides {
		fn(cmd, &getIndexReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newGetIndex())
	})
}

// start list-indexes command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var listIndexesOverrides []func(
	*cobra.Command,
	*vectorsearch.ListIndexesRequest,
)

func newListIndexes() *cobra.Command {
	cmd := &cobra.Command{}

	var listIndexesReq vectorsearch.ListIndexesRequest

	// TODO: short flags

	cmd.Flags().StringVar(&listIndexesReq.PageToken, "page-token", listIndexesReq.PageToken, `Token for pagination.`)

	cmd.Use = "list-indexes ENDPOINT_NAME"
	cmd.Short = `List indexes.`
	cmd.Long = `List indexes.
  
  List all indexes in the given endpoint.

  Arguments:
    ENDPOINT_NAME: Name of the endpoint`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(1)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		listIndexesReq.EndpointName = args[0]

		response, err := w.VectorSearchIndexes.ListIndexesAll(ctx, listIndexesReq)
		if err != nil {
			return err
		}
		return cmdio.Render(ctx, response)
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range listIndexesOverrides {
		fn(cmd, &listIndexesReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newListIndexes())
	})
}

// start query-index command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var queryIndexOverrides []func(
	*cobra.Command,
	*vectorsearch.QueryVectorIndexRequest,
)

func newQueryIndex() *cobra.Command {
	cmd := &cobra.Command{}

	var queryIndexReq vectorsearch.QueryVectorIndexRequest
	var queryIndexJson flags.JsonFlag

	// TODO: short flags
	cmd.Flags().Var(&queryIndexJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	cmd.Flags().StringVar(&queryIndexReq.FiltersJson, "filters-json", queryIndexReq.FiltersJson, `JSON string representing query filters.`)
	cmd.Flags().IntVar(&queryIndexReq.NumResults, "num-results", queryIndexReq.NumResults, `Number of results to return.`)
	cmd.Flags().StringVar(&queryIndexReq.QueryText, "query-text", queryIndexReq.QueryText, `Query text.`)
	// TODO: array: query_vector

	cmd.Use = "query-index INDEX_NAME"
	cmd.Short = `Query an index.`
	cmd.Long = `Query an index.
  
  Query the specified vector index.

  Arguments:
    INDEX_NAME: Name of the vector index to query.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(1)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		if cmd.Flags().Changed("json") {
			err = queryIndexJson.Unmarshal(&queryIndexReq)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("please provide command input in JSON format by specifying the --json flag")
		}
		queryIndexReq.IndexName = args[0]

		response, err := w.VectorSearchIndexes.QueryIndex(ctx, queryIndexReq)
		if err != nil {
			return err
		}
		return cmdio.Render(ctx, response)
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range queryIndexOverrides {
		fn(cmd, &queryIndexReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newQueryIndex())
	})
}

// start sync-index command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var syncIndexOverrides []func(
	*cobra.Command,
	*vectorsearch.SyncIndexRequest,
)

func newSyncIndex() *cobra.Command {
	cmd := &cobra.Command{}

	var syncIndexReq vectorsearch.SyncIndexRequest

	// TODO: short flags

	cmd.Use = "sync-index INDEX_NAME"
	cmd.Short = `Synchronize an index.`
	cmd.Long = `Synchronize an index.
  
  Triggers a synchronization process for a specified vector index.

  Arguments:
    INDEX_NAME: Name of the vector index to synchronize. Must be a Delta Sync Index.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(1)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		syncIndexReq.IndexName = args[0]

		err = w.VectorSearchIndexes.SyncIndex(ctx, syncIndexReq)
		if err != nil {
			return err
		}
		return nil
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range syncIndexOverrides {
		fn(cmd, &syncIndexReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newSyncIndex())
	})
}

// start upsert-data-vector-index command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var upsertDataVectorIndexOverrides []func(
	*cobra.Command,
	*vectorsearch.UpsertDataVectorIndexRequest,
)

func newUpsertDataVectorIndex() *cobra.Command {
	cmd := &cobra.Command{}

	var upsertDataVectorIndexReq vectorsearch.UpsertDataVectorIndexRequest
	var upsertDataVectorIndexJson flags.JsonFlag

	// TODO: short flags
	cmd.Flags().Var(&upsertDataVectorIndexJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	cmd.Use = "upsert-data-vector-index NAME INPUTS_JSON"
	cmd.Short = `Upsert data into an index.`
	cmd.Long = `Upsert data into an index.
  
  Handles the upserting of data into a specified vector index.

  Arguments:
    NAME: Name of the vector index where data is to be upserted. Must be a Direct
      Vector Access Index.
    INPUTS_JSON: JSON string representing the data to be upserted.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("json") {
			err := cobra.ExactArgs(1)(cmd, args)
			if err != nil {
				return fmt.Errorf("when --json flag is specified, provide only NAME as positional arguments. Provide 'inputs_json' in your JSON input")
			}
			return nil
		}
		check := cobra.ExactArgs(2)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		if cmd.Flags().Changed("json") {
			err = upsertDataVectorIndexJson.Unmarshal(&upsertDataVectorIndexReq)
			if err != nil {
				return err
			}
		}
		upsertDataVectorIndexReq.Name = args[0]
		if !cmd.Flags().Changed("json") {
			upsertDataVectorIndexReq.InputsJson = args[1]
		}

		response, err := w.VectorSearchIndexes.UpsertDataVectorIndex(ctx, upsertDataVectorIndexReq)
		if err != nil {
			return err
		}
		return cmdio.Render(ctx, response)
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range upsertDataVectorIndexOverrides {
		fn(cmd, &upsertDataVectorIndexReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newUpsertDataVectorIndex())
	})
}

// end service VectorSearchIndexes
