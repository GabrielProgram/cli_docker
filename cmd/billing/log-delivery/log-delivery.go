package log_delivery

import (
	"github.com/databricks/bricks/lib/ui"
	"github.com/databricks/bricks/project"
	"github.com/databricks/databricks-sdk-go/service/billing"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "log-delivery",
	Short: `These APIs manage log delivery configurations for this account.`, // TODO: fix FirstSentence logic and append dot to summary
}

var createReq billing.WrappedCreateLogDeliveryConfiguration

func init() {
	Cmd.AddCommand(createCmd)
	// TODO: short flags

	// TODO: complex arg: log_delivery_configuration

}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: `Create a new log delivery configuration Creates a new Databricks log delivery configuration to enable delivery of the specified type of logs to your storage location.`, // TODO: fix logic

	PreRunE: project.Configure, // TODO: improve logic for bundle/non-bundle invocations
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		a := project.Get(ctx).AccountClient()
		response, err := a.LogDelivery.Create(ctx, createReq)
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

var getReq billing.GetLogDeliveryRequest

func init() {
	Cmd.AddCommand(getCmd)
	// TODO: short flags

	getCmd.Flags().StringVar(&getReq.LogDeliveryConfigurationId, "log-delivery-configuration-id", "", `Databricks log delivery configuration ID.`)

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: `Get log delivery configuration Gets a Databricks log delivery configuration object for an account, both specified by ID.`, // TODO: fix logic

	PreRunE: project.Configure, // TODO: improve logic for bundle/non-bundle invocations
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		a := project.Get(ctx).AccountClient()
		response, err := a.LogDelivery.Get(ctx, getReq)
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

var listReq billing.ListLogDeliveryRequest

func init() {
	Cmd.AddCommand(listCmd)
	// TODO: short flags

	listCmd.Flags().StringVar(&listReq.CredentialsId, "credentials-id", "", `Filter by credential configuration ID.`)
	// TODO: complex arg: status
	listCmd.Flags().StringVar(&listReq.StorageConfigurationId, "storage-configuration-id", "", `Filter by storage configuration ID.`)

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: `Get all log delivery configurations Gets all Databricks log delivery configurations associated with an account specified by ID.`, // TODO: fix logic

	PreRunE: project.Configure, // TODO: improve logic for bundle/non-bundle invocations
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		a := project.Get(ctx).AccountClient()
		response, err := a.LogDelivery.ListAll(ctx, listReq)
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

var patchStatusReq billing.UpdateLogDeliveryConfigurationStatusRequest

func init() {
	Cmd.AddCommand(patchStatusCmd)
	// TODO: short flags

	patchStatusCmd.Flags().StringVar(&patchStatusReq.LogDeliveryConfigurationId, "log-delivery-configuration-id", "", `Databricks log delivery configuration ID.`)
	// TODO: complex arg: status

}

var patchStatusCmd = &cobra.Command{
	Use:   "patch-status",
	Short: `Enable or disable log delivery configuration Enables or disables a log delivery configuration.`, // TODO: fix logic

	PreRunE: project.Configure, // TODO: improve logic for bundle/non-bundle invocations
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		a := project.Get(ctx).AccountClient()
		err := a.LogDelivery.PatchStatus(ctx, patchStatusReq)
		if err != nil {
			return err
		}

		return nil
	},
}
