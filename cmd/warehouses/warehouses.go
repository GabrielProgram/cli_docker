package warehouses

import (
	query_history "github.com/databricks/bricks/cmd/warehouses/query-history"
	"github.com/databricks/bricks/lib/sdk"
	"github.com/databricks/bricks/lib/ui"
	"github.com/databricks/databricks-sdk-go/service/warehouses"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "warehouses",
	Short: `A SQL warehouse is a compute resource that lets you run SQL commands on data objects within Databricks SQL.`,
	Long: `A SQL warehouse is a compute resource that lets you run SQL commands on data
  objects within Databricks SQL. Compute resources are infrastructure resources
  that provide processing capabilities in the cloud.`,
}

var createWarehouseReq warehouses.CreateWarehouseRequest

func init() {
	Cmd.AddCommand(createWarehouseCmd)
	// TODO: short flags

	createWarehouseCmd.Flags().IntVar(&createWarehouseReq.AutoStopMins, "auto-stop-mins", 0, `The amount of time in minutes that a SQL Endpoint must be idle (i.e., no RUNNING queries) before it is automatically stopped.`)
	// TODO: complex arg: channel
	createWarehouseCmd.Flags().StringVar(&createWarehouseReq.ClusterSize, "cluster-size", "", `Size of the clusters allocated for this endpoint.`)
	createWarehouseCmd.Flags().StringVar(&createWarehouseReq.CreatorName, "creator-name", "", `endpoint creator name.`)
	createWarehouseCmd.Flags().BoolVar(&createWarehouseReq.EnablePhoton, "enable-photon", false, `Configures whether the endpoint should use Photon optimized clusters.`)
	createWarehouseCmd.Flags().BoolVar(&createWarehouseReq.EnableServerlessCompute, "enable-serverless-compute", false, `Configures whether the endpoint should use Serverless Compute (aka Nephos) Defaults to value in global endpoint settings.`)
	createWarehouseCmd.Flags().StringVar(&createWarehouseReq.InstanceProfileArn, "instance-profile-arn", "", `Deprecated.`)
	createWarehouseCmd.Flags().IntVar(&createWarehouseReq.MaxNumClusters, "max-num-clusters", 0, `Maximum number of clusters that the autoscaler will create to handle concurrent queries.`)
	createWarehouseCmd.Flags().IntVar(&createWarehouseReq.MinNumClusters, "min-num-clusters", 0, `Minimum number of available clusters that will be maintained for this SQL Endpoint.`)
	createWarehouseCmd.Flags().StringVar(&createWarehouseReq.Name, "name", "", `Logical name for the cluster.`)
	createWarehouseCmd.Flags().Var(&createWarehouseReq.SpotInstancePolicy, "spot-instance-policy", `Configurations whether the endpoint should use spot instances.`)
	// TODO: complex arg: tags
	createWarehouseCmd.Flags().Var(&createWarehouseReq.WarehouseType, "warehouse-type", `Warehouse type (Classic/Pro).`)

}

var createWarehouseCmd = &cobra.Command{
	Use:   "create-warehouse",
	Short: `Create a warehouse.`,
	Long: `Create a warehouse.
  
  Creates a new SQL warehouse.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Warehouses.CreateWarehouse(ctx, createWarehouseReq)
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

var deleteWarehouseReq warehouses.DeleteWarehouse

func init() {
	Cmd.AddCommand(deleteWarehouseCmd)
	// TODO: short flags

	deleteWarehouseCmd.Flags().StringVar(&deleteWarehouseReq.Id, "id", "", `Required.`)

}

var deleteWarehouseCmd = &cobra.Command{
	Use:   "delete-warehouse",
	Short: `Delete a warehouse.`,
	Long: `Delete a warehouse.
  
  Deletes a SQL warehouse.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.Warehouses.DeleteWarehouse(ctx, deleteWarehouseReq)
		if err != nil {
			return err
		}

		return nil
	},
}

var editWarehouseReq warehouses.EditWarehouseRequest

