package bundle

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/deploy/terraform"
	"github.com/databricks/cli/bundle/phases"
	"github.com/databricks/cli/cmd/bundle/utils"
	"github.com/databricks/cli/cmd/root"
	"github.com/databricks/cli/libs/flags"
	"github.com/spf13/cobra"
)

func newSummaryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "summary",
		Short:   "Describe the bundle resources and their deployment states",
		Args:    root.NoArgs,
		PreRunE: utils.ConfigureBundleWithVariables,

		// This command is currently intended for the Databricks VSCode extension only
		Hidden: true,
	}

	var forcePull bool
	cmd.Flags().BoolVar(&forcePull, "force-pull", false, "Skip local cache and load the state from the remote workspace")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		b := bundle.Get(cmd.Context())

		err := bundle.Apply(cmd.Context(), b, phases.Initialize())
		if err != nil {
			return err
		}

		cacheDir, err := terraform.Dir(cmd.Context(), b)
		if err != nil {
			return err
		}
		_, stateFileErr := os.Stat(filepath.Join(cacheDir, terraform.TerraformStateFileName))
		_, configFileErr := os.Stat(filepath.Join(cacheDir, terraform.TerraformConfigFileName))
		noCache := errors.Is(stateFileErr, os.ErrNotExist) || errors.Is(configFileErr, os.ErrNotExist)

		if forcePull || noCache {
			err = bundle.Apply(cmd.Context(), b, bundle.Seq(
				terraform.StatePull(),
				terraform.Interpolate(),
				terraform.Write(),
			))
			if err != nil {
				return err
			}
		}

		err = bundle.Apply(cmd.Context(), b, terraform.Load())
		if err != nil {
			return err
		}

		switch root.OutputType(cmd) {
		case flags.OutputText:
			return fmt.Errorf("%w, only json output is supported", errors.ErrUnsupported)
		case flags.OutputJSON:
			buf, err := json.MarshalIndent(b.Config, "", "  ")
			if err != nil {
				return err
			}
			cmd.OutOrStdout().Write(buf)
		default:
			return fmt.Errorf("unknown output type %s", root.OutputType(cmd))
		}

		return nil
	}

	return cmd
}
