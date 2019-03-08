package cmd

import (
	"os/exec"
)

// Exec runs a command on the command line.
func Exec(exe string, args ...string) (string, error) {
	var err error

	var exePath string
	exePath, err = exec.LookPath(exe)
	if err != nil {
		return "", err
	}

	// Create the command
	var cmd = exec.Command(exePath, args...)

	// Execute the command
	var output []byte
	output, err = cmd.CombinedOutput()

	//TODO: Make sure we don't lose error output when the command fails

	return string(output), err
}
