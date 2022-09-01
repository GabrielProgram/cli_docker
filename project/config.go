package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/databricks/bricks/folders"
	"github.com/databricks/databricks-sdk-go/service/clusters"

	"github.com/ghodss/yaml"
)

type Isolation string

const (
	None Isolation = ""
	Soft Isolation = "soft"
)

// ConfigFile is the name of project configuration file
const ConfigFile = "databricks.yml"

type Assertions struct {
	Groups            []string `json:"groups,omitempty"`
	Secrets           []string `json:"secrets,omitempty"`
	ServicePrincipals []string `json:"service_principals,omitempty"`
}

type Project struct {
	Name      string    `json:"name"`              // or do default from folder name?..
	Profile   string    `json:"profile,omitempty"` // rename?
	Isolation Isolation `json:"isolation,omitempty"`

	// development-time vs deployment-time resources
	DevCluster *clusters.ClusterInfo `json:"dev_cluster,omitempty"`

	// Assertions defines a list of configurations expected to be applied
	// to the workspace by a higher-privileged user (or service principal)
	// in order for the deploy command to work, as individual project teams
	// in almost all the cases don’t have admin privileges on Databricks
	// workspaces.
	//
	// This configuration simplifies the flexibility of individual project
	// teams, make jobs deployment easier and portable across environments.
	// This configuration block would contain the following entities to be
	// created by administrator users or admin-level automation, like Terraform
	// and/or SCIM provisioning.
	Assertions *Assertions `json:"assertions,omitempty"`
}

func (p Project) IsDevClusterDefined() bool {
	return reflect.ValueOf(p.DevCluster).IsZero()
}

// IsDevClusterJustReference denotes reference-only clusters.
// This conflicts with Soft isolation. Happens for cost-restricted projects,
// where there's only a single Shared Autoscaling cluster per workspace and
// general users have no ability to create other iteractive clusters.
func (p *Project) IsDevClusterJustReference() bool {
	if p.DevCluster.ClusterName == "" {
		return false
	}
	return reflect.DeepEqual(p.DevCluster, &clusters.ClusterInfo{
		ClusterName: p.DevCluster.ClusterName,
	})
}

// IsDatabricksProject returns true for folders with `databricks.yml`
// in the parent tree
func IsDatabricksProject() bool {
	_, err := findProjectRoot()
	return err == nil
}

func loadProjectConf() (prj Project, err error) {
	root, err := findProjectRoot()
	if err != nil {
		return
	}
	config, err := os.Open(fmt.Sprintf("%s/%s", root, ConfigFile))
	if err != nil {
		return
	}
	raw, err := ioutil.ReadAll(config)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(raw, &prj)
	if err != nil {
		return
	}
	return validateAndApplyProjectDefaults(prj)
}

func validateAndApplyProjectDefaults(prj Project) (Project, error) {
	// defaultCluster := clusters.ClusterInfo{
	// 	NodeTypeID: "smallest",
	// 	SparkVersion: "latest",
	// 	AutoterminationMinutes: 30,
	// }
	return prj, nil
}

func findProjectRoot() (string, error) {
	return folders.FindDirWithLeaf(ConfigFile)
}
