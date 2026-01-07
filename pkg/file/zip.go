package file

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ZipFiles returns a zip file from a list of file paths.
func ZipFiles(filePaths []string) ([]byte, error) {
	zipFile, err := os.CreateTemp("", "zip-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer zipFile.Close()
	defer os.Remove(zipFile.Name())

	// Create zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add each file to the zip
	for _, filePath := range filePaths {
		err := addFileToZip(zipWriter, filePath)
		if err != nil {
			return nil, err
		}
	}

	zipFile.Close()

	b, err := io.ReadAll(zipFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read zip file: %w", err)
	}

	return b, nil
}

// addFileToZip adds a single file to the zip archive
func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Get file info for the header
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", filePath, err)
	}

	// Create zip header
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return fmt.Errorf("failed to create zip header for %s: %w", filePath, err)
	}

	// Use only the base name in the zip (not the full path)
	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	// Create writer for this file
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create zip entry for %s: %w", filePath, err)
	}

	// Copy file content to zip
	_, err = io.Copy(writer, file)
	if err != nil {
		return fmt.Errorf("failed to write file %s to zip: %w", filePath, err)
	}

	return nil
}
