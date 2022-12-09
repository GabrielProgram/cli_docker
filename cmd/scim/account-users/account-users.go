package account_users

import (
	"github.com/databricks/bricks/lib/jsonflag"
	"github.com/databricks/bricks/lib/sdk"
	"github.com/databricks/bricks/lib/ui"
	"github.com/databricks/databricks-sdk-go/service/scim"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "account-users",
	Short: `User identities recognized by Databricks and represented by email addresses.`,
	Long: `User identities recognized by Databricks and represented by email addresses.
  
  Databricks recommends using SCIM provisioning to sync users and groups
  automatically from your identity provider to your Databricks Account. SCIM
  streamlines onboarding a new employee or team by using your identity provider
  to create users and groups in Databricks Account and give them the proper
  level of access. When a user leaves your organization or no longer needs
  access to Databricks Account, admins can terminate the user in your identity
  provider and that user’s account will also be removed from Databricks
  Account. This ensures a consistent offboarding process and prevents
  unauthorized users from accessing sensitive data.`,
}

// start create command

var createReq scim.User
var createJson jsonflag.JsonFlag

func init() {
	Cmd.AddCommand(createCmd)
	// TODO: short flags
	createCmd.Flags().Var(&createJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	createCmd.Flags().BoolVar(&createReq.Active, "active", createReq.Active, `If this user is active.`)
	createCmd.Flags().StringVar(&createReq.DisplayName, "display-name", createReq.DisplayName, `String that represents a concatenation of given and family names.`)
	// TODO: array: emails
	// TODO: array: entitlements
	createCmd.Flags().StringVar(&createReq.ExternalId, "external-id", createReq.ExternalId, ``)
	// TODO: array: groups
	createCmd.Flags().StringVar(&createReq.Id, "id", createReq.Id, `Databricks user ID.`)
	// TODO: complex arg: name
	// TODO: array: roles
	createCmd.Flags().StringVar(&createReq.UserName, "user-name", createReq.UserName, `Email address of the Databricks user.`)

}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: `Create a new user.`,
	Long: `Create a new user.
  
  Creates a new user in the Databricks Account. This new user will also be added
  to the Databricks account.`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = createJson.Unmarshall(&createReq)
		if err != nil {
			return err
		}
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		response, err := a.Users.Create(ctx, createReq)
		if err != nil {
			return err
		}
		return ui.Render(cmd, response)
	},
}

// start delete command

var deleteReq scim.DeleteUserRequest

func init() {
	Cmd.AddCommand(deleteCmd)
	// TODO: short flags

	deleteCmd.Flags().StringVar(&deleteReq.Id, "id", deleteReq.Id, `Unique ID for a user in the Databricks Account.`)

}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: `Delete a user.`,
	Long: `Delete a user.
  
  Deletes a user. Deleting a user from a Databricks Account also removes objects
  associated with the user.`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		err = a.Users.Delete(ctx, deleteReq)
		if err != nil {
			return err
		}
		return nil
	},
}

// start get command

var getReq scim.GetUserRequest

func init() {
	Cmd.AddCommand(getCmd)
	// TODO: short flags

	getCmd.Flags().StringVar(&getReq.Id, "id", getReq.Id, `Unique ID for a user in the Databricks Account.`)

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: `Get user details.`,
	Long: `Get user details.
  
  Gets information for a specific user in Databricks Account.`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		response, err := a.Users.Get(ctx, getReq)
		if err != nil {
			return err
		}
		return ui.Render(cmd, response)
	},
}

// start list command

var listReq scim.ListUsersRequest
var listJson jsonflag.JsonFlag

func init() {
	Cmd.AddCommand(listCmd)
	// TODO: short flags
	listCmd.Flags().Var(&listJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	listCmd.Flags().StringVar(&listReq.Attributes, "attributes", listReq.Attributes, `Comma-separated list of attributes to return in response.`)
	listCmd.Flags().IntVar(&listReq.Count, "count", listReq.Count, `Desired number of results per page.`)
	listCmd.Flags().StringVar(&listReq.ExcludedAttributes, "excluded-attributes", listReq.ExcludedAttributes, `Comma-separated list of attributes to exclude in response.`)
	listCmd.Flags().StringVar(&listReq.Filter, "filter", listReq.Filter, `Query by which the results have to be filtered.`)
	listCmd.Flags().StringVar(&listReq.SortBy, "sort-by", listReq.SortBy, `Attribute to sort the results.`)
	listCmd.Flags().Var(&listReq.SortOrder, "sort-order", `The order to sort the results.`)
	listCmd.Flags().IntVar(&listReq.StartIndex, "start-index", listReq.StartIndex, `Specifies the index of the first result.`)

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: `List users.`,
	Long: `List users.
  
  Gets details for all the users associated with a Databricks Account.`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = listJson.Unmarshall(&listReq)
		if err != nil {
			return err
		}
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		response, err := a.Users.ListAll(ctx, listReq)
		if err != nil {
			return err
		}
		return ui.Render(cmd, response)
	},
}

// start patch command

var patchReq scim.PartialUpdate
var patchJson jsonflag.JsonFlag

func init() {
	Cmd.AddCommand(patchCmd)
	// TODO: short flags
	patchCmd.Flags().Var(&patchJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	patchCmd.Flags().StringVar(&patchReq.Id, "id", patchReq.Id, `Unique ID for a group in the Databricks Account.`)
	// TODO: array: operations

}

var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: `Update user details.`,
	Long: `Update user details.
  
  Partially updates a user resource by applying the supplied operations on
  specific user attributes.`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = patchJson.Unmarshall(&patchReq)
		if err != nil {
			return err
		}
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		err = a.Users.Patch(ctx, patchReq)
		if err != nil {
			return err
		}
		return nil
	},
}

// start update command

var updateReq scim.User
var updateJson jsonflag.JsonFlag

func init() {
	Cmd.AddCommand(updateCmd)
	// TODO: short flags
	updateCmd.Flags().Var(&updateJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	updateCmd.Flags().BoolVar(&updateReq.Active, "active", updateReq.Active, `If this user is active.`)
	updateCmd.Flags().StringVar(&updateReq.DisplayName, "display-name", updateReq.DisplayName, `String that represents a concatenation of given and family names.`)
	// TODO: array: emails
	// TODO: array: entitlements
	updateCmd.Flags().StringVar(&updateReq.ExternalId, "external-id", updateReq.ExternalId, ``)
	// TODO: array: groups
	updateCmd.Flags().StringVar(&updateReq.Id, "id", updateReq.Id, `Databricks user ID.`)
	// TODO: complex arg: name
	// TODO: array: roles
	updateCmd.Flags().StringVar(&updateReq.UserName, "user-name", updateReq.UserName, `Email address of the Databricks user.`)

}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: `Replace a user.`,
	Long: `Replace a user.
  
  Replaces a user's information with the data supplied in request.`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = updateJson.Unmarshall(&updateReq)
		if err != nil {
			return err
		}
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		err = a.Users.Update(ctx, updateReq)
		if err != nil {
			return err
		}
		return nil
	},
}

// end service AccountUsers
