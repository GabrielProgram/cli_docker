package dyn_test

import (
	"testing"

	"github.com/databricks/cli/libs/dyn"
	"github.com/stretchr/testify/assert"
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
	m, ok := zeroValue.AsMap()
	assert.False(t, ok)
	assert.Nil(t, m)

	var intValue = dyn.NewValue(1, dyn.Location{})
	m, ok = intValue.AsMap()
	assert.False(t, ok)
	assert.Nil(t, m)

	var mapValue = dyn.NewValue(
		map[string]dyn.Value{
			"key": dyn.NewValue("value", dyn.Location{File: "file", Line: 1, Column: 2}),
		},
		dyn.Location{File: "file", Line: 1, Column: 2},
	)
	m, ok = mapValue.AsMap()
	assert.True(t, ok)
	assert.Len(t, m, 1)
}

func TestValueIsValid(t *testing.T) {
	var zeroValue dyn.Value
	assert.False(t, zeroValue.IsValid())
	var intValue = dyn.NewValue(1, dyn.Location{})
	assert.True(t, intValue.IsValid())
}
