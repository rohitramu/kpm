package yaml

import (
	"fmt"
	"reflect"

	"github.com/ghodss/yaml"

	"github.com/rohitramu/kpm/src/pkg/utils/log"
)

// BytesToObject populates an object's properties from the contents of a yaml file.
// NOTE: ALWAYS pass objToPopulate as a pointer to the object to populate.
func BytesToObject(yamlBytes []byte, objToPopulate any) error {
	// Don't bother trying to deserialize bytes into an object if there are no bytes.
	if len(yamlBytes) <= 0 {
		return nil
	}

	// NOTE: ALWAYS pass "UnmarshalStrict()" a pointer rather than a real value.
	if reflect.TypeOf(objToPopulate).Kind() != reflect.Pointer {
		log.Panicf("object to unmarshal into must be a pointer type")
	}

	// Unmarshal the YAML into the object.
	// TODO: Use standard yaml package once it supports strict unmarshalling: https://github.com/go-yaml/yaml/issues/460
	var err = yaml.UnmarshalStrict(yamlBytes, objToPopulate, yaml.DisallowUnknownFields)
	if err != nil {
		return err
	}

	return nil
}

// ObjectToBytes converts an object into yaml bytes.
func ObjectToBytes(obj any) ([]byte, error) {
	var err error

	// Check for a nil object
	if obj == nil {
		log.Panicf("Object cannot be nil")
	}

	var result []byte
	result, err = yaml.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize object of type \"%s\" to yaml: %s", reflect.TypeOf(obj), err)
	}

	return result, nil
}
