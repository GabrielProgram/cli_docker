package dyn_test

import (
	"testing"

	"github.com/databricks/cli/libs/dyn"
	assert "github.com/databricks/cli/libs/dyn/dynassert"
)

func TestInvalidValue(t *testing.T) {
	// Assert that the zero value of [dyn.Value] is the invalid value.
	var zero dyn.Value
	assert.Equal(t, zero, dyn.InvalidValue)
}

func TestValueIsAnchor(t *testing.T) {
	var zero dyn.Value
	assert.False(t, zero.IsAnchor())
	mark := zero.MarkAnchor()
	assert.True(t, mark.IsAnchor())
}

func TestValueAsMap(t *testing.T) {
	var zeroValue dyn.Value
	_, ok := zeroValue.AsMap()
	assert.False(t, ok)

	var intValue = dyn.NewValue(1, dyn.Location{})
	_, ok = intValue.AsMap()
	assert.False(t, ok)

	var mapValue = dyn.NewValue(
		map[string]dyn.Value{
			"key": dyn.NewValue("value", dyn.Location{File: "file", Line: 1, Column: 2}),
		},
		dyn.Location{File: "file", Line: 1, Column: 2},
	)
	m, ok := mapValue.AsMap()
	assert.True(t, ok)
	assert.Equal(t, 1, m.Len())
}

func TestValueIsValid(t *testing.T) {
	var zeroValue dyn.Value
	assert.False(t, zeroValue.IsValid())
	var intValue = dyn.NewValue(1, dyn.Location{})
	assert.True(t, intValue.IsValid())
}

func TestAppendYamlLocation(t *testing.T) {
	var v dyn.Value

	// Add new locations
	v = v.AppendYamlLocation(dyn.Location{File: "file1", Line: 1, Column: 2})
	assert.Len(t, v.YamlLocations(), 1)

	v = v.AppendYamlLocation(dyn.Location{File: "file2", Line: 3, Column: 4})
	assert.Len(t, v.YamlLocations(), 2)

	// Ignore empty locations
	v = v.AppendYamlLocation(dyn.Location{})
	assert.Len(t, v.YamlLocations(), 2)

	// Ignore duplicate locations
	v = v.AppendYamlLocation(dyn.Location{File: "file1", Line: 1, Column: 2})
	assert.Len(t, v.YamlLocations(), 2)
}
