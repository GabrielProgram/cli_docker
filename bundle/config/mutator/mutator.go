package mutator

import (
	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/config"
	"github.com/databricks/cli/bundle/config/loader"
	"github.com/databricks/cli/bundle/scripts"
)

func DefaultMutators() []bundle.Mutator {
	return []bundle.Mutator{
		// Execute preinit script before loading any configuration files.
		// It needs to be done before processing configuration files to allow
		// the script to modify the configuration or add own configuration files.
		scripts.Execute(config.ScriptPreInit),

		loader.EntryPoint(),
		loader.ProcessRootIncludes(),

		// Verify that the CLI version is within the specified range.
		VerifyCliVersion(),

		EnvironmentsToTargets(),
		InitializeVariables(),
		DefineDefaultTarget(),
		LoadGitDetails(),
	}
}

func DefaultMutatorsForTarget(target string) []bundle.Mutator {
	return append(
		DefaultMutators(),
		SelectTarget(target),
	)
}