func init() {
	Cmd.AddCommand(editWarehouseCmd)
	// TODO: short flags

	editWarehouseCmd.Flags().IntVar(&editWarehouseReq.AutoStopMins, "auto-stop-mins", 0, `The amount of time in minutes that a SQL Endpoint must be idle (i.e., no RUNNING queries) before it is automatically stopped.`)
	// TODO: complex arg: channel
	editWarehouseCmd.Flags().StringVar(&editWarehouseReq.ClusterSize, "cluster-size", "", `Size of the clusters allocated for this endpoint.`)
	editWarehouseCmd.Flags().StringVar(&editWarehouseReq.CreatorName, "creator-name", "", `endpoint creator name.`)
	editWarehouseCmd.Flags().BoolVar(&editWarehouseReq.EnableDatabricksCompute, "enable-databricks-compute", false, `Configures whether the endpoint should use Databricks Compute (aka Nephos) Deprecated: Use enable_serverless_compute TODO(SC-79930): Remove the field once clients are updated.`)
	editWarehouseCmd.Flags().BoolVar(&editWarehouseReq.EnablePhoton, "enable-photon", false, `Configures whether the endpoint should use Photon optimized clusters.`)
	editWarehouseCmd.Flags().BoolVar(&editWarehouseReq.EnableServerlessCompute, "enable-serverless-compute", false, `Configures whether the endpoint should use Serverless Compute (aka Nephos) Defaults to value in global endpoint settings.`)
	editWarehouseCmd.Flags().StringVar(&editWarehouseReq.Id, "id", "", `Required.`)
	editWarehouseCmd.Flags().StringVar(&editWarehouseReq.InstanceProfileArn, "instance-profile-arn", "", `Deprecated.`)
	editWarehouseCmd.Flags().IntVar(&editWarehouseReq.MaxNumClusters, "max-num-clusters", 0, `Maximum number of clusters that the autoscaler will create to handle concurrent queries.`)
	editWarehouseCmd.Flags().IntVar(&editWarehouseReq.MinNumClusters, "min-num-clusters", 0, `Minimum number of available clusters that will be maintained for this SQL Endpoint.`)
	editWarehouseCmd.Flags().StringVar(&editWarehouseReq.Name, "name", "", `Logical name for the cluster.`)
	editWarehouseCmd.Flags().Var(&editWarehouseReq.SpotInstancePolicy, "spot-instance-policy", `Configurations whether the endpoint should use spot instances.`)
	// TODO: complex arg: tags
	editWarehouseCmd.Flags().Var(&editWarehouseReq.WarehouseType, "warehouse-type", `Warehouse type (Classic/Pro).`)

}

var editWarehouseCmd = &cobra.Command{
	Use:   "edit-warehouse",
	Short: `Update a warehouse.`,
	Long: `Update a warehouse.
  
  Updates the configuration for a SQL warehouse.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.Warehouses.EditWarehouse(ctx, editWarehouseReq)
		if err != nil {
			return err
		}

		return nil
	},
}

var getWarehouseReq warehouses.GetWarehouse

func init() {
	Cmd.AddCommand(getWarehouseCmd)
	// TODO: short flags

	getWarehouseCmd.Flags().StringVar(&getWarehouseReq.Id, "id", "", `Required.`)

}

var getWarehouseCmd = &cobra.Command{
	Use:   "get-warehouse",
	Short: `Get warehouse info.`,
	Long: `Get warehouse info.
  
  Gets the information for a single SQL warehouse.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Warehouses.GetWarehouse(ctx, getWarehouseReq)
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
	Cmd.AddCommand(getWorkspaceWarehouseConfigCmd)

}

