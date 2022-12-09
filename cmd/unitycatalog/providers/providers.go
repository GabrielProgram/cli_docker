package providers

import (
	"github.com/databricks/bricks/lib/sdk"
	"github.com/databricks/bricks/lib/ui"
	"github.com/databricks/databricks-sdk-go/service/unitycatalog"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "providers",
	Short: `Databricks Delta Sharing: Providers REST API.`,
	Long:  `Databricks Delta Sharing: Providers REST API`,
}

var createReq unitycatalog.CreateProvider

func init() {
	Cmd.AddCommand(createCmd)
	// TODO: short flags

	createCmd.Flags().BoolVar(&createReq.ActivatedByProvider, "activated-by-provider", false, `[Create,Update:IGN] Whether this provider is successfully activated by the data provider.`)
	createCmd.Flags().Var(&createReq.AuthenticationType, "authentication-type", `[Create:REQ,Update:IGN] The delta sharing authentication type.`)
	createCmd.Flags().StringVar(&createReq.Comment, "comment", "", `[Create,Update:OPT] Description about the provider.`)
	createCmd.Flags().Int64Var(&createReq.CreatedAt, "created-at", 0, `[Create,Update:IGN] Time at which this Provider was created, in epoch milliseconds.`)
	createCmd.Flags().StringVar(&createReq.CreatedBy, "created-by", "", `[Create,Update:IGN] Username of Provider creator.`)
	createCmd.Flags().StringVar(&createReq.Name, "name", "", `[Create,Update:REQ] The name of the Provider.`)
	// TODO: complex arg: recipient_profile
	createCmd.Flags().StringVar(&createReq.RecipientProfileStr, "recipient-profile-str", "", `[Create,Update:OPT] This field is only present when the authentication type is TOKEN.`)
	createCmd.Flags().StringVar(&createReq.SharingCode, "sharing-code", "", `[Create,Update:IGN] The server-generated one-time sharing code.`)
	createCmd.Flags().Int64Var(&createReq.UpdatedAt, "updated-at", 0, `[Create,Update:IGN] Time at which this Provider was created, in epoch milliseconds.`)
	createCmd.Flags().StringVar(&createReq.UpdatedBy, "updated-by", "", `[Create,Update:IGN] Username of user who last modified Share.`)

}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: `Create an auth provider.`,
	Long: `Create an auth provider.
  
  Creates a new authentication provider minimally based on a name and
  authentication type. The caller must be an admin on the Metastore.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Providers.Create(ctx, createReq)
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

var deleteReq unitycatalog.DeleteProviderRequest

func init() {
	Cmd.AddCommand(deleteCmd)
	// TODO: short flags

	deleteCmd.Flags().StringVar(&deleteReq.Name, "name", "", `Required.`)

}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: `Delete a provider.`,
	Long: `Delete a provider.
  
  Deletes an authentication provider, if the caller is a Metastore admin or is
  the owner of the provider.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.Providers.Delete(ctx, deleteReq)
		if err != nil {
			return err
		}

		return nil
	},
}

var getReq unitycatalog.GetProviderRequest

func init() {
	Cmd.AddCommand(getCmd)
	// TODO: short flags

	getCmd.Flags().StringVar(&getReq.Name, "name", "", `Required.`)

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: `Get a provider.`,
	Long: `Get a provider.
  
  Gets a specific authentication provider. The caller must supply the name of
  the provider, and must either be a Metastore admin or the owner of the
  provider.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Providers.Get(ctx, getReq)
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

var listReq unitycatalog.ListProvidersRequest

func init() {
	Cmd.AddCommand(listCmd)
	// TODO: short flags

	listCmd.Flags().StringVar(&listReq.DataProviderGlobalMetastoreId, "data-provider-global-metastore-id", "", `If not provided, all providers will be returned.`)

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: `List providers.`,
	Long: `List providers.
  
  Gets an array of available authentication providers. The caller must either be
  a Metastore admin or the owner of the providers. Providers not owned by the
  caller are not included in the response.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Providers.ListAll(ctx, listReq)
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

var listSharesReq unitycatalog.ListSharesRequest

func init() {
	Cmd.AddCommand(listSharesCmd)
	// TODO: short flags

	listSharesCmd.Flags().StringVar(&listSharesReq.Name, "name", "", `Required.`)

}

var listSharesCmd = &cobra.Command{
	Use:   "list-shares",
	Short: `List shares.`,
	Long: `List shares.
  
  Gets an array of all shares within the Metastore where:
  
  * the caller is a Metastore admin, or * the caller is the owner.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Providers.ListShares(ctx, listSharesReq)
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

var updateReq unitycatalog.UpdateProvider

func init() {
	Cmd.AddCommand(updateCmd)
	// TODO: short flags

	updateCmd.Flags().BoolVar(&updateReq.ActivatedByProvider, "activated-by-provider", false, `[Create,Update:IGN] Whether this provider is successfully activated by the data provider.`)
	updateCmd.Flags().Var(&updateReq.AuthenticationType, "authentication-type", `[Create:REQ,Update:IGN] The delta sharing authentication type.`)
	updateCmd.Flags().StringVar(&updateReq.Comment, "comment", "", `[Create,Update:OPT] Description about the provider.`)
	updateCmd.Flags().Int64Var(&updateReq.CreatedAt, "created-at", 0, `[Create,Update:IGN] Time at which this Provider was created, in epoch milliseconds.`)
	updateCmd.Flags().StringVar(&updateReq.CreatedBy, "created-by", "", `[Create,Update:IGN] Username of Provider creator.`)
	updateCmd.Flags().StringVar(&updateReq.Name, "name", "", `[Create, Update:REQ] The name of the Provider.`)
	// TODO: complex arg: recipient_profile
	updateCmd.Flags().StringVar(&updateReq.RecipientProfileStr, "recipient-profile-str", "", `[Create,Update:OPT] This field is only present when the authentication type is TOKEN.`)
	updateCmd.Flags().StringVar(&updateReq.SharingCode, "sharing-code", "", `[Create,Update:IGN] The server-generated one-time sharing code.`)
	updateCmd.Flags().Int64Var(&updateReq.UpdatedAt, "updated-at", 0, `[Create,Update:IGN] Time at which this Provider was created, in epoch milliseconds.`)
	updateCmd.Flags().StringVar(&updateReq.UpdatedBy, "updated-by", "", `[Create,Update:IGN] Username of user who last modified Share.`)

}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: `Update a provider.`,
	Long: `Update a provider.
  
  Updates the information for an authentication provider, if the caller is a
  Metastore admin or is the owner of the provider. If the update changes the
  provider name, the caller must be both a Metastore admin and the owner of the
  provider.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.Providers.Update(ctx, updateReq)
		if err != nil {
			return err
		}

		return nil
	},
}

// end service Providers
