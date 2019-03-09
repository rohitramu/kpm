package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"../log"
)

// MaxOutputNameLength is the maximum number of characters allowed in the output name.
const MaxOutputNameLength = 64

// ValidatePackageName validates the given package's name.
func ValidatePackageName(packageName string) error {
	var err error

	if len(strings.TrimSpace(packageName)) == 0 {
		return fmt.Errorf("Package name cannot be empty")
	}

	// Check namespace segments
	var nameSegments = strings.Split(packageName, "/")
	for i, namespaceSegment := range nameSegments {
		// The final segment is the unqualified name, so don't check if it is a valid namespace segment
		if i+1 < len(nameSegments) {
			err = ValidateNamespaceSegment(namespaceSegment)
			if err != nil {
				return err
			}
		}
	}

	// Final name segment is the unqualified package name
	var unqualifiedName = nameSegments[len(nameSegments)-1]

	// Check the unqualified name
	var isValid bool
	isValid, err = CheckRegexMatch(unqualifiedName, "^[a-z](\\.?[a-z0-9])*$")
	if err != nil {
		log.Panic(err)
	}

	// Return an error if the name is not valid
	if !isValid {
		return fmt.Errorf("The unqualified package name (i.e. ignoring the namespace) must consist of lowercase words which may be separated by dots: %s", packageName)
	}

	return nil
}

// ValidatePackageVersion validates the given package version.
func ValidatePackageVersion(packageVersion string, allowWildcards bool) error {
	var err error

	// Check for empty string
	if len(strings.TrimSpace(packageVersion)) == 0 {
		return fmt.Errorf("Package version string cannot be empty")
	}

	// Check for wildcards
	if !allowWildcards && strings.ContainsRune(packageVersion, '*') {
		return fmt.Errorf("Package version cannot contain wildcards: %s", packageVersion)
	}

	// Check for zero version
	var zeroVersion = "0.0.0"
	if packageVersion == zeroVersion {
		return fmt.Errorf("Package version cannot be \"%s\"", zeroVersion)
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

	// Check whether the version string satisfies the regex
	var isValid bool
	isValid, err = CheckRegexMatch(packageVersion, fmt.Sprintf(fullRegex, segmentRegex, segmentRegex, segmentRegex))
	if err != nil {
		log.Panic(err)
	}

	// Return error if the version string did not satisfy the regex
	if !isValid {
		return fmt.Errorf("Package version must solely consist of digits, be in the form \"major.minor.revision\" with no leading zeros in any segment, and be greater than \"0.0.0\": %s", packageVersion)
	}

	return nil
}

// ValidateOutputName validates the output name when generating output.
func ValidateOutputName(outputName string) error {
	var err error

	// Check for empty string
	if outputName == "" {
		return fmt.Errorf("Output name cannot be empty")
	}

	// Check length
	if len(outputName) > MaxOutputNameLength {
		return fmt.Errorf("Output name cannot be longer than %d characters: %s", MaxOutputNameLength, outputName)
	}

	var alphaNumeric = "[a-zA-Z0-9]"
	var symbols = "[-_]"
	var regex = fmt.Sprintf("^%s+(%s?%s)+$", alphaNumeric, symbols, alphaNumeric)
	var matched bool
	matched, err = regexp.MatchString(regex, outputName)
	if err != nil {
		log.Panic(err)
	}
	if !matched {
		return fmt.Errorf("Output name must only consist of letters and numbers, optionally separated by dashes and/or underscores")
	}

	return nil
}

// ValidateNamespaceSegment validates an image namespace's segment.
func ValidateNamespaceSegment(namespaceSegment string) error {
	var err error

	// Check for empty string
	if namespaceSegment == "" {
		return errors.New("Namespace segment cannot be empty")
	}

	// Build the regex
	var regex = "^[a-z0-9]+$"

	// Check if the value satisfies the regex
	var isValid bool
	isValid, err = CheckRegexMatch(namespaceSegment, regex)
	if err != nil {
		log.Panic(err)
	}

	// Return an error if the value doesn't satisfy the regex
	if !isValid {
		return fmt.Errorf("Invalid namespace segment: %s", namespaceSegment)
	}

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
func CheckRegexMatch(stringToCheck string, regex string) (bool, error) {
	var isMatch, err = regexp.MatchString(regex, stringToCheck)
	if err != nil {
		return false, err
	}

	return isMatch, nil
}

// GetStringOrDefault returns testValue if it is not null, otherwise returns defaultValue.
func GetStringOrDefault(testValue *string, defaultValue string) string {
	if testValue == nil {
		return defaultValue
	}

	return *testValue
}

// GetStringOrError returns testValue if it is not null, otherwise returns an error.
func GetStringOrError(testValue *string, valueName string) (string, error) {
	if testValue == nil {
		return "", fmt.Errorf("Value cannot be nil: %s", valueName)
	}

	return *testValue, nil
}