var getWorkspaceWarehouseConfigCmd = &cobra.Command{
	Use:   "get-workspace-warehouse-config",
	Short: `Get the workspace configuration.`,
	Long: `Get the workspace configuration.
  
  Gets the workspace level configuration that is shared by all SQL warehouses in
  a workspace.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Warehouses.GetWorkspaceWarehouseConfig(ctx)
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

var listWarehousesReq warehouses.ListWarehouses

func init() {
	Cmd.AddCommand(listWarehousesCmd)
	// TODO: short flags

	listWarehousesCmd.Flags().IntVar(&listWarehousesReq.RunAsUserId, "run-as-user-id", 0, `Service Principal which will be used to fetch the list of endpoints.`)

}

var listWarehousesCmd = &cobra.Command{
	Use:   "list-warehouses",
	Short: `List warehouses.`,
	Long: `List warehouses.
  
  Lists all SQL warehouses that a user has manager permissions on.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.Warehouses.ListWarehousesAll(ctx, listWarehousesReq)
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

var setWorkspaceWarehouseConfigReq warehouses.SetWorkspaceWarehouseConfigRequest

func init() {
	Cmd.AddCommand(setWorkspaceWarehouseConfigCmd)
	// TODO: short flags

	// TODO: complex arg: channel
	// TODO: complex arg: config_param
	// TODO: array: data_access_config
	setWorkspaceWarehouseConfigCmd.Flags().BoolVar(&setWorkspaceWarehouseConfigReq.EnableDatabricksCompute, "enable-databricks-compute", false, `Enable Serverless compute for SQL Endpoints Deprecated: Use enable_serverless_compute TODO(SC-79930): Remove the field once clients are updated.`)
	setWorkspaceWarehouseConfigCmd.Flags().BoolVar(&setWorkspaceWarehouseConfigReq.EnableServerlessCompute, "enable-serverless-compute", false, `Enable Serverless compute for SQL Endpoints.`)
	// TODO: array: enabled_warehouse_types
	// TODO: complex arg: global_param
	setWorkspaceWarehouseConfigCmd.Flags().StringVar(&setWorkspaceWarehouseConfigReq.GoogleServiceAccount, "google-service-account", "", `GCP only: Google Service Account used to pass to cluster to access Google Cloud Storage.`)
	setWorkspaceWarehouseConfigCmd.Flags().StringVar(&setWorkspaceWarehouseConfigReq.InstanceProfileArn, "instance-profile-arn", "", `AWS Only: Instance profile used to pass IAM role to the cluster.`)
	setWorkspaceWarehouseConfigCmd.Flags().Var(&setWorkspaceWarehouseConfigReq.SecurityPolicy, "security-policy", `Security policy for endpoints.`)
	setWorkspaceWarehouseConfigCmd.Flags().BoolVar(&setWorkspaceWarehouseConfigReq.ServerlessAgreement, "serverless-agreement", false, `Internal.`)
	// TODO: complex arg: sql_configuration_parameters

}

var setWorkspaceWarehouseConfigCmd = &cobra.Command{
	Use:   "set-workspace-warehouse-config",
	Short: `Set the workspace configuration.`,
	Long: `Set the workspace configuration.
  
  Sets the workspace level configuration that is shared by all SQL warehouses in
  a workspace.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.Warehouses.SetWorkspaceWarehouseConfig(ctx, setWorkspaceWarehouseConfigReq)
		if err != nil {
			return err
		}

		return nil
	},
}

var startWarehouseReq warehouses.StartWarehouse

func init() {
	Cmd.AddCommand(startWarehouseCmd)
	// TODO: short flags

	startWarehouseCmd.Flags().StringVar(&startWarehouseReq.Id, "id", "", `Required.`)

}

var startWarehouseCmd = &cobra.Command{
	Use:   "start-warehouse",
	Short: `Start a warehouse.`,
	Long: `Start a warehouse.
  
  Starts a SQL warehouse.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.Warehouses.StartWarehouse(ctx, startWarehouseReq)
		if err != nil {
			return err
		}

		return nil
	},
}

var stopWarehouseReq warehouses.StopWarehouse

func init() {
	Cmd.AddCommand(stopWarehouseCmd)
	// TODO: short flags

	stopWarehouseCmd.Flags().StringVar(&stopWarehouseReq.Id, "id", "", `Required.`)

}

var stopWarehouseCmd = &cobra.Command{
	Use:   "stop-warehouse",
	Short: `Stop a warehouse.`,
	Long: `Stop a warehouse.
  
  Stops a SQL warehouse.`,

	PreRunE: sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		err := w.Warehouses.StopWarehouse(ctx, stopWarehouseReq)
		if err != nil {
			return err
		}

		return nil
	},
}

// end service Warehouses

func init() {
	Cmd.PersistentFlags().String("profile", "", "~/.databrickscfg profile")

	Cmd.AddCommand(query_history.Cmd)
}
