package docker

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"../logger"
	"github.com/mholt/archiver"
)

// extractTar extracts a tar file into the destination directory, creating directories as required.
func extractTar(tarData io.ReadCloser, tarDirToExtract string, destinationDir string) error {
	var err error
	var ok bool

	// Convert path inside tar to a path that is supported by the current OS
	var tarPath = filepath.FromSlash(tarDirToExtract)

	// Make sure we have an absolute path for the destination directory
	if !filepath.IsAbs(destinationDir) {
		logger.Default.Error.Panicln(fmt.Sprintf("Destination directory is not an absolute path: %s", destinationDir))
	}

	// Create a new instance of an in-memory tar file
	var tarFile = archiver.Tar{}

	// Open the tar file for reading from the tar data stream
	err = tarFile.Open(tarData, -1)
	if err != nil {
		return fmt.Errorf("Failed to open tar file stream: %s", err)
	}

	for file, eofErr := tarFile.Read(); eofErr != io.EOF; file, eofErr = tarFile.Read() {
		// Don't try to write objects that aren't regular files
		if file.IsDir() || !file.Mode().IsRegular() {
			continue
		}

		// Get file header
		var header *tar.Header
		header, ok = file.Header.(*tar.Header)
		if !ok {
			return fmt.Errorf("Found invalid file header for file: %s", file.Name())
		}

		// Replace the path separators with the separator for the current OS
		var relativePath = filepath.FromSlash(header.Name)

		// Get the relative path
		relativePath, err = filepath.Rel(tarPath, relativePath)
		if err != nil {
			return fmt.Errorf("Failed to get relative path inside tar file: %s", err)
		}

		logger.Default.Verbose.Println(fmt.Sprintf("Extracting: %s", header.Name))

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
		defer func() {
			var closeErr = outFile.Close()
			if closeErr != nil {
				if err != nil {
					err = fmt.Errorf("Failed to close output file stream:\n%s\n%s", closeErr, err)
				} else {
					err = closeErr
				}
			}
		}()
		if err != nil {
			return err
		}

		// Write the file to the filesystem
		var numBytesWritten int64
		numBytesWritten, err = io.Copy(outFile, file.ReadCloser)

		defer logger.Default.Verbose.Println(fmt.Sprintf("Extracted %d bytes: %s", numBytesWritten, outputPath))
	}

	return err
}
