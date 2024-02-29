package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/databricks/cli/cmd/root"
	"github.com/databricks/cli/libs/cmdio"
	"github.com/databricks/cli/libs/env"
	"github.com/databricks/cli/libs/flags"
	"github.com/databricks/databricks-sdk-go/config"
	"github.com/spf13/cobra"
)

var authTemplate = fmt.Sprintf(`{{"Workspace:" | bold}} {{.Details.Host}}.
{{if .Details.Username }}{{"User:" | bold}} {{.Details.Username}}.{{end}}
{{"Authenticated with:" | bold}} {{.Details.AuthType}}
{{if .Details.AccountID }}{{"Account ID:" | bold}} {{.Details.AccountID}}.{{end}}
-----
%s
`, configurationTemplate)

var errorTemplate = fmt.Sprintf(`Unable to authenticate: {{.Error}}
-----
%s
`, configurationTemplate)

const configurationTemplate = `Configuration:
  {{- $details := .Details}}
  {{- range $k, $v := $details.Configuration }}
  {{if $v.AuthTypeMismatch}}x{{else}}✓{{end}} {{$k | bold}}: {{$v.Value}}
  {{- if $v.Source }}
  {{- " (from" | italic}} {{$v.Source.String | italic}}
  {{- if $v.AuthTypeMismatch}}, {{ "not used for auth type " | red | italic }}{{$details.AuthType | red | italic}}{{end}})
  {{- end}}
  {{- end}}
`

func newDescribeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describes which auth credentials is being used by the CLI and identify their source",
	}

	var showSensitive bool
	var isAccount bool
	cmd.Flags().BoolVar(&showSensitive, "sensitive", false, "Include sensitive fields like passwords and tokens in the output")
	cmd.Flags().BoolVar(&isAccount, "account", false, "Describe the account client instead of the workspace client")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var status *authStatus
		var err error
		if isAccount {
			status, err = getAccountAuthStatus(cmd, args, showSensitive, func(cmd *cobra.Command, args []string) (*config.Config, error) {
				err := root.MustAccountClient(cmd, args)
				return root.ConfigUsed(cmd.Context()), err
			})
		} else {
			status, err = getWorkspaceAuthStatus(cmd, args, showSensitive, func(cmd *cobra.Command, args []string) (*config.Config, error) {
				err := root.MustWorkspaceClient(cmd, args)
				return root.ConfigUsed(cmd.Context()), err
			})
		}

		if err != nil {
			return err
		}

		if status.Error != nil {
			return render(ctx, cmd, status, errorTemplate)
		}

		return render(ctx, cmd, status, authTemplate)
	}

	return cmd
}

type tryAuth func(cmd *cobra.Command, args []string) (*config.Config, error)

func getWorkspaceAuthStatus(cmd *cobra.Command, args []string, showSensitive bool, fn tryAuth) (*authStatus, error) {
	cfg, err := fn(cmd, args)
	ctx := cmd.Context()
	if err != nil {
		return &authStatus{
			Status:  "error",
			Error:   err,
			Details: getAuthDetails(ctx, cmd, cfg, showSensitive),
		}, nil
	}

	w := root.WorkspaceClient(ctx)
	me, err := w.CurrentUser.Me(ctx)
	if err != nil {
		return nil, err
	}

	status := authStatus{
		Status:  "success",
		Details: getAuthDetails(ctx, cmd, w.Config, showSensitive),
	}
	status.Details.Username = me.UserName

	return &status, nil
}

func getAccountAuthStatus(cmd *cobra.Command, args []string, showSensitive bool, fn tryAuth) (*authStatus, error) {
	cfg, err := fn(cmd, args)
	ctx := cmd.Context()
	if err != nil {
		return &authStatus{
			Status:  "error",
			Error:   err,
			Details: getAuthDetails(ctx, cmd, cfg, showSensitive),
		}, nil
	}

	a := root.AccountClient(ctx)
	status := authStatus{
		Status:  "success",
		Details: getAuthDetails(ctx, cmd, a.Config, showSensitive),
	}

	status.Details.AccountID = a.Config.AccountID
	status.Details.Username = a.Config.Username
	return &status, nil
}

func render(ctx context.Context, cmd *cobra.Command, status *authStatus, template string) error {
	switch root.OutputType(cmd) {
	case flags.OutputText:
		return cmdio.RenderWithTemplate(ctx, status, "", template)
	case flags.OutputJSON:
		buf, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return err
		}
		cmd.OutOrStdout().Write(buf)
	default:
		return fmt.Errorf("unknown output type %s", root.OutputType(cmd))
	}

	return nil
}

type authStatus struct {
	Status  string      `json:"status"`
	Error   error       `json:"error,omitempty"`
	Details authDetails `json:"details"`
}

type authDetails struct {
	AuthType      string                `json:"auth_type"`
	Username      string                `json:"username,omitempty"`
	Host          string                `json:"host,omitempty"`
	AccountID     string                `json:"account_id,omitempty"`
	Configuration map[string]attrConfig `json:"configuration"`
}

type attrConfig struct {
	Value            string         `json:"value"`
	Source           *config.Source `json:"source"`
	AuthTypeMismatch bool           `json:"auth_type_mismatch"`
}

func getAuthDetails(ctx context.Context, cmd *cobra.Command, cfg *config.Config, showSensitive bool) authDetails {
	attrSet := make(map[string]attrConfig, 0)
	for _, a := range config.ConfigAttributes {
		if a.IsZero(cfg) {
			continue
		}

		attrSet[a.Name] = attrConfig{
			Value:            getValue(cfg, a, showSensitive),
			Source:           getSource(ctx, cmd, cfg, &a),
			AuthTypeMismatch: a.IsAuthAttribute() && !slices.Contains(a.Auth, cfg.AuthType),
		}
	}

	return authDetails{
		AuthType:      cfg.AuthType,
		Host:          cfg.Host,
		Configuration: attrSet,
	}
}

func getValue(cfg *config.Config, a config.ConfigAttribute, showSensitive bool) string {
	if !showSensitive && a.Sensitive {
		return "********"
	}

	return a.GetString(cfg)
}

func getSource(ctx context.Context, cmd *cobra.Command, cfg *config.Config, a *config.ConfigAttribute) *config.Source {
	v := a.GetString(cfg)

	// Check if the value is from an environment variable
	for _, envVar := range a.EnvVars {
		if v == env.Get(ctx, envVar) {
			return &config.Source{
				Type: config.SourceEnv,
				Name: envVar,
			}
		}
	}

	// Check if the value is from command flags
	// Only profile and host can be set like this
	// The rest of the values are set from the config file
	if a.Name == "profile" && cmd.Flag("profile").Changed {
		return &config.Source{
			Type: config.SourceType("flag"),
			Name: "--profile",
		}
	}

	if a.Name == "host" && cmd.Flag("host").Changed {
		return &config.Source{
			Type: config.SourceType("flag"),
			Name: "--host",
		}
	}

	return a.Source
}
