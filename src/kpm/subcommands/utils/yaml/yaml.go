package yaml

import (
	"github.com/ghodss/yaml"

	"../logger"
)

// BytesToObject populates an object's properties from the contents of a yaml file.
func BytesToObject(yamlBytes []byte, objToPopulate interface{}) {
	// Don't bother trying to deserialize bytes into an object if there are no bytes
	if yamlBytes != nil && len(yamlBytes) > 0 {
		// NOTE: ALWAYS pass "UnmarshalStrict()" a pointer rather than a real value
		var err = yaml.UnmarshalStrict(yamlBytes, objToPopulate)
		if err != nil {
			logger.Default.Error.Fatalln(err)
		}
	}
}
