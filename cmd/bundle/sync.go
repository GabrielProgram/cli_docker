package bundle

import (
	"fmt"
	"time"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/deploy/files"
	"github.com/databricks/cli/bundle/phases"
	"github.com/databricks/cli/cmd/bundle/utils"
	"github.com/databricks/cli/cmd/root"
	"github.com/databricks/cli/libs/log"
	"github.com/databricks/cli/libs/sync"
	"github.com/spf13/cobra"
)

type syncFlags struct {
	interval time.Duration
	full     bool
	watch    bool
}

func (f *syncFlags) syncOptionsFromBundle(cmd *cobra.Command, b *bundle.Bundle) (*sync.SyncOptions, error) {
	opts, err := files.GetSyncOptions(cmd.Context(), b)
	if err != nil {
		return nil, fmt.Errorf("cannot get sync options: %w", err)
	}

	opts.Full = f.full
	opts.PollInterval = f.interval
	return opts, nil
}

func newSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync [flags]",
		Short: "Synchronize bundle tree to the workspace",
		Args:  root.NoArgs,

		PreRunE: utils.ConfigureBundleWithVariables,
	}

	var f syncFlags
	cmd.Flags().DurationVar(&f.interval, "interval", 1*time.Second, "file system polling interval (for --watch)")
	cmd.Flags().BoolVar(&f.full, "full", false, "perform full synchronization (default is incremental)")
	cmd.Flags().BoolVar(&f.watch, "watch", false, "watch local file system for changes")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		b := bundle.Get(cmd.Context())

		// Run initialize phase to make sure paths are set.
		err := bundle.Apply(cmd.Context(), b, phases.Initialize())
		if err != nil {
			return err
		}

		opts, err := f.syncOptionsFromBundle(cmd, b)
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		s, err := sync.New(ctx, *opts)
		if err != nil {
			return err
		}

		log.Infof(ctx, "Remote file sync location: %v", opts.RemotePath)

		if f.watch {
			return s.RunContinuous(ctx)
		}

		return s.RunOnce(ctx)
	}

	return cmd
}
