// Code generated from OpenAPI specs by Databricks SDK Generator. DO NOT EDIT.

package secrets

import (
	"fmt"

	"github.com/databricks/cli/cmd/root"
	"github.com/databricks/cli/libs/cmdio"
	"github.com/databricks/cli/libs/flags"
	"github.com/databricks/databricks-sdk-go/service/workspace"
	"github.com/spf13/cobra"
)

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var cmdOverrides []func(*cobra.Command)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: `The Secrets API allows you to manage secrets, secret scopes, and access permissions.`,
		Long: `The Secrets API allows you to manage secrets, secret scopes, and access
  permissions.
  
  Sometimes accessing data requires that you authenticate to external data
  sources through JDBC. Instead of directly entering your credentials into a
  notebook, use Databricks secrets to store your credentials and reference them
  in notebooks and jobs.
  
  Administrators, secret creators, and users granted permission can read
  Databricks secrets. While Databricks makes an effort to redact secret values
  that might be displayed in notebooks, it is not possible to prevent such users
  from reading secrets.`,
		GroupID: "workspace",
		Annotations: map[string]string{
			"package": "workspace",
		},
	}

	// Apply optional overrides to this command.
	for _, fn := range cmdOverrides {
		fn(cmd)
	}

	return cmd
}

// start create-scope command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var createScopeOverrides []func(
	*cobra.Command,
	*workspace.CreateScope,
)

func newCreateScope() *cobra.Command {
	cmd := &cobra.Command{}

	var createScopeReq workspace.CreateScope
	var createScopeJson flags.JsonFlag

	// TODO: short flags
	cmd.Flags().Var(&createScopeJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	// TODO: complex arg: backend_azure_keyvault
	cmd.Flags().StringVar(&createScopeReq.InitialManagePrincipal, "initial-manage-principal", createScopeReq.InitialManagePrincipal, `The principal that is initially granted MANAGE permission to the created scope.`)
	cmd.Flags().Var(&createScopeReq.ScopeBackendType, "scope-backend-type", `The backend type the scope will be created with.`)

	cmd.Use = "create-scope SCOPE"
	cmd.Short = `Create a new secret scope.`
	cmd.Long = `Create a new secret scope.
  
  The scope name must consist of alphanumeric characters, dashes, underscores,
  and periods, and may not exceed 128 characters. The maximum number of scopes
  in a workspace is 100.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(1)
		if cmd.Flags().Changed("json") {
			check = cobra.ExactArgs(0)
		}
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		if cmd.Flags().Changed("json") {
			err = createScopeJson.Unmarshal(&createScopeReq)
			if err != nil {
				return err
			}
		} else {
			createScopeReq.Scope = args[0]
		}

		err = w.Secrets.CreateScope(ctx, createScopeReq)
		if err != nil {
			return err
		}
		return nil
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range createScopeOverrides {
		fn(cmd, &createScopeReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newCreateScope())
	})
}

// start delete-acl command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var deleteAclOverrides []func(
	*cobra.Command,
	*workspace.DeleteAcl,
)

func newDeleteAcl() *cobra.Command {
	cmd := &cobra.Command{}

	var deleteAclReq workspace.DeleteAcl
	var deleteAclJson flags.JsonFlag

	// TODO: short flags
	cmd.Flags().Var(&deleteAclJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	cmd.Use = "delete-acl SCOPE PRINCIPAL"
	cmd.Short = `Delete an ACL.`
	cmd.Long = `Delete an ACL.
  
  Deletes the given ACL on the given scope.
  
  Users must have the MANAGE permission to invoke this API. Throws
  RESOURCE_DOES_NOT_EXIST if no such secret scope, principal, or ACL exists.
  Throws PERMISSION_DENIED if the user does not have permission to make this
  API call.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(2)
		if cmd.Flags().Changed("json") {
			check = cobra.ExactArgs(0)
		}
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		if cmd.Flags().Changed("json") {
			err = deleteAclJson.Unmarshal(&deleteAclReq)
			if err != nil {
				return err
			}
		} else {
			deleteAclReq.Scope = args[0]
			deleteAclReq.Principal = args[1]
		}

		err = w.Secrets.DeleteAcl(ctx, deleteAclReq)
		if err != nil {
			return err
		}
		return nil
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range deleteAclOverrides {
		fn(cmd, &deleteAclReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newDeleteAcl())
	})
}

