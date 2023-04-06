package files

import (
	"context"
	"fmt"

	"github.com/databricks/bricks/bundle"
	"github.com/databricks/bricks/libs/cmdio"
	"github.com/databricks/databricks-sdk-go/service/workspace"
	"github.com/fatih/color"
)

type delete struct{}

func (m *delete) Name() string {
	return "files.Delete"
}

func (m *delete) Apply(ctx context.Context, b *bundle.Bundle) ([]bundle.Mutator, error) {
	// Do not delete files if terraform destroy was not consented
	if !b.Plan.IsEmpty && !b.Plan.ConfirmApply {
		return nil, nil
	}

	cmdio.LogMutatorEvent(ctx, m.Name(), cmdio.MutatorRunning, "Starting deletion of remote bundle files")
	cmdio.LogMutatorEvent(ctx, m.Name(), cmdio.MutatorRunning, fmt.Sprintf("Bundle remote deployment location: %s", b.Config.Workspace.Root))
	// logger, ok := cmdio.FromContext(ctx)
	// if !ok {
	// 	return nil, fmt.Errorf("no logger found at destroy mutator")
	// }
	// bytes, _ := json.MarshalIndent(logger, "", "  ")
	// fmt.Println("\n\n\n\n\n" + string(bytes) + "\n\n\n\n\n\n")

	red := color.New(color.FgRed).SprintFunc()
	if !b.AutoApprove {
		proceed, err := cmdio.Ask(ctx, fmt.Sprintf("\n%s and all files in it will be %s Proceed?: ", b.Config.Workspace.Root, red("deleted permanently!")))
		if err != nil {
			return nil, err
		}
		if !proceed {
			return nil, nil
		}
	}

	err := b.WorkspaceClient().Workspace.Delete(ctx, workspace.Delete{
		Path:      b.Config.Workspace.Root,
		Recursive: true,
	})
	if err != nil {
		return nil, err
	}

	cmdio.LogMutatorEvent(ctx, m.Name(), cmdio.MutatorCompleted, "Successfully deleted files!")
	return nil, nil
}

func Delete() bundle.Mutator {
	return &delete{}
}
