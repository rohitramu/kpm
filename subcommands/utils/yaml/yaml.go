package yaml

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ghodss/yaml"
)

// BytesToObject populates an object's properties from the contents of a yaml file.
// NOTE: ALWAYS pass objToPopulate as a pointer.
func BytesToObject(yamlBytes []byte, objToPopulate interface{}) error {
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
func ObjectToBytes(obj interface{}) ([]byte, error) {
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

// FromObject takes an interface, marshals it to yaml, and returns a string.
func FromObject(obj interface{}) (string, error) {
	var err error

	var data []byte
	data, err = yaml.Marshal(obj)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ToObject converts a YAML document into a map[string]interface{}.
func ToObject(yamlString string) (map[string]interface{}, error) {
	var err error

	var result = map[string]interface{}{}
	err = yaml.Unmarshal([]byte(yamlString), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
