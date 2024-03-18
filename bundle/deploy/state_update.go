package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/deploy/files"
	"github.com/databricks/cli/internal/build"
	"github.com/databricks/cli/libs/log"
)

type stateUpdate struct {
}

func (s *stateUpdate) Name() string {
	return "deploy:state-update"
}

func (s *stateUpdate) Apply(ctx context.Context, b *bundle.Bundle) error {
	state, err := load(ctx, b)
	if err != nil {
		return err
	}

	// Increment the state sequence.
	state.Seq = state.Seq + 1

	// Update timestamp.
	state.Timestamp = time.Now().UTC()

	// Update the CLI version and deployment state version.
	state.CliVersion = build.GetInfo().Version
	state.Version = DeploymentStateVersion

	// Get the current file list.
	sync, err := files.GetSync(ctx, b)
	if err != nil {
		return err
	}

	files, err := sync.GetFileList(ctx)
	if err != nil {
		return err
	}

	// Update the state with the current file list.
	fl, err := FromSlice(files)
	if err != nil {
		return err
	}
	state.Files = fl

	statePath, err := getPathToStateFile(ctx, b)
	if err != nil {
		return err
	}
	// Write the state back to the file.
	f, err := os.OpenFile(statePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		log.Infof(ctx, "Unable to open deployment state file: %s", err)
		return err
	}
	defer f.Close()

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, bytes.NewReader(data))
	if err != nil {
		return err
	}

	return nil
}

func StateUpdate() bundle.Mutator {
	return &stateUpdate{}
}

func load(ctx context.Context, b *bundle.Bundle) (*DeploymentState, error) {
	// If the file does not exist, return a new DeploymentState.
	statePath, err := getPathToStateFile(ctx, b)
	if err != nil {
		return nil, err
	}

	log.Infof(ctx, "Loading deployment state from %s", statePath)
	f, err := os.Open(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Infof(ctx, "No deployment state file found")
			return &DeploymentState{
				Version:    DeploymentStateVersion,
				CliVersion: build.GetInfo().Version,
			}, nil
		}
		return nil, err
	}
	defer f.Close()
	return loadState(f)
}
