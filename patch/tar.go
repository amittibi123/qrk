package patch

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CreateTarFromDirectory takes a source directory path and writes its contents
// into a single uncompressed .tar file using streaming.
func CreateTarFromDirectory(sourceDir string, outputTarPath string) error {
	// 1. Create the output tar file on disk
	tarFile, err := os.Create(outputTarPath)
	if err != nil {
		return fmt.Errorf("failed to create output tar file: %w", err)
	}
	defer tarFile.Close()

	// 2. Create a tar writer that wraps our file
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	// 3. Walk through the source directory recursively
	err = filepath.Walk(sourceDir, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a standard tar header based on the file information
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return fmt.Errorf("failed to create tar header for %s: %w", currentPath, err)
		}

		// Update the header name to use the relative path inside the directory
		relPath, err := filepath.Rel(sourceDir, currentPath)
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write the header to the tar archive
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header for %s: %w", currentPath, err)
		}

		// If the current item is a directory, we only need the header, so skip writing data
		if info.IsDir() {
			return nil
		}

		// 4. Stream the file content into the tar writer
		fileToTar, err := os.Open(currentPath)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", currentPath, err)
		}
		defer fileToTar.Close()

		// io.Copy streams the data directly using an internal buffer (efficient RAM usage)
		_, err = io.Copy(tarWriter, fileToTar)
		if err != nil {
			return fmt.Errorf("failed to copy file data for %s: %w", currentPath, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the directory: %w", err)
	}

	return nil
}
