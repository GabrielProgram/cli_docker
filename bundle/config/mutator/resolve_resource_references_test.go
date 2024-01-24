package mutator

import (
	"context"
	"fmt"
	"testing"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/config"
	"github.com/databricks/cli/bundle/config/variable"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/databricks/databricks-sdk-go/experimental/mocks"
	"github.com/databricks/databricks-sdk-go/service/compute"
)

func TestResolveClusterReference(t *testing.T) {
	clusterRef1 := "Some Custom Cluster"
	clusterRef2 := "Some Other Name"
	justString := "random string"
	b := &bundle.Bundle{
		Config: config.Root{
			Variables: map[string]*variable.Variable{
				"my-cluster-id-1": {
					Lookup: &variable.Lookup{
						Cluster: clusterRef1,
					},
				},
				"my-cluster-id-2": {
					Lookup: &variable.Lookup{
						Cluster: clusterRef2,
					},
				},
				"some-variable": {
					Value: &justString,
				},
			},
		},
	}

	m := mocks.NewMockWorkspaceClient(t)
	b.SetWorkpaceClient(m.WorkspaceClient)
	clusterApi := m.GetMockClustersAPI()
	clusterApi.EXPECT().GetByClusterName(mock.Anything, clusterRef1).Return(&compute.ClusterDetails{
		ClusterId: "1234-5678-abcd",
	}, nil)
	clusterApi.EXPECT().GetByClusterName(mock.Anything, clusterRef2).Return(&compute.ClusterDetails{
		ClusterId: "9876-5432-xywz",
	}, nil)

	err := bundle.Apply(context.Background(), b, ResolveResourceReferences())
	require.NoError(t, err)
	require.Equal(t, "1234-5678-abcd", *b.Config.Variables["my-cluster-id-1"].Value)
	require.Equal(t, "9876-5432-xywz", *b.Config.Variables["my-cluster-id-2"].Value)
}

func TestResolveNonExistentClusterReference(t *testing.T) {
	clusterRef := "Random"
	justString := "random string"
	b := &bundle.Bundle{
		Config: config.Root{
			Variables: map[string]*variable.Variable{
				"my-cluster-id": {
					Lookup: &variable.Lookup{
						Cluster: clusterRef,
					},
				},
				"some-variable": {
					Value: &justString,
				},
			},
		},
	}

	m := mocks.NewMockWorkspaceClient(t)
	b.SetWorkpaceClient(m.WorkspaceClient)
	clusterApi := m.GetMockClustersAPI()
	clusterApi.EXPECT().GetByClusterName(mock.Anything, clusterRef).Return(nil, fmt.Errorf("ClusterDetails named '%s' does not exist", clusterRef))

	err := bundle.Apply(context.Background(), b, ResolveResourceReferences())
	require.ErrorContains(t, err, "failed to resolve cluster: Random, err: ClusterDetails named 'Random' does not exist")
}

func TestNoLookupIfVariableIsSet(t *testing.T) {
	clusterRef := "donotexist"
	b := &bundle.Bundle{
		Config: config.Root{
			Variables: map[string]*variable.Variable{
				"my-cluster-id": {
					Lookup: &variable.Lookup{
						Cluster: clusterRef,
					},
				},
			},
		},
	}

	m := mocks.NewMockWorkspaceClient(t)
	b.SetWorkpaceClient(m.WorkspaceClient)

	b.Config.Variables["my-cluster-id"].Set("random value")

	err := bundle.Apply(context.Background(), b, ResolveResourceReferences())
	require.NoError(t, err)
	require.Equal(t, "random value", *b.Config.Variables["my-cluster-id"].Value)
}
