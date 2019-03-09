package cmd

import (
	"fmt"
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

	log.Verbose(fmt.Sprintf("%s %s:\n%s", exe, strings.Join(args, " "), string(output)))

	//TODO: Make sure we don't lose error output when the command fails

	return string(output), err
}
