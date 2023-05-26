package internal

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/databricks/cli/cmd/workspace"
)

var clusterId string

func TestAccClustersList(t *testing.T) {
	t.Log(GetEnvOrSkipTest(t, "CLOUD_ENV"))

	stdout, stderr := RequireSuccessfulRun(t, "clusters", "list")
	outStr := stdout.String()
	assert.Contains(t, outStr, "ID")
	assert.Contains(t, outStr, "Name")
	assert.Contains(t, outStr, "State")
	assert.Equal(t, "", stderr.String())

	idRegExp := regexp.MustCompile(`[0-9]{4}\-[0-9]{6}-[a-z0-9]{8}`)
	clusterId = idRegExp.FindString(outStr)
	fmt.Println(clusterId)
	assert.NotEmpty(t, clusterId)
}

func TestAccClustersGet(t *testing.T) {
	t.Log(GetEnvOrSkipTest(t, "CLOUD_ENV"))

	stdout, stderr := RequireSuccessfulRun(t, "clusters", "get", clusterId)
	outStr := stdout.String()
	assert.Contains(t, outStr, fmt.Sprintf(`"cluster_id":"%s"`, clusterId))
	assert.Equal(t, "", stderr.String())
}
