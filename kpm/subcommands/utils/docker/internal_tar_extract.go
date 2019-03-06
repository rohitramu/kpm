package docker

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"../logger"
	"github.com/mholt/archiver"
)

func extractTar(tarData io.ReadCloser, destinationDir string) error {
	var err error
	var ok bool

	// Make sure we have an absolute path for the destination directory
	if !filepath.IsAbs(destinationDir) {
		logger.Default.Error.Panicln(fmt.Sprintf("Destination directory is not an absolute path: %s", destinationDir))
	}

	// Create a new instance of an in-memory tar file
	var tarFile = archiver.Tar{}

	// Open the tar file for reading from the tar data stream
	err = tarFile.Open(tarData, 0)
	if err != nil {
		return err
	}

	var file archiver.File
	for file, err = tarFile.Read(); err != io.EOF; file, err = tarFile.Read() {
		// Don't try to write objects that aren't regular files
		if file.IsDir() || !file.Mode().IsRegular() {
			continue
		}

		var header *tar.Header
		header, ok = file.Header.(*tar.Header)
		if !ok {
			return fmt.Errorf("Found invalid file header for file: %s", file.Name())
		}

		// Convert relative path to Unix path to have predictable path separators
		var relativePath = filepath.ToSlash(header.Name)

		// Replace the path separators with the separator for the current OS
		relativePath = strings.Replace(relativePath, "/", string(filepath.Separator), -1)

		// Join the path with the destination directory to get the final output path
		var outputPath = filepath.Join(destinationDir, relativePath)

		// Get the directory for the file
		var outputDir = filepath.Dir(outputPath)

		// Create the directory if it doesn't exist
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return err
		}

		// Create a new file on the filesystem to write to
		var outFile *os.File
		outFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Write the file to the filesystem
		var numBytesWritten int64
		numBytesWritten, err = io.Copy(outFile, file.ReadCloser)

		logger.Default.Verbose.Println(fmt.Sprintf("Extracted %d bytes: %s", numBytesWritten, outputPath))
	}

	return nil
}
