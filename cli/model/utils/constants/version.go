package constants

// The version of the CLI app.
//
// Can be set during build: `go build -ldflags "-X cli.commands.versionString v<my_version_number>"`.
var VersionString = "0.0.0"

// When this binary was built.
var BuildTimestampUTC = ""

// The Git commit hash that this binary was built from.
var GitCommitHash = ""

// If true, the source code was modified from the commit hash specified in GitCommitHash.
var IsSourceModified = true
