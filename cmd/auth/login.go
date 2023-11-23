package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/databricks/cli/libs/auth"
	"github.com/databricks/cli/libs/cmdio"
	"github.com/databricks/cli/libs/databrickscfg"
	"github.com/databricks/cli/libs/databrickscfg/cfgpickers"
	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/databricks-sdk-go/config"
	"github.com/spf13/cobra"
)

func configureHost(ctx context.Context, persistentAuth *auth.PersistentAuth, args []string, argIndex int) error {
	if len(args) > argIndex {
		persistentAuth.Host = args[argIndex]
		return nil
	}

	host, err := promptForHost(ctx)
	if err != nil {
		return err
	}
	persistentAuth.Host = host
	return nil
}

const minimalDbConnectVersion = "13.1"

func newLoginCommand(persistentAuth *auth.PersistentAuth) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [HOST]",
		Short: "Authenticate this machine",
	}

	var loginTimeout time.Duration
	var configureCluster bool
	cmd.Flags().DurationVar(&loginTimeout, "timeout", auth.DefaultTimeout,
		"Timeout for completing login challenge in the browser")
	cmd.Flags().BoolVar(&configureCluster, "configure-cluster", false,
		"Prompts to configure cluster")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var profileName string
		profileFlag := cmd.Flag("profile")
		if profileFlag != nil && profileFlag.Value.String() != "" {
			profileName = profileFlag.Value.String()
		} else if cmdio.IsInTTY(ctx) {
			prompt := cmdio.Prompt(ctx)
			prompt.Label = "Databricks Profile Name"
			prompt.Default = persistentAuth.ProfileName()
			prompt.AllowEdit = true
			profile, err := prompt.Run()
			if err != nil {
				return err
			}
			profileName = profile
		}

		err := setHost(ctx, profileName, persistentAuth, args)
		if err != nil {
			return err
		}
		defer persistentAuth.Close()

		// We need the config without the profile before it's used to initialise new workspace client below.
		// Otherwise it will complain about non existing profile because it was not yet saved.
		cfg := config.Config{
			Host:     persistentAuth.Host,
			AuthType: "databricks-cli",
		}
		if cfg.IsAccountClient() && persistentAuth.AccountID == "" {
			accountId, err := promptForAccountID(ctx)
			if err != nil {
				return err
			}
			persistentAuth.AccountID = accountId
		}
		cfg.AccountID = persistentAuth.AccountID

		ctx, cancel := context.WithTimeout(ctx, loginTimeout)
		defer cancel()

		err = persistentAuth.Challenge(ctx)
		if err != nil {
			return err
		}

		if configureCluster {
			w, err := databricks.NewWorkspaceClient((*databricks.Config)(&cfg))
			if err != nil {
				return err
			}
			ctx := cmd.Context()
			clusterID, err := cfgpickers.AskForCluster(ctx, w,
				cfgpickers.WithDatabricksConnect(minimalDbConnectVersion))
			if err != nil {
				return err
			}
			cfg.ClusterID = clusterID
		}

		if profileName != "" {
			err = databrickscfg.SaveToProfile(ctx, &config.Config{
				Profile:   profileName,
				Host:      cfg.Host,
				AuthType:  cfg.AuthType,
				AccountID: cfg.AccountID,
				ClusterID: cfg.ClusterID,
			})
			if err != nil {
				return err
			}

			cmdio.LogString(ctx, fmt.Sprintf("Profile %s was successfully saved", profileName))
		}

		return nil
	}

	return cmd
}

func setHost(ctx context.Context, profileName string, persistentAuth *auth.PersistentAuth, args []string) error {
	// If the chosen profile has a hostname and the user hasn't specified a host, infer the host from the profile.
	_, profiles, err := databrickscfg.LoadProfiles(ctx, func(p databrickscfg.Profile) bool {
		return p.Name == profileName
	})
	// Tolerate ErrNoConfiguration here, as we will write out a configuration as part of the login flow.
	if !errors.Is(err, databrickscfg.ErrNoConfiguration) {
		return err
	}
	if persistentAuth.Host == "" {
		if len(profiles) > 0 && profiles[0].Host != "" {
			persistentAuth.Host = profiles[0].Host
		} else {
			configureHost(ctx, persistentAuth, args, 0)
		}
	}
	return nil
}
