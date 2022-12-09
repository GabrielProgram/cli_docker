package recipients

import (
	"github.com/databricks/bricks/lib/sdk"
	"github.com/databricks/bricks/lib/ui"
	"github.com/databricks/databricks-sdk-go/service/unitycatalog"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "recipients",
	Short: `Databricks Delta Sharing: Recipients REST API.`,
	Long:  `Databricks Delta Sharing: Recipients REST API`,
}

var createReq unitycatalog.CreateRecipient

func init() {
	Cmd.AddCommand(createCmd)
	// TODO: short flags

	createCmd.Flags().Var(&createReq.AuthenticationType, "authentication-type", `The delta sharing authentication type.`)
	createCmd.Flags().StringVar(&createReq.Comment, "comment", createReq.Comment, `Description about the recipient.`)
	// TODO: complex arg: ip_access_list
	createCmd.Flags().StringVar(&createReq.Name, "name", createReq.Name, `Name of Recipient.`)
	createCmd.Flags().StringVar(&createReq.SharingCode, "sharing-code", createReq.SharingCode, `The one-time sharing code provided by the data recipient.`)

}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: `Create a share recipient.`,
	Long: `Create a share recipient.
  
  Creates a new recipient with the delta sharing authentication type in the
  Metastore. The caller must be a Metastore admin or has the CREATE RECIPIENT
  privilege on the Metastore.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Recipients.Create(ctx, createReq)
		if err != nil {
			return err
		}
		return ui.Render(cmd, response)
	},
}

var deleteReq unitycatalog.DeleteRecipientRequest

func init() {
	Cmd.AddCommand(deleteCmd)
	// TODO: short flags

	deleteCmd.Flags().StringVar(&deleteReq.Name, "name", deleteReq.Name, `Required.`)

}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: `Delete a share recipient.`,
	Long: `Delete a share recipient.
  
  Deletes the specified recipient from the Metastore. The caller must be the
  owner of the recipient.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.Recipients.Delete(ctx, deleteReq)
		if err != nil {
			return err
		}
		return nil
	},
}

var getReq unitycatalog.GetRecipientRequest

func init() {
	Cmd.AddCommand(getCmd)
	// TODO: short flags

	getCmd.Flags().StringVar(&getReq.Name, "name", getReq.Name, `Required.`)

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: `Get a share recipient.`,
	Long: `Get a share recipient.
  
  Gets a share recipient from the Metastore if:
  
  * the caller is the owner of the share recipient, or: * is a Metastore admin`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Recipients.Get(ctx, getReq)
		if err != nil {
			return err
		}
		return ui.Render(cmd, response)
	},
}

var listReq unitycatalog.ListRecipientsRequest

func init() {
	Cmd.AddCommand(listCmd)
	// TODO: short flags

	listCmd.Flags().StringVar(&listReq.DataRecipientGlobalMetastoreId, "data-recipient-global-metastore-id", listReq.DataRecipientGlobalMetastoreId, `If not provided, all recipients will be returned.`)

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: `List share recipients.`,
	Long: `List share recipients.
  
  Gets an array of all share recipients within the current Metastore where:
  
  * the caller is a Metastore admin, or * the caller is the owner.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Recipients.ListAll(ctx, listReq)
		if err != nil {
			return err
		}
		return ui.Render(cmd, response)
	},
}

var rotateTokenReq unitycatalog.RotateRecipientToken

func init() {
	Cmd.AddCommand(rotateTokenCmd)
	// TODO: short flags

	rotateTokenCmd.Flags().Int64Var(&rotateTokenReq.ExistingTokenExpireInSeconds, "existing-token-expire-in-seconds", rotateTokenReq.ExistingTokenExpireInSeconds, `Required.`)
	rotateTokenCmd.Flags().StringVar(&rotateTokenReq.Name, "name", rotateTokenReq.Name, `Required.`)

}

var rotateTokenCmd = &cobra.Command{
	Use:   "rotate-token",
	Short: `Rotate a token.`,
	Long: `Rotate a token.
  
  Refreshes the specified recipient's delta sharing authentication token with
  the provided token info. The caller must be the owner of the recipient.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Recipients.RotateToken(ctx, rotateTokenReq)
		if err != nil {
			return err
		}
		return ui.Render(cmd, response)
	},
}

var sharePermissionsReq unitycatalog.SharePermissionsRequest

func init() {
	Cmd.AddCommand(sharePermissionsCmd)
	// TODO: short flags

	sharePermissionsCmd.Flags().StringVar(&sharePermissionsReq.Name, "name", sharePermissionsReq.Name, `Required.`)

}

var sharePermissionsCmd = &cobra.Command{
	Use:   "share-permissions",
	Short: `Get share permissions.`,
	Long: `Get share permissions.
  
  Gets the share permissions for the specified Recipient. The caller must be a
  Metastore admin or the owner of the Recipient.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Recipients.SharePermissions(ctx, sharePermissionsReq)
		if err != nil {
			return err
		}
		return ui.Render(cmd, response)
	},
}

var updateReq unitycatalog.UpdateRecipient

func init() {
	Cmd.AddCommand(updateCmd)
	// TODO: short flags

	updateCmd.Flags().Var(&updateReq.AuthenticationType, "authentication-type", `The delta sharing authentication type.`)
	updateCmd.Flags().StringVar(&updateReq.Comment, "comment", updateReq.Comment, `Description about the recipient.`)
	// TODO: complex arg: ip_access_list
	updateCmd.Flags().StringVar(&updateReq.Name, "name", updateReq.Name, `Name of Recipient.`)

}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: `Update a share recipient.`,
	Long: `Update a share recipient.
  
  Updates an existing recipient in the Metastore. The caller must be a Metastore
  admin or the owner of the recipient. If the recipient name will be updated,
  the user must be both a Metastore admin and the owner of the recipient.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.Recipients.Update(ctx, updateReq)
		if err != nil {
			return err
		}
		return nil
	},
}

// end service Recipients
