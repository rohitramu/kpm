package templates

import (
	"github.com/rohitramu/kpm/src/pkg/utils/yaml"
)

// FuncNameToYaml is the name of the template function which converts objects to yaml strings.
const FuncNameToYaml = "toYaml"

// FuncNameFromYaml is the name of the template function which converts yaml strings to objects.
const FuncNameFromYaml = "fromYaml"

// ToYamlFunc converts objects to yaml strings.
func ToYamlFunc(value any) (string, error) {
	var err error
	var resultBytes []byte
	resultBytes, err = yaml.ObjectToBytes(value)
	if err != nil {
		return "", err
	}

	return string(resultBytes), nil
}

// FromYamlFunc converts yaml strings to objects.
func FromYamlFunc(yamlString string) (*map[string]any, error) {
	var err error
	var resultObj = new(map[string]any)
	err = yaml.BytesToObject([]byte(yamlString), yamlString)
	if err != nil {
		return nil, err
	}

	return resultObj, nil
}
