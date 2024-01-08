package mutator

import (
	"context"
	"testing"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/config"
	"github.com/databricks/cli/bundle/config/variable"
	"github.com/databricks/databricks-sdk-go/service/compute"
	"github.com/stretchr/testify/require"
)

type MockClusterService struct{}

// ChangeOwner implements compute.ClustersService.
func (MockClusterService) ChangeOwner(ctx context.Context, request compute.ChangeClusterOwner) error {
	panic("unimplemented")
}

// Create implements compute.ClustersService.
func (MockClusterService) Create(ctx context.Context, request compute.CreateCluster) (*compute.CreateClusterResponse, error) {
	panic("unimplemented")
}

// Delete implements compute.ClustersService.
func (MockClusterService) Delete(ctx context.Context, request compute.DeleteCluster) error {
	panic("unimplemented")
}

// Edit implements compute.ClustersService.
func (MockClusterService) Edit(ctx context.Context, request compute.EditCluster) error {
	panic("unimplemented")
}

// Events implements compute.ClustersService.
func (MockClusterService) Events(ctx context.Context, request compute.GetEvents) (*compute.GetEventsResponse, error) {
	panic("unimplemented")
}

// Get implements compute.ClustersService.
func (MockClusterService) Get(ctx context.Context, request compute.GetClusterRequest) (*compute.ClusterDetails, error) {
	panic("unimplemented")
}

// GetPermissionLevels implements compute.ClustersService.
func (MockClusterService) GetPermissionLevels(ctx context.Context, request compute.GetClusterPermissionLevelsRequest) (*compute.GetClusterPermissionLevelsResponse, error) {
	panic("unimplemented")
}

// GetPermissions implements compute.ClustersService.
func (MockClusterService) GetPermissions(ctx context.Context, request compute.GetClusterPermissionsRequest) (*compute.ClusterPermissions, error) {
	panic("unimplemented")
}

// List implements compute.ClustersService.
func (MockClusterService) List(ctx context.Context, request compute.ListClustersRequest) (*compute.ListClustersResponse, error) {
	return &compute.ListClustersResponse{
		Clusters: []compute.ClusterDetails{
			{ClusterId: "1234-5678-abcd", ClusterName: "Some Custom Cluster"},
			{ClusterId: "9876-5432-xywz", ClusterName: "Some Other Name"},
		},
	}, nil
}

// ListNodeTypes implements compute.ClustersService.
func (MockClusterService) ListNodeTypes(ctx context.Context) (*compute.ListNodeTypesResponse, error) {
	panic("unimplemented")
}

// ListZones implements compute.ClustersService.
func (MockClusterService) ListZones(ctx context.Context) (*compute.ListAvailableZonesResponse, error) {
	panic("unimplemented")
}

// PermanentDelete implements compute.ClustersService.
func (MockClusterService) PermanentDelete(ctx context.Context, request compute.PermanentDeleteCluster) error {
	panic("unimplemented")
}

// Pin implements compute.ClustersService.
func (MockClusterService) Pin(ctx context.Context, request compute.PinCluster) error {
	panic("unimplemented")
}

// Resize implements compute.ClustersService.
func (MockClusterService) Resize(ctx context.Context, request compute.ResizeCluster) error {
	panic("unimplemented")
}

// Restart implements compute.ClustersService.
func (MockClusterService) Restart(ctx context.Context, request compute.RestartCluster) error {
	panic("unimplemented")
}

// SetPermissions implements compute.ClustersService.
func (MockClusterService) SetPermissions(ctx context.Context, request compute.ClusterPermissionsRequest) (*compute.ClusterPermissions, error) {
	panic("unimplemented")
}

// SparkVersions implements compute.ClustersService.
func (MockClusterService) SparkVersions(ctx context.Context) (*compute.GetSparkVersionsResponse, error) {
	panic("unimplemented")
}

// Start implements compute.ClustersService.
func (MockClusterService) Start(ctx context.Context, request compute.StartCluster) error {
	panic("unimplemented")
}

// Unpin implements compute.ClustersService.
func (MockClusterService) Unpin(ctx context.Context, request compute.UnpinCluster) error {
	panic("unimplemented")
}

// UpdatePermissions implements compute.ClustersService.
func (MockClusterService) UpdatePermissions(ctx context.Context, request compute.ClusterPermissionsRequest) (*compute.ClusterPermissions, error) {
	panic("unimplemented")
}

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

	b.WorkspaceClient().Clusters.WithImpl(MockClusterService{})

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

	b.WorkspaceClient().Clusters.WithImpl(MockClusterService{})

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

	b.WorkspaceClient().Clusters.WithImpl(MockClusterService{})
	b.Config.Variables["my-cluster-id"].Set("random value")

	err := bundle.Apply(context.Background(), b, ResolveResourceReferences())
	require.NoError(t, err)
	require.Equal(t, "random value", *b.Config.Variables["my-cluster-id"].Value)
}
