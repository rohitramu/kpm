package cmd

import (
	"os/exec"
	"strings"

	"../log"
)

// Exec runs a command on the command line.
func Exec(exe string, args ...string) (string, error) {
	var err error

	// Create the command
	var cmd = exec.Command(exe, args...)

	// Execute the command
	var output []byte
	output, err = cmd.CombinedOutput()

	log.Verbose("%s %s\n%s", exe, strings.Join(args, " "), string(output))

	return string(output), err
}
