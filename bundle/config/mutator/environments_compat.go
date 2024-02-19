package mutator

import (
	"context"
	"fmt"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/libs/dyn"
)

type environmentsToTargets struct{}

func EnvironmentsToTargets() bundle.Mutator {
	return &environmentsToTargets{}
}

func (m *environmentsToTargets) Name() string {
	return "EnvironmentsToTargets"
}

func (m *environmentsToTargets) Apply(ctx context.Context, b *bundle.Bundle) error {
	// Short circuit if the "environments" key is not set.
	// This is the common case.
	if b.Config.Environments == nil {
		return nil
	}

	// The "environments" key is set; validate and rewrite it to "targets".
	return b.Config.Mutate(func(v dyn.Value) (dyn.Value, error) {
		environments := v.Get("environments")
		targets := v.Get("targets")

		// Return an error if both "environments" and "targets" are set.
		if environments != dyn.NilValue && targets != dyn.NilValue {
			return dyn.NilValue, fmt.Errorf(
				"both 'environments' and 'targets' are specified; only 'targets' should be used: %s",
				environments.Location().String(),
			)
		}

		// Rewrite "environments" to "targets".
		if environments != dyn.NilValue && targets == dyn.NilValue {
			nv, err := dyn.Set(v, "targets", environments)
			if err != nil {
				return dyn.NilValue, err
			}
			// Drop the "environments" key.
			return dyn.Walk(nv, func(p dyn.Path, v dyn.Value) (dyn.Value, error) {
				switch len(p) {
				case 0:
					return v, nil
				case 1:
					if p[0] == dyn.Key("environments") {
						return v, dyn.ErrDrop
					}
				}
				return v, dyn.ErrSkip
			})
		}

		return v, nil
	})
}
