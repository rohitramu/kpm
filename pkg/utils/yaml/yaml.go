package yaml

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ghodss/yaml"
)

// BytesToObject populates an object's properties from the contents of a yaml file.
// NOTE: ALWAYS pass objToPopulate as a pointer.
func BytesToObject(yamlBytes []byte, objToPopulate any) error {
	// Don't bother trying to deserialize bytes into an object if there are no bytes
	if yamlBytes != nil && len(yamlBytes) > 0 {
		// NOTE: ALWAYS pass "UnmarshalStrict()" a pointer rather than a real value
		var err = yaml.UnmarshalStrict(yamlBytes, objToPopulate, yaml.DisallowUnknownFields)
		if err != nil {
			return err
		}
	}

	return nil
}

// ObjectToBytes converts an object into yaml bytes.
func ObjectToBytes(obj any) ([]byte, error) {
	var err error

	// Check for a nil object
	if obj == nil {
		log.Panic("Object cannot be nil")
	}

	var result []byte
	result, err = yaml.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize object of type \"%s\" to yaml: %s", reflect.TypeOf(obj), err)
	}

	return result, nil
}
