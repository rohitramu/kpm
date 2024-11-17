package constants

// The version of the CLI app.
//
// Can be set during build: `go build -ldflags="-X 'github.com/rohitramu/kpm/cli/model/utils/constants.VersionString=0.0.0'"`
var VersionString = "0.0.0"

// When this binary was built.
var BuildTimestampUTC = ""

// The Git commit hash that this binary was built from.
var GitCommitHash = ""

// If true, the source code was modified from the commit hash specified in GitCommitHash.
var IsSourceModified = true
