package ip_access_lists

import (
	"github.com/databricks/bricks/lib/ui"
	"github.com/databricks/bricks/project"
	"github.com/databricks/databricks-sdk-go/service/ipaccesslists"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "ip-access-lists",
	Short: `The IP Access List API enables Databricks admins to configure IP access lists for a workspace.`, // TODO: fix FirstSentence logic and append dot to summary
}

var createReq ipaccesslists.CreateIpAccessList

func init() {
	Cmd.AddCommand(createCmd)
	// TODO: short flags

	// TODO: complex arg: ip_addresses
	createCmd.Flags().StringVar(&createReq.Label, "label", "", `Label for the IP access list.`)
	// TODO: complex arg: list_type

}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: `Create access list Creates an IP access list for this workspace.`, // TODO: fix logic

	PreRunE: project.Configure, // TODO: improve logic for bundle/non-bundle invocations
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := project.Get(ctx).WorkspacesClient()
		response, err := w.IpAccessLists.Create(ctx, createReq)
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

var deleteReq ipaccesslists.Delete

func init() {
	Cmd.AddCommand(deleteCmd)
	// TODO: short flags

	deleteCmd.Flags().StringVar(&deleteReq.IpAccessListId, "ip-access-list-id", "", `The ID for the corresponding IP access list to modify.`)

}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: `Delete access list Deletes an IP access list, specified by its list ID.`, // TODO: fix logic

	PreRunE: project.Configure, // TODO: improve logic for bundle/non-bundle invocations
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := project.Get(ctx).WorkspacesClient()
		err := w.IpAccessLists.Delete(ctx, deleteReq)
		if err != nil {
			return err
		}

		return nil
	},
}

var getReq ipaccesslists.Get

func init() {
	Cmd.AddCommand(getCmd)
	// TODO: short flags

	getCmd.Flags().StringVar(&getReq.IpAccessListId, "ip-access-list-id", "", `The ID for the corresponding IP access list to modify.`)

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: `Get access list Gets an IP access list, specified by its list ID.`, // TODO: fix logic

	PreRunE: project.Configure, // TODO: improve logic for bundle/non-bundle invocations
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := project.Get(ctx).WorkspacesClient()
		response, err := w.IpAccessLists.Get(ctx, getReq)
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
	Short: `Get access lists Gets all IP access lists for the specified workspace.`, // TODO: fix logic

	PreRunE: project.Configure, // TODO: improve logic for bundle/non-bundle invocations
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := project.Get(ctx).WorkspacesClient()
		response, err := w.IpAccessLists.ListAll(ctx)
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

var replaceReq ipaccesslists.ReplaceIpAccessList

func init() {
	Cmd.AddCommand(replaceCmd)
	// TODO: short flags

	replaceCmd.Flags().BoolVar(&replaceReq.Enabled, "enabled", false, `Specifies whether this IP access list is enabled.`)
	replaceCmd.Flags().StringVar(&replaceReq.IpAccessListId, "ip-access-list-id", "", `The ID for the corresponding IP access list to modify.`)
	// TODO: complex arg: ip_addresses
	replaceCmd.Flags().StringVar(&replaceReq.Label, "label", "", `Label for the IP access list.`)
	replaceCmd.Flags().StringVar(&replaceReq.ListId, "list-id", "", `Universally unique identifier(UUID) of the IP access list.`)
	// TODO: complex arg: list_type

}

var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: `Replace access list Replaces an IP access list, specified by its ID.`, // TODO: fix logic

	PreRunE: project.Configure, // TODO: improve logic for bundle/non-bundle invocations
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := project.Get(ctx).WorkspacesClient()
		err := w.IpAccessLists.Replace(ctx, replaceReq)
		if err != nil {
			return err
		}

		return nil
	},
}

var updateReq ipaccesslists.UpdateIpAccessList

func init() {
	Cmd.AddCommand(updateCmd)
	// TODO: short flags

	updateCmd.Flags().BoolVar(&updateReq.Enabled, "enabled", false, `Specifies whether this IP access list is enabled.`)
	updateCmd.Flags().StringVar(&updateReq.IpAccessListId, "ip-access-list-id", "", `The ID for the corresponding IP access list to modify.`)
	// TODO: complex arg: ip_addresses
	updateCmd.Flags().StringVar(&updateReq.Label, "label", "", `Label for the IP access list.`)
	updateCmd.Flags().StringVar(&updateReq.ListId, "list-id", "", `Universally unique identifier(UUID) of the IP access list.`)
	// TODO: complex arg: list_type

}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: `Update access list Updates an existing IP access list, specified by its ID.`, // TODO: fix logic

	PreRunE: project.Configure, // TODO: improve logic for bundle/non-bundle invocations
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := project.Get(ctx).WorkspacesClient()
		err := w.IpAccessLists.Update(ctx, updateReq)
		if err != nil {
			return err
		}

		return nil
	},
}
