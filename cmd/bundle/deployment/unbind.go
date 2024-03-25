package deployment

import (
	"context"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/phases"
	"github.com/databricks/cli/cmd/bundle/utils"
	"github.com/databricks/cli/cmd/root"
	"github.com/spf13/cobra"
)

func newUnbindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unbind KEY",
		Short:   "Unbind bundle-defined resources from its managed remote resource",
		Args:    root.ExactArgs(1),
		PreRunE: utils.ConfigureBundleWithVariables,
	}

	var forceLock bool
	cmd.Flags().BoolVar(&forceLock, "force-lock", false, "Force acquisition of deployment lock.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		b := bundle.Get(ctx)
		resource, err := b.Config.Resources.FindResourceByConfigKey(args[0])
		if err != nil {
			return err
		}

		bundle.ApplyFunc(ctx, b, func(context.Context, *bundle.Bundle) error {
			b.Config.Bundle.Deployment.Lock.Force = forceLock
			return nil
		})

		return bundle.Apply(cmd.Context(), b, bundle.Seq(
			phases.Initialize(),
			phases.Unbind(resource.TerraformResourceName(), args[0]),
		))
	}

	return cmd
}
