package template_package

import (
	"fmt"
	"os/user"
	"path"
	"regexp"

	"github.com/rohitramu/kpm/pkg/utils/constants"
	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"golang.org/x/exp/maps"
)

func GenerateSampleTemplatePackage(packageAbsolutePath string, packageName string) error {
	var err error

	var destinationDirIsEmpty bool
	destinationDirIsEmpty, err = files.DirIsEmpty(packageAbsolutePath, "template package")
	if err != nil {
		return err
	}
	if !destinationDirIsEmpty {
		return fmt.Errorf("cannot generate template package in \"%s\" because it is not empty", packageAbsolutePath)
	}

	// Get current user info.
	var currentUser string
	if currentOsUser, err := user.Current(); err != nil {
		currentUser = "foobar"
		log.Warningf("Failed to retrieve current user info.")
	} else {
		if regexp, err := regexp.Compile("[a-zA-Z0-9_-]+"); err != nil {
			log.Warningf("Failed to compile regex for sanitizing username.")
		} else {
			var matches = regexp.FindAllString(currentOsUser.Username, 100)
			if len(matches) > 0 {
				currentUser = matches[len(matches)-1]
			}
		}
	}

	// Define the files to be created (path mapped to human-friendly file description and content).
	var genFiles = map[string][2]string{
		// Package info
		path.Join(packageAbsolutePath, constants.PackageInfoFileName): {
			"package info",
			fmt.Sprintf(`# %s
name: %s/helloworld
version: 0.0.1
`, constants.PackageInfoFileName, currentUser),
		},

		// Interface
		path.Join(packageAbsolutePath, constants.InterfaceFileName): {
			"interface",
			fmt.Sprintf(`# %s
username: {{ .name.first }} {{ .name.last }}
`, constants.InterfaceFileName),
		},

		// Parameters
		path.Join(packageAbsolutePath, constants.ParametersFileName): {
			"parameters",
			fmt.Sprintf(`# %s
name:
  first: Foo
  last: Bar
`, constants.ParametersFileName),
		},

		// Template
		path.Join(packageAbsolutePath, constants.TemplatesDirName, "hello.txt"): {
			"\"hello\" template",
			fmt.Sprintf(`# %s
Hello, {{ .values.username }}!

Are you enjoying KPM?
Let me know if you have any issues: https://github.com/rohitramu/kpm/issues

Kind regards,
Rohit (creator of KPM)
`, "hello.txt"),
		},
	}

	for _, filePath := range maps.Keys(genFiles) {
		var entry = genFiles[filePath]
		var lowercaseHumanFriendlyName, content = entry[0], entry[1]

		if err = files.CreateFile(filePath, lowercaseHumanFriendlyName, content); err != nil {
			return fmt.Errorf("failed to generate template package: %s", err)
		}
	}

	return nil
}
