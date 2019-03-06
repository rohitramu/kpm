package docker

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"../logger"
	"../types"
	"github.com/mholt/archiver"
)

// createTar creates a new in-memory Docker tar file which can be used to build Docker images by
// making a request to the Docker daemon.
func createTar(dockerfileName string, dockerfile string, dirToCopy string) (*bytes.Buffer, error) {
	var err error

	// Get the base name of the directory to copy so we can use it as the parent directory at the destination
	var dirToCopyBaseName = filepath.Base(dirToCopy)

	// Create byte stream
	var byteStream = new(bytes.Buffer)

	var tarFile = archiver.Tar{
		ImplicitTopLevelFolder: false,
		ContinueOnError:        false,
	}

	// Create a new instance of the tar file in memory
	err = tarFile.Create(byteStream)
	if err != nil {
		return nil, err
	}
	defer tarFile.Close()

	// Add Dockerfile
	err = tarFile.Write(archiver.File{
		FileInfo: archiver.FileInfo{
			FileInfo: types.MockFileInfo{
				MockName:    dockerfileName,
				MockSize:    int64(len(dockerfile)),
				MockMode:    os.ModePerm,
				MockModTime: time.Now(),
			},
			CustomName: dockerfileName,
		},
		ReadCloser: ioutil.NopCloser(bytes.NewBufferString(dockerfile)),
	})

	// Walk the package directory and add all files
	err = filepath.Walk(dirToCopy, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Don't add this file to the tar archive if it is a directory or symbolic link
		if !fileInfo.Mode().IsRegular() {
			return nil
		}

		// Get the internal file path
		var internalFilePath string
		internalFilePath, err = filepath.Rel(dirToCopy, filePath)
		if err != nil {
			return err
		}

		// Prepend the full package name to the path so the files don't all get copied to the root
		internalFilePath = filepath.Join(dirToCopyBaseName, internalFilePath)

		// Replace all backslashes with forward slashes since Docker uses Unix file paths
		internalFilePath = filepath.ToSlash(internalFilePath)

		logger.Default.Verbose.Println(fmt.Sprintf("Include: %s", internalFilePath))

		// Open the file for reading
		var fileReader *os.File
		fileReader, err = os.Open(filePath)
		if err != nil {
			return err
		}

		// Create the file header
		var header *tar.Header
		header, err = tar.FileInfoHeader(fileInfo, "")
		if err != nil {
			return err
		}
		header.Name = internalFilePath + "/" + header.Name

		// Write the file to the archive
		err = tarFile.Write(archiver.File{
			FileInfo: archiver.FileInfo{
				FileInfo:   fileInfo,
				CustomName: internalFilePath,
			},
			Header:     header,
			ReadCloser: fileReader,
		})
		fileReader.Close()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	logger.Default.Verbose.Println(fmt.Sprintf("Tar file size: %d bytes", len(byteStream.Bytes())))

	var tempOutputDir = filepath.Join(os.TempDir(), "kpm")
	err = os.MkdirAll(tempOutputDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	var tempOutputPath = filepath.Join(tempOutputDir, fmt.Sprintf("%s.tar", dirToCopyBaseName))
	err = ioutil.WriteFile(tempOutputPath, byteStream.Bytes(), os.ModePerm)
	if err != nil {
		panic(err)
	}
	logger.Default.Verbose.Println(fmt.Sprintf("Wrote temp file to: %s", tempOutputPath))

	return byteStream, nil
}
