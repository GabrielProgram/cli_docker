package config

import "time"

type Value struct {
	v any
	l Location

	// Whether or not this value is an anchor.
	// If this node doesn't map to a type, we don't need to warn about it.
	anchor bool
}

// NilValue is equal to the zero-value of Value.
var NilValue = Value{}

// NewValue constructs a new Value with the given value and location.
func NewValue(v any, loc Location) Value {
	return Value{
		v: v,
		l: loc,
	}
}

func (v Value) AsMap() (map[string]Value, bool) {
	m, ok := v.v.(map[string]Value)
	return m, ok
}

func (v Value) Location() Location {
	return v.l
}

func (v Value) AsAny() any {
	switch vv := v.v.(type) {
	case map[string]Value:
		m := make(map[string]any)
		for k, v := range vv {
			m[k] = v.AsAny()
		}
		return m
	case []Value:
		a := make([]any, len(vv))
		for i, v := range vv {
			a[i] = v.AsAny()
		}
		return a
	case string:
		return vv
	case bool:
		return vv
	case int:
		return vv
	case int32:
		return vv
	case int64:
		return vv
	case float32:
		return vv
	case float64:
		return vv
	case time.Time:
		return vv
	case nil:
		return nil
	default:
		// Panic because we only want to deal with known types.
		panic("not handled")
	}
}

func (v Value) Get(key string) Value {
	m, ok := v.AsMap()
	if !ok {
		return NilValue
	}

	vv, ok := m[key]
	if !ok {
		return NilValue
	}

	return vv
}

func (v Value) Index(i int) Value {
	s, ok := v.v.([]Value)
	if !ok {
		return NilValue
	}

	if i < 0 || i >= len(s) {
		return NilValue
	}

	return s[i]
}

func (v Value) MarkAnchor() Value {
	return Value{
		v: v.v,
		l: v.l,

		anchor: true,
	}
}

func (v Value) IsAnchor() bool {
	return v.anchor
}
