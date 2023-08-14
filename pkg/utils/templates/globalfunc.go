package templates

// GetGlobalFuncMap returns the template functions which can be used in any context.
func GetGlobalFuncMap() map[string]interface{} {
	return map[string]interface{}{
		FuncNameIndex:    IndexFunc,
		FuncNameFromYaml: FromYamlFunc,
		FuncNameToYaml:   ToYamlFunc,
	}
}
