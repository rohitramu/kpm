package templates

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/types"
)

// FuncNameIndex is the name of the "index" template function.
const FuncNameIndex = "index"

// IndexFunc gets a single value from a generic map (of any depth), given an ordered list of keys.
func IndexFunc(data interface{}, keys ...string) (interface{}, error) {
	var ok bool

	// If there are no keys, return the object as-is
	if len(keys) == 0 {
		return data, nil
	}

	// Make sure the data is either a map or a value type
	var currentMap *types.GenericMap
	currentMap, ok = data.(*types.GenericMap)
	if !ok {
		if len(keys) == 0 {
			// Data is a value type and it was expected (i.e. no keys were supplied), so return it as-is
			return data, nil
		}

		// Unexpected data type - it is not a map, but keys were supplied
		return nil, fmt.Errorf("Invalid object supplied to the \"%s\" function: %s", FuncNameIndex, data)
	}

	var result interface{}
	for _, key := range keys {
		// Make sure that the current map is not nil, otherwise there were too many keys provided
		if currentMap == nil {
			return nil, fmt.Errorf("Too many keys provided to the \"%s\" function - a key was provided, but the data is not a map: %s", FuncNameIndex, key)
		}

		// Get the value
		result, ok = (*currentMap)[key]
		if !ok {
			return nil, fmt.Errorf("Missing key found in the object provided to the \"%s\" function: %s", FuncNameIndex, key)
		}

		// Try to assign the next map if the type is a map
		currentMap, ok = result.(*types.GenericMap)
		if !ok {
			// If the type is not a map, set this to nil so we don't reuse the old map
			currentMap = nil
		}
	}

	return result, nil
}
