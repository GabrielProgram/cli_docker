package mutator

import (
	"context"
	"fmt"
	"testing"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/config"
	"github.com/databricks/cli/internal/build"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	currentVersion string
	constraint     string
	expectedError  string
}

func TestVerifyCliVersion(t *testing.T) {
	testCases := []testCase{
		{
			currentVersion: "0.0.1",
		},
		{
			currentVersion: "0.0.1",
			constraint:     "0.100.0",
			expectedError:  "Databricks CLI version constraint not satisfied. Required: 0.100.0, current: 0.0.1",
		},
		{
			currentVersion: "0.0.1",
			constraint:     ">= 0.100.0",
			expectedError:  "Databricks CLI version constraint not satisfied. Required: >= 0.100.0, current: 0.0.1",
		},
		{
			currentVersion: "0.100.0",
			constraint:     "0.100.0",
		},
		{
			currentVersion: "0.100.1",
			constraint:     "0.100.0",
			expectedError:  "Databricks CLI version constraint not satisfied. Required: 0.100.0, current: 0.100.1",
		},
		{
			currentVersion: "0.100.1",
			constraint:     ">= 0.100.0",
		},
		{
			currentVersion: "0.100.0",
			constraint:     "<= 1.0.0",
		},
		{
			currentVersion: "1.0.0",
			constraint:     "<= 1.0.0",
		},
		{
			currentVersion: "1.0.0",
			constraint:     "<= 0.100.0",
			expectedError:  "Databricks CLI version constraint not satisfied. Required: <= 0.100.0, current: 1.0.0",
		},
		{
			currentVersion: "0.99.0",
			constraint:     ">= 0.100.0, <= 0.100.2",
			expectedError:  "Databricks CLI version constraint not satisfied. Required: >= 0.100.0, <= 0.100.2, current: 0.99.0",
		},
		{
			currentVersion: "0.100.0",
			constraint:     ">= 0.100.0, <= 0.100.2",
		},
		{
			currentVersion: "0.100.1",
			constraint:     ">= 0.100.0, <= 0.100.2",
		},
		{
			currentVersion: "0.100.2",
			constraint:     ">= 0.100.0, <= 0.100.2",
		},
		{
			currentVersion: "0.101.0",
			constraint:     ">= 0.100.0, <= 0.100.2",
			expectedError:  "Databricks CLI version constraint not satisfied. Required: >= 0.100.0, <= 0.100.2, current: 0.101.0",
		},
		{
			currentVersion: "0.100.0-beta",
			constraint:     ">= 0.100.0, <= 0.100.2",
			expectedError:  "Databricks CLI version constraint not satisfied. Required: >= 0.100.0, <= 0.100.2, current: 0.100.0-beta",
		},
		{
			currentVersion: "0.100.0-beta",
			constraint:     ">= 0.100.0-0, <= 0.100.2-0",
		},
		{
			currentVersion: "0.100.1-beta",
			constraint:     ">= 0.100.0-0, <= 0.100.2-0",
		},
		{
			currentVersion: "0.100.3-beta",
			constraint:     ">= 0.100.0, <= 0.100.2",
			expectedError:  "Databricks CLI version constraint not satisfied. Required: >= 0.100.0, <= 0.100.2, current: 0.100.3-beta",
		},
		{
			currentVersion: "0.100.123",
			constraint:     "0.100.*",
		},
		{
			currentVersion: "0.100.123",
			constraint:     "^0.100",
		},
	}

	t.Cleanup(func() {
		build.SetBuildVersion(build.DefaultSemver)
	})

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("testcase #%d", i), func(t *testing.T) {
			build.SetBuildVersion(tc.currentVersion)
			b := &bundle.Bundle{
				Config: config.Root{
					Bundle: config.Bundle{
						DatabricksCliVersion: tc.constraint,
					},
				},
			}
			diags := bundle.Apply(context.Background(), b, VerifyCliVersion())
			if tc.expectedError != "" {
				require.NotEmpty(t, diags)
				require.Equal(t, tc.expectedError, diags.Error().Error())
			} else {
				require.Empty(t, diags)
			}
		})
	}
}
