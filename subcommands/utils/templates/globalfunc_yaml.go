package templates

import (
	"github.com/rohitramu/kpm/subcommands/utils/yaml"
)

// FuncNameToYaml is the name of the template function which converts objects to yaml strings.
const FuncNameToYaml = "toYaml"

// FuncNameFromYaml is the name of the template function which converts yaml strings to objects.
const FuncNameFromYaml = "fromYaml"

// ToYamlFunc converts objects to yaml strings.
func ToYamlFunc(value interface{}) (string, error) {
	return yaml.FromObject(value)
}

// FromYamlFunc converts yaml strings to objects.
func FromYamlFunc(yamlString string) (map[string]interface{}, error) {
	return yaml.ToObject(yamlString)
}
