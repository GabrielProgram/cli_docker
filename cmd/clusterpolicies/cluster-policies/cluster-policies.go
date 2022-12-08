package cluster_policies

import (
	"github.com/databricks/bricks/lib/sdk"
	"github.com/databricks/bricks/lib/ui"
	"github.com/databricks/databricks-sdk-go/service/clusterpolicies"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "cluster-policies",
	Short: `Cluster policy limits the ability to configure clusters based on a set of rules.`,
}

var createReq clusterpolicies.CreatePolicy

func init() {
	Cmd.AddCommand(createCmd)
	// TODO: short flags

	createCmd.Flags().StringVar(&createReq.Definition, "definition", "", `Policy definition document expressed in Databricks Cluster Policy Definition Language.`)
	createCmd.Flags().StringVar(&createReq.Name, "name", "", `Cluster Policy name requested by the user.`)

}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: `Create a new policy.`,
	Long: `Create a new policy.
  
  Creates a new policy with prescribed settings.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.ClusterPolicies.Create(ctx, createReq)
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

var deleteReq clusterpolicies.DeletePolicy

func init() {
	Cmd.AddCommand(deleteCmd)
	// TODO: short flags

	deleteCmd.Flags().StringVar(&deleteReq.PolicyId, "policy-id", "", `The ID of the policy to delete.`)

}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: `Delete a cluster policy.`,
	Long: `Delete a cluster policy.
  
  Delete a policy for a cluster. Clusters governed by this policy can still run,
  but cannot be edited.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.ClusterPolicies.Delete(ctx, deleteReq)
		if err != nil {
			return err
		}

		return nil
	},
}

var editReq clusterpolicies.EditPolicy

func init() {
	Cmd.AddCommand(editCmd)
	// TODO: short flags

	editCmd.Flags().StringVar(&editReq.Definition, "definition", "", `Policy definition document expressed in Databricks Cluster Policy Definition Language.`)
	editCmd.Flags().StringVar(&editReq.Name, "name", "", `Cluster Policy name requested by the user.`)
	editCmd.Flags().StringVar(&editReq.PolicyId, "policy-id", "", `The ID of the policy to update.`)

}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: `Update a cluster policy.`,
	Long: `Update a cluster policy.
  
  Update an existing policy for cluster. This operation may make some clusters
  governed by the previous policy invalid.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.ClusterPolicies.Edit(ctx, editReq)
		if err != nil {
			return err
		}

		return nil
	},
}

var getReq clusterpolicies.Get

func init() {
	Cmd.AddCommand(getCmd)
	// TODO: short flags

	getCmd.Flags().StringVar(&getReq.PolicyId, "policy-id", "", `Canonical unique identifier for the cluster policy.`)

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: `Get entity.`,
	Long: `Get entity.
  
  Get a cluster policy entity. Creation and editing is available to admins only.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.ClusterPolicies.Get(ctx, getReq)
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
	Short: `Get a cluster policy.`,
	Long: `Get a cluster policy.
  
  Returns a list of policies accessible by the requesting user.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.ClusterPolicies.ListAll(ctx)
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
