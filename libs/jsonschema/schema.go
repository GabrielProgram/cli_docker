package jsonschema

// defines schema for a json object
type Schema struct {
	// Type of the object
	Type Type `json:"type,omitempty"`

	// Description of the object. This is rendered as inline documentation in the
	// IDE. This is manually injected here using schema.Docs
	Description string `json:"description,omitempty"`

	// Schemas for the fields of an struct. The keys are the first json tag.
	// The values are the schema for the type of the field
	Properties map[string]*Schema `json:"properties,omitempty"`

	// The schema for all values of an array
	Items *Schema `json:"items,omitempty"`

	// The schema for any properties not mentioned in the Schema.Properties field.
	// this validates maps[string]any in bundle configuration
	// OR
	// A boolean type with value false. Setting false here validates that all
	// properties in the config have been defined in the json schema as properties
	//
	// Its type during runtime will either be *Schema or bool
	AdditionalProperties any `json:"additionalProperties,omitempty"`

	// Required properties for the object. Any fields missing the "omitempty"
	// json tag will be included
	Required []string `json:"required,omitempty"`

	// URI to a json schema
	Reference *string `json:"$ref,omitempty"`

	// Default value for the property / object
	Default any `json:"default,omitempty"`
}

type Type string

const (
	InvalidType Type = "invalid"
	BooleanType Type = "boolean"
	StringType  Type = "string"
	NumberType  Type = "number"
	ObjectType  Type = "object"
	ArrayType   Type = "array"
	IntegerType Type = "integer"
)
