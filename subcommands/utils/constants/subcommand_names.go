package constants

// ListCmdName is the name of the subcommand which lists the packages that are available for use.
const ListCmdName = "list"

// RemoveCmdName is the name of the subcommand which removes a package from the local KPM repository.
const RemoveCmdName = "remove"

// PurgeCmdName is the name of the subcommand which purges the local KPM repository.
const PurgeCmdName = "purge"

// PackCmdName is the name of the subcommand which makes a template package ready for use.
const PackCmdName = "pack"

// UnpackCmdName is the name of the subcommand which exports a template package.
const UnpackCmdName = "unpack"

// InspectCmdName is the name of the subcommand which outputs the contents of the default parameters file in a package.
const InspectCmdName = "inspect"

// RunCmdName is the name of the subcommand which generates output using a template package.
const RunCmdName = "run"

// DockerCmdName is the name of the subcommand which contains subcommands for interacting with KPM packages in Docker registries.
const DockerCmdName = "docker"

// VersionsCmdName is the name of the subcommand which prints all tags in a remote Docker repository.
const VersionsCmdName = "versions"

// PushCmdName is the name of the subcommand which pushes a template package to a remote repository.
const PushCmdName = "push"

// PullCmdName is the name of the subcommand which pulls a template package from a remote repository.
const PullCmdName = "pull"
