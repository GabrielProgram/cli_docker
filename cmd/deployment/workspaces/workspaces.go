package workspaces

import (
	"github.com/databricks/bricks/lib/sdk"
	"github.com/databricks/bricks/lib/ui"
	"github.com/databricks/databricks-sdk-go/service/deployment"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "workspaces",
	Short: `These APIs manage workspaces for this account.`,
}

var createReq deployment.CreateWorkspaceRequest

func init() {
	Cmd.AddCommand(createCmd)
	// TODO: short flags

	createCmd.Flags().StringVar(&createReq.AwsRegion, "aws-region", "", `The AWS region of the workspace's data plane.`)
	createCmd.Flags().StringVar(&createReq.Cloud, "cloud", "", `The cloud provider which the workspace uses.`)
	// TODO: complex arg: cloud_resource_bucket
	createCmd.Flags().StringVar(&createReq.CredentialsId, "credentials-id", "", `ID of the workspace's credential configuration object.`)
	createCmd.Flags().StringVar(&createReq.DeploymentName, "deployment-name", "", `The deployment name defines part of the subdomain for the workspace.`)
	createCmd.Flags().StringVar(&createReq.Location, "location", "", `The Google Cloud region of the workspace data plane in your Google account.`)
	createCmd.Flags().StringVar(&createReq.ManagedServicesCustomerManagedKeyId, "managed-services-customer-managed-key-id", "", `The ID of the workspace's managed services encryption key configuration object.`)
	// TODO: complex arg: network
	createCmd.Flags().StringVar(&createReq.NetworkId, "network-id", "", `The ID of the workspace's network configuration object.`)
	createCmd.Flags().Var(&createReq.PricingTier, "pricing-tier", `The pricing tier of the workspace.`)
	createCmd.Flags().StringVar(&createReq.PrivateAccessSettingsId, "private-access-settings-id", "", `ID of the workspace's private access settings object.`)
	createCmd.Flags().StringVar(&createReq.StorageConfigurationId, "storage-configuration-id", "", `The ID of the workspace's storage configuration object.`)
	createCmd.Flags().StringVar(&createReq.StorageCustomerManagedKeyId, "storage-customer-managed-key-id", "", `The ID of the workspace's storage encryption key configuration object.`)
	createCmd.Flags().StringVar(&createReq.WorkspaceName, "workspace-name", "", `The workspace's human-readable name.`)

}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: `Create a new workspace.`,
	Long: `Create a new workspace.
  
  Creates a new workspace using a credential configuration and a storage
  configuration, an optional network configuration (if using a customer-managed
  VPC), an optional managed services key configuration (if using
  customer-managed keys for managed services), and an optional storage key
  configuration (if using customer-managed keys for storage). The key
  configurations used for managed services and storage encryption can be the
  same or different.
  
  **Important**: This operation is asynchronous. A response with HTTP status
  code 200 means the request has been accepted and is in progress, but does not
  mean that the workspace deployed successfully and is running. The initial
  workspace status is typically PROVISIONING. Use the workspace ID
  (workspace_id) field in the response to identify the new workspace and make
  repeated GET requests with the workspace ID and check its status. The
  workspace becomes available when the status changes to RUNNING.
  
  You can share one customer-managed VPC with multiple workspaces in a single
  account. It is not required to create a new VPC for each workspace. However,
  you **cannot** reuse subnets or Security Groups between workspaces. If you
  plan to share one VPC with multiple workspaces, make sure you size your VPC
  and subnets accordingly. Because a Databricks Account API network
  configuration encapsulates this information, you cannot reuse a Databricks
  Account API network configuration across workspaces.\nFor information about
  how to create a new workspace with this API **including error handling**, see
  [Create a new workspace using the Account API].
  
  **Important**: Customer-managed VPCs, PrivateLink, and customer-managed keys
  are supported on a limited set of deployment and subscription types. If you
  have questions about availability, contact your Databricks
  representative.\n\nThis operation is available only if your account is on the
  E2 version of the platform or on a select custom plan that allows multiple
  workspaces per account.
  
  [Create a new workspace using the Account API]: http://docs.databricks.com/administration-guide/account-api/new-workspace.html`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		response, err := a.Workspaces.Create(ctx, createReq)
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

var deleteReq deployment.DeleteWorkspaceRequest

func init() {
	Cmd.AddCommand(deleteCmd)
	// TODO: short flags

	deleteCmd.Flags().Int64Var(&deleteReq.WorkspaceId, "workspace-id", 0, `Workspace ID.`)

}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: `Delete workspace.`,
	Long: `Delete workspace.
  
  Terminates and deletes a Databricks workspace. From an API perspective,
  deletion is immediate. However, it might take a few minutes for all workspaces
  resources to be deleted, depending on the size and number of workspace
  resources.
  
  This operation is available only if your account is on the E2 version of the
  platform or on a select custom plan that allows multiple workspaces per
  account.`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		err := a.Workspaces.Delete(ctx, deleteReq)
		if err != nil {
			return err
		}

		return nil
	},
}

