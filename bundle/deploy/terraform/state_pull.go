package terraform

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/deploy"
	"github.com/databricks/cli/libs/diag"
	"github.com/databricks/cli/libs/log"
)

type tfState struct {
	Serial  int    `json:"serial"`
	Lineage string `json:"lineage"`
}

type statePull struct {
	filerFactory deploy.FilerFactory
}

func (l *statePull) Name() string {
	return "terraform:state-pull"
}

func (l *statePull) remoteState(ctx context.Context, b *bundle.Bundle) (*tfState, []byte, error) {
	f, err := l.filerFactory(b)
	if err != nil {
		return nil, nil, err
	}

	remote, err := f.Read(ctx, TerraformStateFileName)
	if err != nil {
		return nil, nil, err
	}
	defer remote.Close()

	remoteContent, err := io.ReadAll(remote)
	if err != nil {
		return nil, nil, err
	}

	remoteState := &tfState{}
	err = json.Unmarshal(remoteContent, remoteState)
	if err != nil {
		return nil, nil, err
	}

	return remoteState, remoteContent, nil
}

func (l *statePull) localState(ctx context.Context, b *bundle.Bundle) (*tfState, error) {
	dir, err := Dir(ctx, b)
	if err != nil {
		return nil, err
	}

	localStatePath := filepath.Join(dir, TerraformStateFileName)
	local, err := os.Open(localStatePath)
	if err != nil {
		return nil, err
	}
	defer local.Close()

	localContent, err := io.ReadAll(local)
	if err != nil {
		return nil, err
	}

	localState := &tfState{}
	err = json.Unmarshal(localContent, localState)
	if err != nil {
		return nil, err
	}

	return localState, nil
}

func (l *statePull) Apply(ctx context.Context, b *bundle.Bundle) diag.Diagnostics {
	dir, err := Dir(ctx, b)
	if err != nil {
		return diag.FromErr(err)
	}

	localStatePath := filepath.Join(dir, TerraformStateFileName)

	// Case: remote state file does not exist. In this case we should not use the
	// local state file because we cannot guarantee it corresponds to the same deployment,
	// that is the same root_path and workspace host.
	remoteState, remoteContent, err := l.remoteState(ctx, b)
	if errors.Is(err, fs.ErrNotExist) {
		log.Infof(ctx, "Remote state file does not exist. Invalidating local terraform state.")
		os.Remove(localStatePath)
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to read remote state file: %v", err)
	}

	// Case: Local host does not exist. In this case we should rely on the remote state file.
	localState, err := l.localState(ctx, b)
	if errors.Is(err, fs.ErrNotExist) {
		log.Infof(ctx, "Local state file does not exist. Using remote terraform state. Invalidating local terraform state.")
		err := os.WriteFile(localStatePath, remoteContent, 0600)
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.Errorf("failed to read local state file: %v", err)
	}

	// If the lineages do not match, the terraform state files do not correspond to the same deployment.
	if localState.Lineage != remoteState.Lineage {
		log.Infof(ctx, "Remote and local state lineages do not match. Using remote terraform state. Invalidating local terraform state.")
		err := os.WriteFile(localStatePath, remoteContent, 0600)
		return diag.FromErr(err)
	}

	// If the remote state is newer than the local state, we should use the remote state.
	if remoteState.Serial > localState.Serial {
		log.Infof(ctx, "Remote state is newer than local state. Using remote terraform state.")
		err := os.WriteFile(localStatePath, remoteContent, 0600)
		return diag.FromErr(err)
	}

	return nil
}

func StatePull() bundle.Mutator {
	return &statePull{deploy.StateFiler}
}
