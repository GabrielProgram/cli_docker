package convert

import (
	"reflect"
	"slices"
	"strings"
	"unicode"

	"github.com/databricks/cli/libs/dyn"
)

var skipFields = []string{"Format"}

// Converts a struct to map. Skips any nil fields.
// It uses `skipFields` to skip unnecessary fields.
// Uses `order` to define the order of keys in resulting outout
func ConvertToMapValue(strct any, order *dyn.Order, dst map[string]dyn.Value) (dyn.Value, error) {
	itemValue := reflect.ValueOf(strct)
	if itemValue.Kind() == reflect.Pointer {
		itemValue = itemValue.Elem()
	}
	for i := 0; i < itemValue.NumField(); i++ {
		if itemValue.Field(i).IsZero() {
			continue
		}

		f := itemValue.Type().Field(i)
		if slices.Contains(skipFields, f.Name) {
			continue
		}

		// If the field is not defined as json field, we're skipping it
		k, isJson := dyn.ConfigKey(strct, f.Name)
		if !isJson {
			continue
		}

		// If the value is already defined in destination, it means it was
		// manually set due to custom ordering or other customisation required
		// So we're skipping processing it again
		if _, ok := dst[k]; ok {
			continue
		}

		ref := dyn.NilValue
		nv, err := FromTyped(itemValue.Field(i).Interface(), ref)
		if err != nil {
			return dyn.NilValue, err
		}

		if nv.Kind() != dyn.KindNil {
			nv.SetLocation(dyn.Location{Line: order.Get(f.Name)})
			dst[k] = nv
		}
	}

	return dyn.V(dst), nil
}

func replaceNonAlphanumeric(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return r
	}
	return '_'
}

// We leave the full range of unicode letters in tact, but remove all "special" characters,
// including spaces and dots, which are not supported in e.g. experiment names or YAML keys.
func NormaliseString(name string) string {
	name = strings.ToLower(name)
	return strings.Map(replaceNonAlphanumeric, name)
}
