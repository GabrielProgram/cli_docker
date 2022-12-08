package external_locations

import (
	"github.com/databricks/bricks/lib/sdk"
	"github.com/databricks/bricks/lib/ui"
	"github.com/databricks/databricks-sdk-go/service/unitycatalog"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "external-locations",
	Short: `An external location is an object that combines a cloud storage path with a storage credential that authorizes access to the cloud storage path.`,
	Long: `An external location is an object that combines a cloud storage path with a
  storage credential that authorizes access to the cloud storage path. Each
  storage location is subject to Unity Catalog access-control policies that
  control which users and groups can access the credential. If a user does not
  have access to a storage location in Unity Catalog, the request fails and
  Unity Catalog does not attempt to authenticate to your cloud tenant on the
  user’s behalf.
  
  Databricks recommends using external locations rather than using storage
  credentials directly.
  
  To create external locations, you must be a metastore admin or a user with the
  CREATE EXTERNAL LOCATION privilege.`,
}

var createReq unitycatalog.CreateExternalLocation

func init() {
	Cmd.AddCommand(createCmd)
	// TODO: short flags

	createCmd.Flags().StringVar(&createReq.Comment, "comment", "", `[Create:OPT Update:OPT] User-provided free-form text description.`)
	createCmd.Flags().Int64Var(&createReq.CreatedAt, "created-at", 0, `[Create,Update:IGN] Time at which this External Location was created, in epoch milliseconds.`)
	createCmd.Flags().StringVar(&createReq.CreatedBy, "created-by", "", `[Create,Update:IGN] Username of External Location creator.`)
	createCmd.Flags().StringVar(&createReq.CredentialId, "credential-id", "", `[Create,Update:IGN] Unique ID of the location's Storage Credential.`)
	createCmd.Flags().StringVar(&createReq.CredentialName, "credential-name", "", `[Create:REQ Update:OPT] Current name of the Storage Credential this location uses.`)
	createCmd.Flags().StringVar(&createReq.MetastoreId, "metastore-id", "", `[Create,Update:IGN] Unique identifier of Metastore hosting the External Location.`)
	createCmd.Flags().StringVar(&createReq.Name, "name", "", `[Create:REQ Update:OPT] Name of the External Location.`)
	createCmd.Flags().StringVar(&createReq.Owner, "owner", "", `[Create:IGN Update:OPT] The owner of the External Location.`)
	createCmd.Flags().Int64Var(&createReq.UpdatedAt, "updated-at", 0, `[Create,Update:IGN] Time at which this was last modified, in epoch milliseconds.`)
	createCmd.Flags().StringVar(&createReq.UpdatedBy, "updated-by", "", `[Create,Update:IGN] Username of user who last modified the External Location.`)
	createCmd.Flags().StringVar(&createReq.Url, "url", "", `[Create:REQ Update:OPT] Path URL of the External Location.`)

}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: `Create an external location.`,
	Long: `Create an external location.
  
  Creates a new External Location entry in the Metastore. The caller must be a
  Metastore admin or have the CREATE EXTERNAL LOCATION privilege on the
  Metastore.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.ExternalLocations.Create(ctx, createReq)
		if err != nil {
			return err
		}

		pretty, err := ui.MarshalJSON(response)
		if err != nil {
			return err
		}
		cmd.OutOrStdout().Write(pretty)

		return nil
	},
}

var deleteReq unitycatalog.DeleteExternalLocationRequest

func init() {
	Cmd.AddCommand(deleteCmd)
	// TODO: short flags

	deleteCmd.Flags().BoolVar(&deleteReq.Force, "force", false, `Force deletion even if there are dependent external tables or mounts.`)
	deleteCmd.Flags().StringVar(&deleteReq.Name, "name", "", `Required.`)

}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: `Delete an external location.`,
	Long: `Delete an external location.
  
  Deletes the specified external location from the Metastore. The caller must be
  the owner of the external location.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.ExternalLocations.Delete(ctx, deleteReq)
		if err != nil {
			return err
		}

		return nil
	},
}

