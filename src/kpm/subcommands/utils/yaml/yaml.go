package yaml

import (
	"github.com/ghodss/yaml"

	"../logger"
	"../types"
)

// BytesToMap generates an generic object from the contents of a yaml file.
func BytesToMap(yamlBytes []byte) *types.GenericMap {
	// NOTE: ALWAYS pass "UnmarshalStrict()" a pointer rather than a real value
	var result = &types.GenericMap{}
	var err = yaml.UnmarshalStrict(yamlBytes, result)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	return result
}

// BytesToPackageInfo generates a PackageInfo object from the contents of a yaml file.
func BytesToPackageInfo(packageInfoBytes []byte) *types.PackageInfo {
	// NOTE: ALWAYS pass "UnmarshalStrict()" a pointer rather than a real value
	var result = new(types.PackageInfo)
	var err = yaml.UnmarshalStrict(packageInfoBytes, result)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	return result
}
