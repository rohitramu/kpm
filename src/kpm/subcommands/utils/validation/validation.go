package validation

import (
	"fmt"
	"regexp"
	"strings"

	"../logger"
)

// ValidatePackageName validates the given package's name.
func ValidatePackageName(packageName string) error {
	if len(strings.TrimSpace(packageName)) == 0 {
		return fmt.Errorf("Package name cannot be empty")
	}

	var isValid = CheckRegexMatch(packageName, "^[a-z](\\.?[a-z0-9])*$")

	if !isValid {
		return fmt.Errorf("Package name must consist of lowercase words which may be separated by dots: %s", packageName)
	}

	return nil
}

// ValidatePackageVersion validates the given package version.
func ValidatePackageVersion(packageVersion string, allowWildcards bool) error {
	if len(strings.TrimSpace(packageVersion)) == 0 {
		return fmt.Errorf("Package version string cannot be empty")
	}

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

// ValidateOutputName validates the output name when generating output.
func ValidateOutputName(outputName string) error {
	//TODO: Add validation
	return nil
}

// ExtractNameAndVersionFromFullPackageName returns the name and version of a template package, given the full package name.
func ExtractNameAndVersionFromFullPackageName(fullPackageName string) (packageName string, packageVersion string, err error) {
	// Split the file name to get the name and version
	var splitFileName = strings.SplitN(fullPackageName, "-", 2)

	// Check that this is a valid package name
	if len(splitFileName) != 2 {
		return "", "", fmt.Errorf("Full package name is an invalid format: %s", fullPackageName)
	}

	packageName = splitFileName[0]
	packageVersion = splitFileName[1]

	// Validate the package name
	if err := ValidatePackageName(packageName); err != nil {
		return "", "", err
	}

	// Validate the package version
	if err := ValidatePackageVersion(packageVersion, false); err != nil {
		return "", "", err
	}

	return packageName, packageVersion, nil
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
