package artifacts

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/artifacts/whl"
	"github.com/databricks/cli/bundle/config"
	"github.com/databricks/cli/libs/cmdio"
	"github.com/databricks/databricks-sdk-go/service/workspace"
)

type mutatorFactory = func(name string) bundle.Mutator

var buildMutators map[config.ArtifactType]mutatorFactory = map[config.ArtifactType]mutatorFactory{
	config.ArtifactPythonWheel: whl.Build,
}

var uploadMutators map[config.ArtifactType]mutatorFactory = map[config.ArtifactType]mutatorFactory{}

func getBuildMutator(t config.ArtifactType, name string) bundle.Mutator {
	mutatorFactory, ok := buildMutators[t]
	if !ok {
		mutatorFactory = BasicBuild
	}

	return mutatorFactory(name)
}

func getUploadMutator(t config.ArtifactType, name string) bundle.Mutator {
	mutatorFactory, ok := uploadMutators[t]
	if !ok {
		mutatorFactory = BasicUpload
	}

	return mutatorFactory(name)
}

// Basic Build defines a general build mutator which builds artifact based on artifact.BuildCommand
type basicBuild struct {
	name string
}

func BasicBuild(name string) bundle.Mutator {
	return &basicBuild{name: name}
}

func (m *basicBuild) Name() string {
	return fmt.Sprintf("artifacts.Build(%s)", m.name)
}

func (m *basicBuild) Apply(ctx context.Context, b *bundle.Bundle) error {
	artifact, ok := b.Config.Artifacts[m.name]
	if !ok {
		return fmt.Errorf("artifact doesn't exist: %s", m.name)
	}

	cmdio.LogString(ctx, fmt.Sprintf("artifacts.Build(%s): Building...", m.name))

	out, err := artifact.Build(ctx)
	if err != nil {
		return fmt.Errorf("artifacts.Build(%s): %w, output: %s", m.name, err, out)
	}
	cmdio.LogString(ctx, fmt.Sprintf("artifacts.Build(%s): Build succeeded", m.name))

	return nil
}

// Basic Upload defines a general upload mutator which uploads artifact as a library to workspace
type basicUpload struct {
	name string
}

func BasicUpload(name string) bundle.Mutator {
	return &basicUpload{name: name}
}

func (m *basicUpload) Name() string {
	return fmt.Sprintf("artifacts.Build(%s)", m.name)
}

func (m *basicUpload) Apply(ctx context.Context, b *bundle.Bundle) error {
	artifact, ok := b.Config.Artifacts[m.name]
	if !ok {
		return fmt.Errorf("artifact doesn't exist: %s", m.name)
	}

	if artifact.File == "" {
		return fmt.Errorf("artifact source is not configured: %s", m.name)
	}

	cmdio.LogString(ctx, fmt.Sprintf("artifacts.Upload(%s): Uploading...", m.name))

	r, err := uploadArtifact(ctx, artifact, b)
	if err != nil {
		return fmt.Errorf("artifacts.Upload(%s): %w", m.name, err)
	}

	artifact.RemotePath = r
	cmdio.LogString(ctx, fmt.Sprintf("artifacts.Upload(%s): Upload succeeded", m.name))
	return nil
}

// Function to upload artifact as a library to Workspace
// Currenly it does not work correctly because Workspace.Import API can not import libraries
// but only notebooks or files
func uploadArtifact(ctx context.Context, a *config.Artifact, b *bundle.Bundle) (string, error) {
	raw, err := os.ReadFile(a.File)
	if err != nil {
		return "", fmt.Errorf("unable to read %s: %w", a.File, errors.Unwrap(err))
	}

	artifactPath := b.Config.Workspace.ArtifactsPath
	if artifactPath == "" {
		return "", fmt.Errorf("remote artifact path not configured")
	}

	remotePath := path.Join(artifactPath, path.Base(a.File))

	// Make sure target directory exists.
	err = b.WorkspaceClient().Workspace.MkdirsByPath(ctx, path.Dir(remotePath))
	if err != nil {
		return "", fmt.Errorf("unable to create directory for %s: %w", remotePath, err)
	}

	// Import to workspace.
	err = b.WorkspaceClient().Workspace.Import(ctx, workspace.Import{
		Path:      remotePath,
		Overwrite: true,
		Format:    workspace.ImportFormatAuto,
		Content:   base64.StdEncoding.EncodeToString(raw),
	})
	if err != nil {
		return "", fmt.Errorf("unable to import %s: %w", remotePath, err)
	}

	return remotePath, nil
}
