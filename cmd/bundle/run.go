package bundle

import (
	"encoding/json"
	"fmt"

	"github.com/databricks/bricks/bundle"
	"github.com/databricks/bricks/bundle/deploy/terraform"
	"github.com/databricks/bricks/bundle/phases"
	"github.com/databricks/bricks/bundle/run"
	"github.com/databricks/bricks/cmd/root"
	"github.com/databricks/bricks/libs/flags"
	"github.com/spf13/cobra"
)

var runOptions run.Options

var runCmd = &cobra.Command{
	Use:   "run [flags] KEY",
	Short: "Run a workload (e.g. a job or a pipeline)",

	Args:    cobra.ExactArgs(1),
	PreRunE: root.MustConfigureBundle,
	RunE: func(cmd *cobra.Command, args []string) error {
		b := bundle.Get(cmd.Context())
		err := bundle.Apply(cmd.Context(), b, []bundle.Mutator{
			phases.Initialize(),
			terraform.Load(),
		})
		if err != nil {
			return err
		}

		runner, err := run.Find(b, args[0])
		if err != nil {
			return err
		}

		// TODO: address in PR that I chose not to add output-task flag
		// Rather the right thing to do would be maye
		output, err := runner.Run(cmd.Context(), &runOptions)
		if output != nil {
			switch outputType {
			case flags.OutputText:
				resultString, err := output.String()
				if err != nil {
					return err
				}
				fmt.Print(resultString)
			case flags.OutputJSON:
				b, err := json.MarshalIndent(output, "", "  ")
				if err != nil {
					return err
				}
				fmt.Print(string(b))
			default:
				return fmt.Errorf("unknown output type %s", outputType)
			}
		}
		if err != nil {
			return err
		}
		return nil
	},

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		err := root.MustConfigureBundle(cmd, args)
		if err != nil {
			cobra.CompErrorln(err.Error())
			return nil, cobra.ShellCompDirectiveError
		}

		// No completion in the context of a bundle.
		// Source and destination paths are taken from bundle configuration.
		b := bundle.GetOrNil(cmd.Context())
		if b == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return run.ResourceCompletions(b), cobra.ShellCompDirectiveNoFileComp
	},
}

var outputType flags.Output = flags.OutputText

func init() {
	runOptions.Define(runCmd.Flags())
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().Var(&outputType, "output", "type of output format")
}