var getReq deployment.GetWorkspaceRequest

func init() {
	Cmd.AddCommand(getCmd)
	// TODO: short flags

	getCmd.Flags().Int64Var(&getReq.WorkspaceId, "workspace-id", 0, `Workspace ID.`)

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: `Get workspace.`,
	Long: `Get workspace.
  
  Gets information including status for a Databricks workspace, specified by ID.
  In the response, the workspace_status field indicates the current status.
  After initial workspace creation (which is asynchronous), make repeated GET
  requests with the workspace ID and check its status. The workspace becomes
  available when the status changes to RUNNING.
  
  For information about how to create a new workspace with this API **including
  error handling**, see [Create a new workspace using the Account API].
  
  This operation is available only if your account is on the E2 version of the
  platform or on a select custom plan that allows multiple workspaces per
  account.
  
  [Create a new workspace using the Account API]: http://docs.databricks.com/administration-guide/account-api/new-workspace.html`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		response, err := a.Workspaces.Get(ctx, getReq)
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

var getWorkspaceKeyHistoryReq deployment.GetWorkspaceKeyHistoryRequest

func init() {
	Cmd.AddCommand(getWorkspaceKeyHistoryCmd)
	// TODO: short flags

	getWorkspaceKeyHistoryCmd.Flags().Int64Var(&getWorkspaceKeyHistoryReq.WorkspaceId, "workspace-id", 0, `Workspace ID.`)

}

var getWorkspaceKeyHistoryCmd = &cobra.Command{
	Use:   "get-workspace-key-history",
	Short: `Get the history of a workspace's associations with keys.`,
	Long: `Get the history of a workspace's associations with keys.
  
  Gets a list of all associations with key configuration objects for the
  specified workspace that encapsulate customer-managed keys that encrypt
  managed services, workspace storage, or in some cases both.
  
  **Important**: In the current implementation, keys cannot be rotated or
  removed from a workspace. It is possible for a workspace to show a storage
  customer-managed key having been attached and then detached if the workspace
  was updated to use the key and the update operation failed.
  
  **Important**: Customer-managed keys are supported only for some deployment
  types and subscription types.
  
  This operation is available only if your account is on the E2 version of the
  platform.`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		response, err := a.Workspaces.GetWorkspaceKeyHistory(ctx, getWorkspaceKeyHistoryReq)
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
	Short: `Get all workspaces.`,
	Long: `Get all workspaces.
  
  Gets a list of all workspaces associated with an account, specified by ID.
  
  This operation is available only if your account is on the E2 version of the
  platform or on a select custom plan that allows multiple workspaces per
  account.`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		response, err := a.Workspaces.List(ctx)
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

var updateReq deployment.UpdateWorkspaceRequest

func init() {
	Cmd.AddCommand(updateCmd)
	// TODO: short flags

	updateCmd.Flags().StringVar(&updateReq.AwsRegion, "aws-region", "", `The AWS region of the workspace's data plane (for example, us-west-2).`)
	updateCmd.Flags().StringVar(&updateReq.CredentialsId, "credentials-id", "", `ID of the workspace's credential configuration object.`)
	updateCmd.Flags().StringVar(&updateReq.ManagedServicesCustomerManagedKeyId, "managed-services-customer-managed-key-id", "", `The ID of the workspace's managed services encryption key configuration object.`)
	updateCmd.Flags().StringVar(&updateReq.NetworkId, "network-id", "", `The ID of the workspace's network configuration object.`)
	updateCmd.Flags().StringVar(&updateReq.StorageConfigurationId, "storage-configuration-id", "", `The ID of the workspace's storage configuration object.`)
	updateCmd.Flags().StringVar(&updateReq.StorageCustomerManagedKeyId, "storage-customer-managed-key-id", "", `The ID of the key configuration object for workspace storage.`)
	updateCmd.Flags().Int64Var(&updateReq.WorkspaceId, "workspace-id", 0, `Workspace ID.`)

}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: `Update workspace configuration.`,
	Long: `Update workspace configuration.
  
  Updates a workspace configuration for either a running workspace or a failed
  workspace. The elements that can be updated varies between these two use
  cases.
  
  ### Update a failed workspace You can update a Databricks workspace
  configuration for failed workspace deployment for some fields, but not all
  fields. For a failed workspace, this request supports updates to the following
  fields only: - Credential configuration ID - Storage configuration ID -
  Network configuration ID. Used only if you use customer-managed VPC. - Key
  configuration ID for managed services (control plane storage, such as notebook
  source and Databricks SQL queries). Used only if you use customer-managed keys
  for managed services. - Key configuration ID for workspace storage (root S3
  bucket and, optionally, EBS volumes). Used only if you use customer-managed
  keys for workspace storage. **Important**: If the workspace was ever in the
  running state, even if briefly before becoming a failed workspace, you cannot
  add a new key configuration ID for workspace storage.
  
  After calling the PATCH operation to update the workspace configuration,
  make repeated GET requests with the workspace ID and check the workspace
  status. The workspace is successful if the status changes to RUNNING.
  
  For information about how to create a new workspace with this API **including
  error handling**, see [Create a new workspace using the Account API].
  
  ### Update a running workspace You can update a Databricks workspace
  configuration for running workspaces for some fields, but not all fields. For
  a running workspace, this request supports updating the following fields only:
  - Credential configuration ID
  
  - Network configuration ID. Used only if you already use use customer-managed
  VPC. This change is supported only if you specified a network configuration ID
  in your original workspace creation. In other words, you cannot switch from a
  Databricks-managed VPC to a customer-managed VPC. **Note**: You cannot use a
  network configuration update in this API to add support for PrivateLink (in
  Public Preview). To add PrivateLink to an existing workspace, contact your
  Databricks representative.
  
  - Key configuration ID for managed services (control plane storage, such as
  notebook source and Databricks SQL queries). Databricks does not directly
  encrypt the data with the customer-managed key (CMK). Databricks uses both the
  CMK and the Databricks managed key (DMK) that is unique to your workspace to
  encrypt the Data Encryption Key (DEK). Databricks uses the DEK to encrypt your
  workspace's managed services persisted data. If the workspace does not already
  have a CMK for managed services, adding this ID enables managed services
  encryption for new or updated data. Existing managed services data that
  existed before adding the key remains not encrypted with the DEK until it is
  modified. If the workspace already has customer-managed keys for managed
  services, this request rotates (changes) the CMK keys and the DEK is
  re-encrypted with the DMK and the new CMK. - Key configuration ID for
  workspace storage (root S3 bucket and, optionally, EBS volumes). You can set
  this only if the workspace does not already have a customer-managed key
  configuration for workspace storage.
  
  **Important**: For updating running workspaces, this API is unavailable on
  Mondays, Tuesdays, and Thursdays from 4:30pm-7:30pm PST due to routine
  maintenance. Plan your workspace updates accordingly. For questions about this
  schedule, contact your Databricks representative.
  
  **Important**: To update a running workspace, your workspace must have no
  running cluster instances, which includes all-purpose clusters, job clusters,
  and pools that might have running clusters. Terminate all cluster instances in
  the workspace before calling this API.
  
  ### Wait until changes take effect. After calling the PATCH operation to
  update the workspace configuration, make repeated GET requests with the
  workspace ID and check the workspace status and the status of the fields. *
  For workspaces with a Databricks-managed VPC, the workspace status becomes
  PROVISIONING temporarily (typically under 20 minutes). If the workspace
  update is successful, the workspace status changes to RUNNING. Note that you
  can also check the workspace status in the [Account Console]. However, you
  cannot use or create clusters for another 20 minutes after that status change.
  This results in a total of up to 40 minutes in which you cannot create
  clusters. If you create or use clusters before this time interval elapses,
  clusters do not launch successfully, fail, or could cause other unexpected
  behavior.
  
  * For workspaces with a customer-managed VPC, the workspace status stays at
  status RUNNING and the VPC change happens immediately. A change to the
  storage customer-managed key configuration ID might take a few minutes to
  update, so continue to check the workspace until you observe that it has been
  updated. If the update fails, the workspace might revert silently to its
  original configuration. After the workspace has been updated, you cannot use
  or create clusters for another 20 minutes. If you create or use clusters
  before this time interval elapses, clusters do not launch successfully, fail,
  or could cause other unexpected behavior.
  
  If you update the _storage_ customer-managed key configurations, it takes 20
  minutes for the changes to fully take effect. During the 20 minute wait, it is
  important that you stop all REST API calls to the DBFS API. If you are
  modifying _only the managed services key configuration_, you can omit the 20
  minute wait.
  
  **Important**: Customer-managed keys and customer-managed VPCs are supported
  by only some deployment types and subscription types. If you have questions
  about availability, contact your Databricks representative.
  
  This operation is available only if your account is on the E2 version of the
  platform or on a select custom plan that allows multiple workspaces per
  account.
  
  [Account Console]: https://docs.databricks.com/administration-guide/account-settings-e2/account-console-e2.html
  [Create a new workspace using the Account API]: http://docs.databricks.com/administration-guide/account-api/new-workspace.html`,

	PreRunE: sdk.PreAccountClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		a := sdk.AccountClient(ctx)
		err := a.Workspaces.Update(ctx, updateReq)
		if err != nil {
			return err
		}

		return nil
	},
}