// start delete-scope command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var deleteScopeOverrides []func(
	*cobra.Command,
	*workspace.DeleteScope,
)

func newDeleteScope() *cobra.Command {
	cmd := &cobra.Command{}

	var deleteScopeReq workspace.DeleteScope
	var deleteScopeJson flags.JsonFlag

	// TODO: short flags
	cmd.Flags().Var(&deleteScopeJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	cmd.Use = "delete-scope SCOPE"
	cmd.Short = `Delete a secret scope.`
	cmd.Long = `Delete a secret scope.
  
  Deletes a secret scope.
  
  Throws RESOURCE_DOES_NOT_EXIST if the scope does not exist. Throws
  PERMISSION_DENIED if the user does not have permission to make this API
  call.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(1)
		if cmd.Flags().Changed("json") {
			check = cobra.ExactArgs(0)
		}
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		if cmd.Flags().Changed("json") {
			err = deleteScopeJson.Unmarshal(&deleteScopeReq)
			if err != nil {
				return err
			}
		} else {
			deleteScopeReq.Scope = args[0]
		}

		err = w.Secrets.DeleteScope(ctx, deleteScopeReq)
		if err != nil {
			return err
		}
		return nil
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range deleteScopeOverrides {
		fn(cmd, &deleteScopeReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newDeleteScope())
	})
}

// start delete-secret command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var deleteSecretOverrides []func(
	*cobra.Command,
	*workspace.DeleteSecret,
)

func newDeleteSecret() *cobra.Command {
	cmd := &cobra.Command{}

	var deleteSecretReq workspace.DeleteSecret
	var deleteSecretJson flags.JsonFlag

	// TODO: short flags
	cmd.Flags().Var(&deleteSecretJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	cmd.Use = "delete-secret SCOPE KEY"
	cmd.Short = `Delete a secret.`
	cmd.Long = `Delete a secret.
  
  Deletes the secret stored in this secret scope. You must have WRITE or
  MANAGE permission on the secret scope.
  
  Throws RESOURCE_DOES_NOT_EXIST if no such secret scope or secret exists.
  Throws PERMISSION_DENIED if the user does not have permission to make this
  API call.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(2)
		if cmd.Flags().Changed("json") {
			check = cobra.ExactArgs(0)
		}
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		if cmd.Flags().Changed("json") {
			err = deleteSecretJson.Unmarshal(&deleteSecretReq)
			if err != nil {
				return err
			}
		} else {
			deleteSecretReq.Scope = args[0]
			deleteSecretReq.Key = args[1]
		}

		err = w.Secrets.DeleteSecret(ctx, deleteSecretReq)
		if err != nil {
			return err
		}
		return nil
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range deleteSecretOverrides {
		fn(cmd, &deleteSecretReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newDeleteSecret())
	})
}

// start get-acl command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var getAclOverrides []func(
	*cobra.Command,
	*workspace.GetAclRequest,
)

