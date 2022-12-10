// Code generated from OpenAPI specs by Databricks SDK Generator. DO NOT EDIT.

package m_lflow_metrics

import (
	"github.com/databricks/bricks/lib/sdk"
	"github.com/databricks/bricks/lib/ui"
	"github.com/databricks/databricks-sdk-go/service/mlflow"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "m-lflow-metrics",
}

// start get-history command

var getHistoryReq mlflow.GetHistoryRequest

func init() {
	Cmd.AddCommand(getHistoryCmd)
	// TODO: short flags

	getHistoryCmd.Flags().StringVar(&getHistoryReq.MetricKey, "metric-key", getHistoryReq.MetricKey, `Name of the metric.`)
	getHistoryCmd.Flags().StringVar(&getHistoryReq.RunId, "run-id", getHistoryReq.RunId, `ID of the run from which to fetch metric values.`)
	getHistoryCmd.Flags().StringVar(&getHistoryReq.RunUuid, "run-uuid", getHistoryReq.RunUuid, `[Deprecated, use run_id instead] ID of the run from which to fetch metric values.`)

}

var getHistoryCmd = &cobra.Command{
	Use:   "get-history",
	Short: `Get all history.`,
	Long: `Get all history.
  
  Gets a list of all values for the specified metric for a given run.`,

	Annotations: map[string]string{},
	PreRunE:     sdk.PreWorkspaceClient,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := sdk.WorkspaceClient(ctx)
		response, err := w.MLflowMetrics.GetHistory(ctx, getHistoryReq)
		if err != nil {
			return err
		}
		return ui.Render(cmd, response)
	},
}

// end service MLflowMetrics
