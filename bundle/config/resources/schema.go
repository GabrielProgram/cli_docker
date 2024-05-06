package resources

import (
	"github.com/databricks/databricks-sdk-go/service/catalog"
)

type Schema struct {
	// List of grants to apply on this schema.
	Grants []Grant `json:"grants,omitempty"`

	// Full name of the schema (catalog_name.schema_name). This value is read from
	// the terraform state after deployment succeeds.
	ID string `json:"id,omitempty" bundle:"readonly"`

	*catalog.CreateSchema

	ModifiedStatus ModifiedStatus `json:"modified_status,omitempty" bundle:"internal"`
}