func newGetAcl() *cobra.Command {
	cmd := &cobra.Command{}

	var getAclReq workspace.GetAclRequest

	// TODO: short flags

	cmd.Use = "get-acl SCOPE PRINCIPAL"
	cmd.Short = `Get secret ACL details.`
	cmd.Long = `Get secret ACL details.
  
  Gets the details about the given ACL, such as the group and permission. Users
  must have the MANAGE permission to invoke this API.
  
  Throws RESOURCE_DOES_NOT_EXIST if no such secret scope exists. Throws
  PERMISSION_DENIED if the user does not have permission to make this API
  call.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(2)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		getAclReq.Scope = args[0]
		getAclReq.Principal = args[1]

		response, err := w.Secrets.GetAcl(ctx, getAclReq)
		if err != nil {
			return err
		}
		return cmdio.Render(ctx, response)
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range getAclOverrides {
		fn(cmd, &getAclReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newGetAcl())
	})
}

// start get-secret command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var getSecretOverrides []func(
	*cobra.Command,
	*workspace.GetSecretRequest,
)

func newGetSecret() *cobra.Command {
	cmd := &cobra.Command{}

	var getSecretReq workspace.GetSecretRequest

	// TODO: short flags

	cmd.Use = "get-secret SCOPE KEY"
	cmd.Short = `Get a secret.`
	cmd.Long = `Get a secret.
  
  Gets the bytes representation of a secret value for the specified scope and
  key.
  
  Users need the READ permission to make this call.
  
  Note that the secret value returned is in bytes. The interpretation of the
  bytes is determined by the caller in DBUtils and the type the data is decoded
  into.
  
  Throws PERMISSION_DENIED if the user does not have permission to make this
  API call. Throws RESOURCE_DOES_NOT_EXIST if no such secret or secret scope
  exists.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(2)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		getSecretReq.Scope = args[0]
		getSecretReq.Key = args[1]

		response, err := w.Secrets.GetSecret(ctx, getSecretReq)
		if err != nil {
			return err
		}
		return cmdio.Render(ctx, response)
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range getSecretOverrides {
		fn(cmd, &getSecretReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newGetSecret())
	})
}

// start list-acls command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var listAclsOverrides []func(
	*cobra.Command,
	*workspace.ListAclsRequest,
)

func newListAcls() *cobra.Command {
	cmd := &cobra.Command{}

	var listAclsReq workspace.ListAclsRequest

	// TODO: short flags

	cmd.Use = "list-acls SCOPE"
	cmd.Short = `Lists ACLs.`
	cmd.Long = `Lists ACLs.
  
  List the ACLs for a given secret scope. Users must have the MANAGE
  permission to invoke this API.
  
  Throws RESOURCE_DOES_NOT_EXIST if no such secret scope exists. Throws
  PERMISSION_DENIED if the user does not have permission to make this API
  call.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(1)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		listAclsReq.Scope = args[0]

		response, err := w.Secrets.ListAclsAll(ctx, listAclsReq)
		if err != nil {
			return err
		}
		return cmdio.Render(ctx, response)
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range listAclsOverrides {
		fn(cmd, &listAclsReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newListAcls())
	})
}

// start list-scopes command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var listScopesOverrides []func(
	*cobra.Command,
)

func newListScopes() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "list-scopes"
	cmd.Short = `List all scopes.`
	cmd.Long = `List all scopes.
  
  Lists all secret scopes available in the workspace.
  
  Throws PERMISSION_DENIED if the user does not have permission to make this
  API call.`

	cmd.Annotations = make(map[string]string)

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)
		response, err := w.Secrets.ListScopesAll(ctx)
		if err != nil {
			return err
		}
		return cmdio.Render(ctx, response)
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range listScopesOverrides {
		fn(cmd)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newListScopes())
	})
}

// start list-secrets command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var listSecretsOverrides []func(
	*cobra.Command,
	*workspace.ListSecretsRequest,
)

func newListSecrets() *cobra.Command {
	cmd := &cobra.Command{}

	var listSecretsReq workspace.ListSecretsRequest

	// TODO: short flags

	cmd.Use = "list-secrets SCOPE"
	cmd.Short = `List secret keys.`
	cmd.Long = `List secret keys.
  
  Lists the secret keys that are stored at this scope. This is a metadata-only
  operation; secret data cannot be retrieved using this API. Users need the READ
  permission to make this call.
  
  The lastUpdatedTimestamp returned is in milliseconds since epoch. Throws
  RESOURCE_DOES_NOT_EXIST if no such secret scope exists. Throws
  PERMISSION_DENIED if the user does not have permission to make this API
  call.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(1)
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		listSecretsReq.Scope = args[0]

		response, err := w.Secrets.ListSecretsAll(ctx, listSecretsReq)
		if err != nil {
			return err
		}
		return cmdio.Render(ctx, response)
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range listSecretsOverrides {
		fn(cmd, &listSecretsReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newListSecrets())
	})
}

