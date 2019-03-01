package validation

import (
	"fmt"
	"regexp"

	"../logger"
)

// ValidatePackageName validates the given package's name.
func ValidatePackageName(packageName string) error {
	var isValid = CheckRegexMatch(packageName, "^[a-z](\\.?[a-z0-9])*$")

	if !isValid {
		return fmt.Errorf("Package name must consist of lowercase words which may be separated by dots: %s", packageName)
	}

	return nil
}

// ValidatePackageVersion validates the given package version.
func ValidatePackageVersion(packageVersion string, allowWildcards bool) error {
	// Regex for each segment that has a valid integer
	var segmentRegex = "(0|[1-9][0-9]*)"

	// Overall regex - pick one depending on whether wildcards are allowed or not
	var fullRegex string
	if allowWildcards {
		fullRegex = "^(\\*|%s\\.(\\*|%s\\.(\\*|%s)))$"
	} else {
		fullRegex = "^%s\\.%s\\.%s$"
	}

	var isValid = packageVersion != "0.0.0" && CheckRegexMatch(packageVersion, fmt.Sprintf(fullRegex, segmentRegex, segmentRegex, segmentRegex))

	if !isValid {
		return fmt.Errorf("Package version must solely consist of digits, be in the form \"major.minor.revision\" with no leading zeros, and be greater than \"0.0.0\": %s", packageVersion)
	}

	return nil
}

// GetPackageNameWithVersion returns the full package name with version.
func GetPackageNameWithVersion(packageName string, packageVersion string) string {
	return fmt.Sprintf("%s-%s", packageName, packageVersion)
}

// CheckRegexMatch checks whether a string satisfies the given regex expression.
func CheckRegexMatch(stringToCheck string, regex string) bool {
	var isMatch, err = regexp.MatchString(regex, stringToCheck)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	return isMatch
}

// GetStringOrDefault returns testValue if it is not null, otherwise returns defaultValue.
func GetStringOrDefault(testValue *string, defaultValue string) string {
	if testValue == nil {
		return defaultValue
	}

	return *testValue
}

// GetStringOrFail returns testValue if it is not null, otherwise throws a fatal error.
func GetStringOrFail(testValue *string, valueName string) string {
	if testValue == nil {
		logger.Default.Error.Fatalln(fmt.Sprintf("Value cannot be nil: %s", valueName))
	}

	return *testValue
}
