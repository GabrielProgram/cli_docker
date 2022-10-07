package internal

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/databricks/bricks/cmd/api"
)

func TestAccApiGet(t *testing.T) {
	t.Log(GetEnvOrSkipTest(t, "CLOUD_ENV"))

	stdout, _, err := run(t, "api", "get", "/api/2.0/preview/scim/v2/Me")
	require.NoError(t, err)

	// Deserialize SCIM API response.
	var out map[string]any
	err = json.Unmarshal(stdout.Bytes(), &out)
	require.NoError(t, err)

	// Assert that the output somewhat makes sense for the SCIM API.
	assert.Equal(t, true, out["active"])
	assert.NotNil(t, out["id"])
}

func TestAccApiPost(t *testing.T) {
	t.Log(GetEnvOrSkipTest(t, "CLOUD_ENV"))

	dbfsPath := filepath.Join("/tmp/bricks/integration", RandomName("api-post"))
	requestPath := writeFile(t, "body.json", fmt.Sprintf(`{
		"path": "%s"
	}`, dbfsPath))

	// Post to mkdir
	{
		_, _, err := run(t, "api", "post", "--body=@"+requestPath, "/api/2.0/dbfs/mkdirs")
		require.NoError(t, err)
	}

	// Post to delete
	{
		_, _, err := run(t, "api", "post", "--body=@"+requestPath, "/api/2.0/dbfs/delete")
		require.NoError(t, err)
	}
}