var getReq unitycatalog.GetExternalLocationRequest

func init() {
	Cmd.AddCommand(getCmd)
	// TODO: short flags

	getCmd.Flags().StringVar(&getReq.Name, "name", "", `Required.`)

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: `Get an external location.`,
	Long: `Get an external location.
  
  Gets an external location from the Metastore. The caller must be either a
  Metastore admin, the owner of the external location, or has an appropriate
  privilege level on the Metastore.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.ExternalLocations.Get(ctx, getReq)
		if err != nil {
			return err
		}

		pretty, err := ui.MarshalJSON(response)
		if err != nil {
			return err
		}
		cmd.OutOrStdout().Write(pretty)

		return nil
	},
}

func init() {
	Cmd.AddCommand(listCmd)

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: `List external locations.`,
	Long: `List external locations.
  
  Gets an array of External Locations (ExternalLocationInfo objects) from the
  Metastore. The caller must be a Metastore admin, is the owner of the external
  location, or has privileges to access the external location.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.ExternalLocations.ListAll(ctx)
		if err != nil {
			return err
		}

		pretty, err := ui.MarshalJSON(response)
		if err != nil {
			return err
		}
		cmd.OutOrStdout().Write(pretty)

		return nil
	},
}

var updateReq unitycatalog.UpdateExternalLocation

func init() {
	Cmd.AddCommand(updateCmd)
	// TODO: short flags

	updateCmd.Flags().StringVar(&updateReq.Comment, "comment", "", `[Create:OPT Update:OPT] User-provided free-form text description.`)
	updateCmd.Flags().Int64Var(&updateReq.CreatedAt, "created-at", 0, `[Create,Update:IGN] Time at which this External Location was created, in epoch milliseconds.`)
	updateCmd.Flags().StringVar(&updateReq.CreatedBy, "created-by", "", `[Create,Update:IGN] Username of External Location creator.`)
	updateCmd.Flags().StringVar(&updateReq.CredentialId, "credential-id", "", `[Create,Update:IGN] Unique ID of the location's Storage Credential.`)
	updateCmd.Flags().StringVar(&updateReq.CredentialName, "credential-name", "", `[Create:REQ Update:OPT] Current name of the Storage Credential this location uses.`)
	updateCmd.Flags().BoolVar(&updateReq.Force, "force", false, `TODO: SC-90063 re-add 'force' parameter in backward-compatible way for DBR (not removed below as it still works with CLI) Optional.`)
	updateCmd.Flags().StringVar(&updateReq.MetastoreId, "metastore-id", "", `[Create,Update:IGN] Unique identifier of Metastore hosting the External Location.`)
	updateCmd.Flags().StringVar(&updateReq.Name, "name", "", `[Create:REQ Update:OPT] Name of the External Location.`)
	updateCmd.Flags().StringVar(&updateReq.Owner, "owner", "", `[Create:IGN Update:OPT] The owner of the External Location.`)
	updateCmd.Flags().Int64Var(&updateReq.UpdatedAt, "updated-at", 0, `[Create,Update:IGN] Time at which this was last modified, in epoch milliseconds.`)
	updateCmd.Flags().StringVar(&updateReq.UpdatedBy, "updated-by", "", `[Create,Update:IGN] Username of user who last modified the External Location.`)
	updateCmd.Flags().StringVar(&updateReq.Url, "url", "", `[Create:REQ Update:OPT] Path URL of the External Location.`)

}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: `Update an external location.`,
	Long: `Update an external location.
  
  Updates an external location in the Metastore. The caller must be the owner of
  the externa location, or be a Metastore admin. In the second case, the admin
  can only update the name of the external location.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.ExternalLocations.Update(ctx, updateReq)
		if err != nil {
			return err
		}

		return nil
	},
}
