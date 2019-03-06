package templatefuncs

import (
	"fmt"

	"../types"
)

// Index gets a single value from a generic map (of any depth) given an ordered list of keys
func Index(data interface{}, keys ...string) (interface{}, error) {
	if len(keys) == 0 {
		return data, nil
	}

	// Make sure the data is either a map or a value type
	var currentMap, ok = data.(*types.GenericMap)
	if !ok {
		if len(keys) == 0 {
			// Data is a value type and it was expected (i.e. no keys were supplied), so return it as-is
			return data, nil
		}

		// Unexpected data type - it is not a map, but keys were supplied
		return nil, fmt.Errorf("Invalid object for index: %s", data)
	}

	var result interface{}
	for _, key := range keys {
		// Make sure that the current map is not nil, otherwise there were too many keys provided
		if currentMap == nil {
			return nil, fmt.Errorf("Too many keys provided - a key was provided, but the data is not a map: %s", key)
		}

		// Get the value
		result, ok := (*currentMap)[key]
		if !ok {
			return nil, fmt.Errorf("Missing key: %s", key)
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