// start put-acl command

// Slice with functions to override default command behavior.
// Functions can be added from the `init()` function in manually curated files in this directory.
var putAclOverrides []func(
	*cobra.Command,
	*workspace.PutAcl,
)

func newPutAcl() *cobra.Command {
	cmd := &cobra.Command{}

	var putAclReq workspace.PutAcl
	var putAclJson flags.JsonFlag

	// TODO: short flags
	cmd.Flags().Var(&putAclJson, "json", `either inline JSON string or @path/to/file.json with request body`)

	cmd.Use = "put-acl SCOPE PRINCIPAL PERMISSION"
	cmd.Short = `Create/update an ACL.`
	cmd.Long = `Create/update an ACL.
  
  Creates or overwrites the Access Control List (ACL) associated with the given
  principal (user or group) on the specified scope point.
  
  In general, a user or group will use the most powerful permission available to
  them, and permissions are ordered as follows:
  
  * MANAGE - Allowed to change ACLs, and read and write to this secret scope.
  * WRITE - Allowed to read and write to this secret scope. * READ - Allowed
  to read this secret scope and list what secrets are available.
  
  Note that in general, secret values can only be read from within a command on
  a cluster (for example, through a notebook). There is no API to read the
  actual secret value material outside of a cluster. However, the user's
  permission will be applied based on who is executing the command, and they
  must have at least READ permission.
  
  Users must have the MANAGE permission to invoke this API.
  
  The principal is a user or group name corresponding to an existing Databricks
  principal to be granted or revoked access.
  
  Throws RESOURCE_DOES_NOT_EXIST if no such secret scope exists. Throws
  RESOURCE_ALREADY_EXISTS if a permission for the principal already exists.
  Throws INVALID_PARAMETER_VALUE if the permission or principal is invalid.
  Throws PERMISSION_DENIED if the user does not have permission to make this
  API call.`

	cmd.Annotations = make(map[string]string)

	cmd.Args = func(cmd *cobra.Command, args []string) error {
		check := cobra.ExactArgs(3)
		if cmd.Flags().Changed("json") {
			check = cobra.ExactArgs(0)
		}
		return check(cmd, args)
	}

	cmd.PreRunE = root.MustWorkspaceClient
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()
		w := root.WorkspaceClient(ctx)

		if cmd.Flags().Changed("json") {
			err = putAclJson.Unmarshal(&putAclReq)
			if err != nil {
				return err
			}
		} else {
			putAclReq.Scope = args[0]
			putAclReq.Principal = args[1]
			_, err = fmt.Sscan(args[2], &putAclReq.Permission)
			if err != nil {
				return fmt.Errorf("invalid PERMISSION: %s", args[2])
			}
		}

		err = w.Secrets.PutAcl(ctx, putAclReq)
		if err != nil {
			return err
		}
		return nil
	}

	// Disable completions since they are not applicable.
	// Can be overridden by manual implementation in `override.go`.
	cmd.ValidArgsFunction = cobra.NoFileCompletions

	// Apply optional overrides to this command.
	for _, fn := range putAclOverrides {
		fn(cmd, &putAclReq)
	}

	return cmd
}

func init() {
	cmdOverrides = append(cmdOverrides, func(cmd *cobra.Command) {
		cmd.AddCommand(newPutAcl())
	})
}

// end service Secrets
