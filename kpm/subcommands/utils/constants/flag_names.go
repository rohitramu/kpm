package constants

var (
	// LogLevelFlagName is the minimum severity log level to output.
	LogLevelFlagName = "logLevel"

	// PackageVersionFlagName is the version of a template package.
	PackageVersionFlagName = "packageVersion"

	// ParametersFileFlagName is the file that contains the parameters for a template.
	ParametersFileFlagName = "parametersFile"

	// OutputNameFlagName is the name of the generated configuration.
	OutputNameFlagName = "outputName"

	// OutputDirFlagName is the output directory.
	OutputDirFlagName = "outputDir"

	// KpmHomeDirFlagName is the home directory for KPM, for the current user.
	KpmHomeDirFlagName = "kpmHomeDir"

	// DockerRegistryURLFlagName is the Docker registry URL to use when pushing or pulling a package.
	DockerRegistryURLFlagName = "dockerRegistryUrl"

	// DockerNamespaceFlagName is the docker namespace to use when pushing or pulling a package.
	DockerNamespaceFlagName = "dockerNamespace"
)
