package file

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Save saves the given multipart file to the specified destination (dst). It creates directories if they don't exist.
func Save(file *multipart.FileHeader, dst string) error {
	// Ensure the directory exists
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create the destination file,
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the file content to the destination
	_, err = io.Copy(destFile, src)
	return err
}

type FileInfo struct {
	Content  io.ReadCloser
	MimeType string
}

// Read opens the file at the given destination, determines its MIME type, and returns the content and MIME type.
func Read(dst string) (*FileInfo, error) {
	file, err := os.Open(dst)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("File does not exist")
		}
		return nil, err
	}

	// Read a small portion of the file to determine its MIME type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		file.Close() // Close the file if there's an error
		return nil, err
	}

	// Reset the file cursor to the beginning
	if _, err := file.Seek(0, 0); err != nil {
		file.Close()
		return nil, err
	}

	// Detect MIME type
	mimeType := http.DetectContentType(buffer[:n])

	return &FileInfo{
		Content:  file,
		MimeType: mimeType,
	}, nil
}

// Delete removes the file at the specified destination (dst).
func Delete(dst string) error {
	// Attempt to remove the file
	if err := os.Remove(dst); err != nil {
		if os.IsNotExist(err) {
			return errors.New("File does not exist")
		}
		return err
	}
	return nil
}
