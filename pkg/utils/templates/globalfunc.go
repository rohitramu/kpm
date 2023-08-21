package templates

// GetGlobalFuncMap returns the template functions which can be used in any context.
func GetGlobalFuncMap() map[string]any {
	return map[string]any{
		FuncNameIndex:    IndexFunc,
		FuncNameFromYaml: FromYamlFunc,
		FuncNameToYaml:   ToYamlFunc,
	}
}
