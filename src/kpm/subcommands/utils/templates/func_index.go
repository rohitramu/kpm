package templates

import (
	"errors"
	"fmt"
	"reflect"

	"../types"
)

// FuncIndex gets a single value from a generic map (of any depth) given an ordered list of keys.
func FuncIndex(data types.GenericMap, keys ...string) (interface{}, error) {
	if len(keys) == 0 {
		return data, nil
	}

	var currentMap = data
	var result interface{}
	for _, key := range keys {
		var ok bool
		result, ok = currentMap[key]
		if !ok {
			var message string
			if !ok {
				message = fmt.Sprintf("Missing key of type: %s", reflect.TypeOf(key))
			} else {
				message = fmt.Sprintf("Missing key: %s", key)
			}
			return nil, errors.New(message)
		}

		// Try to assign the next map if the type is a map
		currentMap, ok = result.(types.GenericMap)
		if !ok {
			// If the type is not a map, set this to nil so we don't reuse the old map
			currentMap = nil
		}
	}

	return result, nil
}
